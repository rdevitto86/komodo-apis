package postgres

import (
	"context"
	"errors"
	logger "komodo-forge-apis-go/loggers/runtime"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Option func(*clientOptions)

type clientOptions struct {
 	workerPoolCfg *WorkerPoolConfig
 	batchWriterCfg *BatchWriterConfig
}

type Client struct {
 	pool   *pgxpool.Pool
 	worker *WorkerPool
 	batch  *BatchWriter
}

// Creates an option to configure the worker pool
func WithWorkerPool(cfg WorkerPoolConfig) Option {
 	return func(opts *clientOptions) {
 		config := cfg
 		opts.workerPoolCfg = &config
 	}
}

// Creates an option to configure the batch writer
func WithBatchWriter(cfg BatchWriterConfig) Option {
 	return func(opts *clientOptions) {
 		config := cfg
 		opts.batchWriterCfg = &config
 	}
}

// Creates a new client with the given configuration and options
func New(ctx context.Context, cfg Config, opts ...Option) (*Client, error) {
 	conn, err := cfg.connectionString()
 	if err != nil { return nil, err }

 	pcfg, err := pgxpool.ParseConfig(conn)
 	if err != nil { return nil, err }

 	if cfg.ConnectTimeout > 0 {
 		pcfg.ConnConfig.ConnectTimeout = cfg.ConnectTimeout
 	}
 	if cfg.StatementCacheCap > 0 {
 		pcfg.ConnConfig.StatementCacheCapacity = cfg.StatementCacheCap
 	}
 	if cfg.PreferSimpleProtocol {
 		pcfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
 	}
 	if cfg.TLSConfig != nil {
 		pcfg.ConnConfig.TLSConfig = cfg.TLSConfig
 	}
 	if cfg.ApplicationName != "" {
 		pcfg.ConnConfig.RuntimeParams["application_name"] = cfg.ApplicationName
 	}
 	if cfg.MaxConns > 0 {
 		pcfg.MaxConns = cfg.MaxConns
 	}
 	if cfg.MinConns > 0 {
 		pcfg.MinConns = cfg.MinConns
 	}
 	if cfg.MaxConnLifetime > 0 {
 		pcfg.MaxConnLifetime = cfg.MaxConnLifetime
 	}
 	if cfg.MaxConnIdleTime > 0 {
 		pcfg.MaxConnIdleTime = cfg.MaxConnIdleTime
 	}
 	if cfg.HealthCheckPeriod > 0 {
 		pcfg.HealthCheckPeriod = cfg.HealthCheckPeriod
 	}

 	pool, err := pgxpool.NewWithConfig(ctx, pcfg)
 	if err != nil { return nil, err }

 	if err := pool.Ping(ctx); err != nil {
 		pool.Close()
 		return nil, err
 	}

 	options := &clientOptions{}
 	for _, opt := range opts {
 		if opt != nil {
 			opt(options)
 		}
 	}

 	client := &Client{pool: pool}

	// Initialize worker pool
 	if options.workerPoolCfg != nil {
 		wp, err := NewWorkerPool(*options.workerPoolCfg)
 		if err != nil {
 			pool.Close()
 			return nil, err
 		}
 		client.worker = wp
 	}

	// Initialize batch writer
 	if options.batchWriterCfg != nil {
 		bw, err := NewBatchWriter(pool, *options.batchWriterCfg)
 		if err != nil {
 			if client.worker != nil {
 				_ = client.worker.Shutdown(context.Background())
 			}
 			pool.Close()
 			return nil, err
 		}
 		client.batch = bw
 	}

 	logger.Info("postgres client initialized")
 	return client, nil
}

// Returns the underlying pool
func (client *Client) Pool() *pgxpool.Pool {
 	if client == nil { return nil }
 	return client.pool
}

// Pings the database
func (client *Client) Ping(ctx context.Context) error {
 	if client == nil || client.pool == nil { return ErrNotInitialized }
 	return client.pool.Ping(ctx)
}

// Closes the client
func (client *Client) Close() error {
 	if client == nil { return nil }

 	if client.batch != nil {
 		_ = client.batch.Shutdown(context.Background())
 	}
 	if client.worker != nil {
 		_ = client.worker.Shutdown(context.Background())
 	}
 	if client.pool != nil {
 		client.pool.Close()
 	}
 	return nil
}

// Executes a query and returns a command tag
func (client *Client) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
 	if client == nil || client.pool == nil {
 		return pgconn.CommandTag{}, ErrNotInitialized
 	}
 	return client.pool.Exec(ctx, sql, args...)
}

// Executes a query and returns a set of rows
func (client *Client) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
 	if client == nil || client.pool == nil {
 		return nil, ErrNotInitialized
 	}
 	return client.pool.Query(ctx, sql, args...)
}

// Executes a query and returns a single row
func (client *Client) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if client == nil || client.pool == nil {
		return errRow{err: ErrNotInitialized}
	}
	return client.pool.QueryRow(ctx, sql, args...)
}

// Begins a transaction
func (client *Client) BeginTransaction(ctx context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
 	if client == nil || client.pool == nil {
 		return nil, ErrNotInitialized
 	}
 	return client.pool.BeginTx(ctx, opts)
}

// Starts a transaction and executes the given function
func (client *Client) WithTransaction(ctx context.Context, opts pgx.TxOptions, fn func(pgx.Tx) error) error {
 	if fn == nil {
 		return errors.New("transaction function is required")
 	}

 	tx, err := client.BeginTransaction(ctx, opts)
 	if err != nil { return err }

 	defer tx.Rollback(ctx)

 	if err := fn(tx); err != nil { return err }
 	return tx.Commit(ctx)
}

// Submits a job to the worker pool
func (client *Client) Submit(ctx context.Context, job Job) error {
 	if client == nil { return ErrNotInitialized }

 	if client.worker == nil {
 		if job == nil {
 			return errors.New("nil job")
 		}
 		return job(ctx)
 	}
 	return client.worker.Submit(ctx, job)
}

// Enqueues a job to execute a SQL statement
func (client *Client) EnqueueExec(ctx context.Context, sql string, args ...any) (<-chan error, error) {
 	if client == nil { return nil, ErrNotInitialized }

 	if client.batch == nil {
 		resp := make(chan error, 1)
 		go func() {
 			_, err := client.Exec(ctx, sql, args...)
 			resp <- err
 			close(resp)
 		}()
 		return resp, nil
 	}
 	return client.batch.EnqueueExec(ctx, sql, args...)
}
