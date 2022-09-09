package sqllogger

import (
	"context"
	"database/sql/driver"
)

type ConnQueryTracer interface {
	ConnQueryTrace(connID int64, query string, args []driver.Value) func(rowsID int64, err error)
}

type ConnQueryContextTracer interface {
	ConnQueryContextTrace(ctx context.Context, connID int64, query string, args []driver.NamedValue) func(rowsID int64, err error)
}

type StmtQueryTracer interface {
	StmtQueryTrace(stmtID int64, query string, args []driver.Value) func(rowsID int64, err error)
}

type StmtQueryContextTracer interface {
	StmtQueryContextTrace(ctx context.Context, stmtID int64, query string, args []driver.NamedValue) func(rowsID int64, err error)
}
