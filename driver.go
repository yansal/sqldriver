package sqldriver

import (
	"context"
	"database/sql/driver"
	"time"
)

type Driver struct {
	driver.Driver
	QueryContextFunc
	NextFunc
}

type QueryContextFunc func(string, []driver.NamedValue, time.Duration, error)
type NextFunc func([]driver.Value, time.Duration, error)

func (d Driver) Open(name string) (driver.Conn, error) {
	conn, err := d.Driver.Open(name)
	return wrappedConn{
		Conn:             conn,
		QueryContextFunc: d.QueryContextFunc,
		NextFunc:         d.NextFunc,
	}, err
}

type wrappedConn struct {
	driver.Conn
	QueryContextFunc
	NextFunc
}

func (w wrappedConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (rows driver.Rows, err error) {
	start := time.Now()
	rows, err = w.Conn.(driver.QueryerContext).QueryContext(ctx, query, args)
	if w.QueryContextFunc != nil {
		w.QueryContextFunc(query, args, time.Since(start), err)
	}
	return wrappedRows{
		Rows:     rows,
		NextFunc: w.NextFunc,
	}, err
}

type wrappedRows struct {
	driver.Rows
	NextFunc
}

func (w wrappedRows) Next(dest []driver.Value) (err error) {
	start := time.Now()
	err = w.Rows.Next(dest)
	if w.NextFunc != nil {
		w.NextFunc(dest, time.Since(start), err)
	}
	return err
}
