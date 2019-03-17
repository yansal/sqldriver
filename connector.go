package sqldriver

import (
	"context"
	"database/sql/driver"
)

type Connector struct {
	Connector driver.Connector
	QueryContextFunc
	NextFunc
}

func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	conn, err := c.Connector.Connect(ctx)
	return wrappedConn{
		Conn:             conn,
		QueryContextFunc: c.QueryContextFunc,
		NextFunc:         c.NextFunc,
	}, err
}

func (c *Connector) Driver() driver.Driver {
	return Driver{
		Driver:           c.Connector.Driver(),
		QueryContextFunc: c.QueryContextFunc,
		NextFunc:         c.NextFunc,
	}
}
