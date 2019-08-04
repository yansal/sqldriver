package sqldriver

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/lib/pq"
)

func assertf(t *testing.T, ok bool, format string, args ...interface{}) {
	t.Helper()
	if !ok {
		t.Errorf(format, args...)
	}
}

func TestSelectNow(t *testing.T) {
	pqconnector, err := pq.NewConnector("sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	connector := &Connector{
		Connector: pqconnector,
		QueryContextFunc: func(ctx context.Context, query string, args []driver.NamedValue, duration time.Duration, err error) {
			t.Logf("query=%q args=%+v duration=%s err=%v", query, args, duration, err)
			assertf(t, query == `select now()`, "got %q", query)
			assertf(t, len(args) == 0, "got %d args", len(args))
			assertf(t, err == nil, "got %v", err)
		},
		NextFunc: func(dest []driver.Value, duration time.Duration, err error) {
			t.Logf("dest=%+v duration=%s err=%v", dest, duration, err)
			assertf(t, len(dest) == 1, "got %d args", len(dest))
			ti, ok := dest[0].(time.Time)
			assertf(t, ok, "got %T", dest[0])
			assertf(t, !ti.IsZero(), "got %v", dest[0])
			assertf(t, err == nil, "got %v", err)
		},
	}

	db := sql.OpenDB(connector)
	defer db.Close()

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		t.Fatal(err)
	}
	var (
		query = `select now()`
		now   time.Time
	)
	if err := db.QueryRowContext(ctx, query).Scan(&now); err != nil {
		t.Fatal(err)
	}
	if now.IsZero() {
		t.Error("expected now to be non-zero")
	}
}
