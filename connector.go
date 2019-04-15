package sqldriver

import (
	"context"
	"database/sql/driver"
)

type Connector struct {
	Connector driver.Connector
	BeginTxFunc
	CommitFunc
	NextFunc
	QueryContextFunc
	RollbackFunc
}

func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	conn, err := c.Connector.Connect(ctx)
	return wrappedConn{
		Conn:             conn,
		BeginTxFunc:      c.BeginTxFunc,
		CommitFunc:       c.CommitFunc,
		NextFunc:         c.NextFunc,
		QueryContextFunc: c.QueryContextFunc,
		RollbackFunc:     c.RollbackFunc,
	}, err
}

func (c *Connector) Driver() driver.Driver {
	return Driver{
		Driver:           c.Connector.Driver(),
		BeginTxFunc:      c.BeginTxFunc,
		CommitFunc:       c.CommitFunc,
		NextFunc:         c.NextFunc,
		QueryContextFunc: c.QueryContextFunc,
		RollbackFunc:     c.RollbackFunc,
	}
}
