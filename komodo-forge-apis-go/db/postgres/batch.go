package postgres

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BatchWriterConfig struct {
 	QueueSize    int
 	MaxBatchSize int
 	MaxWait      time.Duration
 	FlushTimeout time.Duration
}

type BatchWriter struct {
 	pool    			*pgxpool.Pool
 	tasks   			chan execTask
 	wg      			sync.WaitGroup
 	closed  			atomic.Bool
 	maxBatch 			int
 	maxWait  			time.Duration
 	flushTimeout 	time.Duration
}

type execTask struct {
 	ctx  context.Context
 	sql  string
 	args []any
 	resp chan error
}

// Creates a new batch writer
func NewBatchWriter(pool *pgxpool.Pool, cfg BatchWriterConfig) (*BatchWriter, error) {
 	if pool == nil {
 		return nil, ErrNotInitialized
 	}
 	if cfg.QueueSize <= 0 {
 		cfg.QueueSize = 2048
 	}
 	if cfg.MaxBatchSize <= 0 {
 		cfg.MaxBatchSize = 128
 	}
 	if cfg.MaxWait <= 0 {
 		cfg.MaxWait = 10 * time.Millisecond
 	}
 	if cfg.FlushTimeout <= 0 {
 		cfg.FlushTimeout = 5 * time.Second
 	}

 	bw := &BatchWriter{
 		pool:         pool,
 		tasks:        make(chan execTask, cfg.QueueSize),
 		maxBatch:     cfg.MaxBatchSize,
 		maxWait:      cfg.MaxWait,
 		flushTimeout: cfg.FlushTimeout,
 	}

 	bw.wg.Add(1)
 	go bw.loop()

 	return bw, nil
}

// Enqueues a job to execute a SQL statement
func (bw *BatchWriter) EnqueueExec(ctx context.Context, sql string, args ...any) (<-chan error, error) {
 	if bw == nil {
 		return nil, ErrNotInitialized
 	}
 	if bw.closed.Load() {
 		return nil, ErrWorkerPoolClosed
 	}

 	resp := make(chan error, 1)
 	task := execTask{ctx: ctx, sql: sql, args: args, resp: resp}

 	select {
 		case bw.tasks <- task:
 			return resp, nil
 		case <-ctx.Done():
 			close(resp)
 			return nil, ctx.Err()
 		default:
 			close(resp)
 			return nil, ErrQueueFull
 	}
}

// Helper function to flush pending tasks 
func (bw *BatchWriter) loop() {
 	defer bw.wg.Done()

 	timer := time.NewTimer(bw.maxWait)
 	defer timer.Stop()

 	var pending []execTask

	// Flush pending tasks
 	flush := func() {
 		if len(pending) == 0 { return }

 		batch := &pgx.Batch{}
 		resps := make([]chan error, 0, len(pending))

		// Build batch
 		for _, t := range pending {
 			if t.ctx != nil {
 				if err := t.ctx.Err(); err != nil {
 					t.resp <- err
 					close(t.resp)
 					continue
 				}
 			}
 			batch.Queue(t.sql, t.args...)
 			resps = append(resps, t.resp)
 		}

 		pending = pending[:0]
 		if len(resps) == 0 { return }

		// Flush batch
 		flushCtx, cancel := context.WithTimeout(context.Background(), bw.flushTimeout)
 		results := bw.pool.SendBatch(flushCtx, batch)

		// Send batch
 		for _, chnl := range resps {
 			_, err := results.Exec()
 			chnl <- err
 			close(chnl)
 		}

 		_ = results.Close()
 		cancel()
 	}

	// Main loop
 	for {
 		select {
 			case t, ok := <-bw.tasks:
 				if !ok {
 					flush()
 					return
 				}

 				pending = append(pending, t)
 				if len(pending) >= bw.maxBatch {
 					flush()
 					if !timer.Stop() {
 						select {
 							case <-timer.C:
 							default:
 						}
 					}

 					timer.Reset(bw.maxWait)
 				}
 			case <-timer.C:
 				flush()
 				timer.Reset(bw.maxWait)
 		}
 	}
}

// Shuts down the batch writer
func (bw *BatchWriter) Shutdown(ctx context.Context) error {
 	if bw == nil { return nil }

 	if bw.closed.CompareAndSwap(false, true) {
 		close(bw.tasks)
 	}

 	done := make(chan struct{})
 	go func() {
 		bw.wg.Wait()
 		close(done)
 	}()

 	select {
 		case <-done:
 			return nil
 		case <-ctx.Done():
 			return ctx.Err()
 	}
}
