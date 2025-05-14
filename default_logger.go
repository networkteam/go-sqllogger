package sqllogger

import (
	"context"
	"database/sql/driver"
)

// StdLogger is an interface to adapt the DefaultSQLLogger to the standard library log.Logger or other log frameworks
type StdLogger interface {
	Printf(format string, args ...interface{})
}

// NewDefaultSQLLogger creates a new default SQL logger with sensible defaults
//
// A *log.Logger can be passed or any other implementation of the StdLogger interface.
func NewDefaultSQLLogger(log StdLogger) *DefaultSQLLogger {
	return &DefaultSQLLogger{
		log:        log,
		Enabled:    true,
		LogConnect: true,
		LogClose:   false,
	}
}

// DefaultSQLLogger is an implementation of the Logger interface logging to a *log.Logger from the standard library
type DefaultSQLLogger struct {
	log StdLogger

	// Enabled sets, whether
	Enabled bool

	LogConnect bool
	LogClose   bool
}

var _ SQLLogger = &DefaultSQLLogger{}

// TxRollback satisfies Logger interface
func (dl *DefaultSQLLogger) TxRollback(txID int64) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("  TX(%d) ► Rollback", txID)
}

// TxCommit satisfies Logger interface
func (dl *DefaultSQLLogger) TxCommit(txID int64) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("  TX(%d) ► Commit", txID)
}

// RowsClose satisfies Logger interface
func (dl *DefaultSQLLogger) RowsClose(rowsID int64) {
	if !dl.Enabled {
		return
	}
	if dl.LogClose {
		dl.log.Printf("ROWS(%d) ► Close", rowsID)
	}
}

// Connect satisfies Logger interface
func (dl *DefaultSQLLogger) Connect(ctx context.Context, connID int64) {
	if !dl.Enabled {
		return
	}
	if !dl.LogConnect {
		return
	}
	dl.log.Printf("Connect → CONN(%d)", connID)
}

// ConnBegin satisfies Logger interface
func (dl *DefaultSQLLogger) ConnBegin(ctx context.Context, connID, txID int64, opts driver.TxOptions) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("CONN(%d) ► Begin -> TX(%d)", connID, txID)
}

// ConnPrepare satisfies Logger interface
func (dl *DefaultSQLLogger) ConnPrepare(connID, stmtID int64, query string) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("CONN(%d) ► Prepare(%s) → STMT(%d)", connID, query, stmtID)
}

// ConnPrepareContext satisfies Logger interface
func (dl *DefaultSQLLogger) ConnPrepareContext(ctx context.Context, connID int64, stmtID int64, query string) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("CONN(%d) ► Prepare(%s) → STMT(%d)", connID, query, stmtID)
}

// ConnQuery satisfies Logger interface
func (dl *DefaultSQLLogger) ConnQuery(connID, rowsID int64, query string, args []driver.Value) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("CONN(%d) ► Query(%s) → ROWS(%d)", connID, query, rowsID)
}

// ConnQueryContext satisfies Logger interface
func (dl *DefaultSQLLogger) ConnQueryContext(ctx context.Context, connID int64, rowsID int64, query string, args []driver.NamedValue) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("CONN(%d) ► Query(%s) → ROWS(%d)", connID, query, rowsID)
}

// ConnExec satisfies Logger interface
func (dl *DefaultSQLLogger) ConnExec(connID int64, query string, args []driver.Value) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("CONN(%d) ► Exec(%s)", connID, query)
}

// ConnExecContext satisfies Logger interface
func (dl *DefaultSQLLogger) ConnExecContext(ctx context.Context, connID int64, query string, args []driver.NamedValue) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("CONN(%d) ► Exec(%s)", connID, query)
}

// ConnClose satisfies Logger interface
func (dl *DefaultSQLLogger) ConnClose(connID int64) {
	if !dl.Enabled {
		return
	}
	if dl.LogClose {
		dl.log.Printf("CONN(%d) ► Close", connID)
	}
}

// StmtExec satisfies Logger interface
func (dl *DefaultSQLLogger) StmtExec(stmtID int64, query string, args []driver.Value) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("STMT(%d) ► Exec(%s)", stmtID, query)
}

// StmtExecContext satisfies Logger interface
func (dl *DefaultSQLLogger) StmtExecContext(ctx context.Context, stmtID int64, query string, args []driver.NamedValue) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("STMT(%d) ► Exec(%s)", stmtID, query)
}

// StmtQuery satisfies Logger interface
func (dl *DefaultSQLLogger) StmtQuery(stmtID, rowsID int64, query string, args []driver.Value) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("STMT(%d) ► Query(%s) → ROWS(%d)", stmtID, query, rowsID)
}

// StmtQueryContext satisfies Logger interface
func (dl *DefaultSQLLogger) StmtQueryContext(ctx context.Context, stmtID int64, rowsID int64, query string, args []driver.NamedValue) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("STMT(%d) ► Query(%s) → ROWS(%d)", stmtID, query, rowsID)
}

// StmtClose satisfies Logger interface
func (dl *DefaultSQLLogger) StmtClose(stmtID int64) {
	if !dl.Enabled {
		return
	}
	if dl.LogClose {
		dl.log.Printf("STMT(%d) ► Close", stmtID)
	}
}
