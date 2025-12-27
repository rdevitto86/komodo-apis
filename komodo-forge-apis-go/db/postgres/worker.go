package postgres

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

type Job func(context.Context) error

type WorkerPoolConfig struct {
 	Workers   int
 	QueueSize int
}

type WorkerPool struct {
 	jobs   chan jobRequest
 	wg     sync.WaitGroup
 	closed atomic.Bool
}

type jobRequest struct {
 	ctx  context.Context
 	job  Job
 	resp chan error
}

// Creates a new worker pool.
func NewWorkerPool(cfg WorkerPoolConfig) (*WorkerPool, error) {
 	if cfg.Workers <= 0 {
 		return nil, errors.New("worker pool workers must be > 0")
 	}
 	if cfg.QueueSize <= 0 {
 		cfg.QueueSize = cfg.Workers * 64
 	}

 	pool := &WorkerPool{jobs: make(chan jobRequest, cfg.QueueSize)}
 	pool.wg.Add(cfg.Workers)

 	for i := 0; i < cfg.Workers; i++ {
 		go pool.workerLoop()
 	}
 	return pool, nil
}

func (pool *WorkerPool) workerLoop() {
 	defer pool.wg.Done()

 	for req := range pool.jobs {
 		if req.job == nil {
 			if req.resp != nil {
 				req.resp <- errors.New("nil job")
 				close(req.resp)
 			}
 			continue
 		}

 		err := req.job(req.ctx)
 		if req.resp != nil {
 			req.resp <- err
 			close(req.resp)
 		}
 	}
}

// Submits a job to the worker pool asynchronously
func (pool *WorkerPool) SubmitAsync(ctx context.Context, job Job) (<-chan error, error) {
 	if pool == nil {
 		return nil, ErrNotInitialized
 	}
 	if pool.closed.Load() {
 		return nil, ErrWorkerPoolClosed
 	}

 	resp := make(chan error, 1)
 	req := jobRequest{ctx: ctx, job: job, resp: resp}

 	select {
 		case pool.jobs <- req:
 			return resp, nil
 		case <-ctx.Done():
 			close(resp)
 			return nil, ctx.Err()
 		default:
 			close(resp)
 			return nil, ErrQueueFull
 	}
}

// Submits a job to the worker pool
func (pool *WorkerPool) Submit(ctx context.Context, job Job) error {
 	chnl, err := pool.SubmitAsync(ctx, job)
 	if err != nil { return err }

 	select {
 		case err := <-chnl:
 			return err
 		case <-ctx.Done():
 			return ctx.Err()
 	}
}

// Shuts down the worker pool
func (pool *WorkerPool) Shutdown(ctx context.Context) error {
 	if pool == nil { return nil }
 	if pool.closed.CompareAndSwap(false, true) {
 		close(pool.jobs)
 	}

	// Wait for all workers to finish
 	done := make(chan struct{})
 	go func() {
 		pool.wg.Wait()
 		close(done)
 	}()

 	select {
 		case <-done:
 			return nil
 		case <-ctx.Done():
 			return ctx.Err()
 	}
}
