package sqldriver

import (
	"context"
	"database/sql/driver"
	"time"
)

type Driver struct {
	driver.Driver
	BeginTxFunc
	CommitFunc
	ExecContextFunc
	NextFunc
	QueryContextFunc
	RollbackFunc
}

type BeginTxFunc func(driver.TxOptions, time.Duration, error)
type CommitFunc func(time.Duration, error)
type ExecContextFunc func(context.Context, string, []driver.NamedValue, time.Duration, error)
type NextFunc func([]driver.Value, time.Duration, error)
type QueryContextFunc func(context.Context, string, []driver.NamedValue, time.Duration, error)
type RollbackFunc func(time.Duration, error)

func (d Driver) Open(name string) (driver.Conn, error) {
	conn, err := d.Driver.Open(name)
	return wrappedConn{
		Conn:             conn,
		BeginTxFunc:      d.BeginTxFunc,
		CommitFunc:       d.CommitFunc,
		ExecContextFunc:  d.ExecContextFunc,
		NextFunc:         d.NextFunc,
		QueryContextFunc: d.QueryContextFunc,
		RollbackFunc:     d.RollbackFunc,
	}, err
}

type wrappedConn struct {
	driver.Conn
	BeginTxFunc
	CommitFunc
	ExecContextFunc
	NextFunc
	QueryContextFunc
	RollbackFunc
}

func (w wrappedConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	start := time.Now()
	tx, err := w.Conn.(driver.ConnBeginTx).BeginTx(ctx, opts)
	if w.BeginTxFunc != nil {
		w.BeginTxFunc(opts, time.Since(start), err)
	}
	return wrappedTx{
		Tx:           tx,
		CommitFunc:   w.CommitFunc,
		RollbackFunc: w.RollbackFunc,
	}, err
}

type wrappedTx struct {
	driver.Tx
	CommitFunc
	RollbackFunc
}

func (w wrappedTx) Commit() error {
	start := time.Now()
	err := w.Tx.Commit()
	if w.CommitFunc != nil {
		w.CommitFunc(time.Since(start), err)
	}
	return err
}

func (w wrappedTx) Rollback() error {
	start := time.Now()
	err := w.Tx.Rollback()
	if w.RollbackFunc != nil {
		w.RollbackFunc(time.Since(start), err)
	}
	return err
}

func (w wrappedConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	start := time.Now()
	rows, err := w.Conn.(driver.QueryerContext).QueryContext(ctx, query, args)
	if w.QueryContextFunc != nil {
		w.QueryContextFunc(ctx, query, args, time.Since(start), err)
	}
	return wrappedRows{
		Rows:     rows,
		NextFunc: w.NextFunc,
	}, err
}

func (w wrappedConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	start := time.Now()
	result, err := w.Conn.(driver.ExecerContext).ExecContext(ctx, query, args)
	if w.ExecContextFunc != nil {
		w.ExecContextFunc(ctx, query, args, time.Since(start), err)
	}
	return result, err
}

type wrappedRows struct {
	driver.Rows
	NextFunc
}

func (w wrappedRows) Next(dest []driver.Value) error {
	start := time.Now()
	err := w.Rows.Next(dest)
	if w.NextFunc != nil {
		w.NextFunc(dest, time.Since(start), err)
	}
	return err
}
