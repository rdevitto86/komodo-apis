package postgres

import "errors"

var (
	ErrNotInitialized   = errors.New("postgres client not initialized")
	ErrWorkerPoolClosed = errors.New("worker pool is closed")
	ErrQueueFull        = errors.New("queue full")
)

type errRow struct{ 
	err error 
}

func (errRow errRow) Scan(_ ...any) error { 
	return errRow.err 
}