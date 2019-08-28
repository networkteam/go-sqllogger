package sqllogger

import (
	"database/sql/driver"
)

// StdLogger is an interface to adapt the DefaultLogger to the standard library log.Logger or other log frameworks
type StdLogger interface {
	Printf(format string, args ...interface{})
}

// NewDefaultLogger creates a new default logger with sensible defaults
func NewDefaultLogger(log StdLogger) *DefaultLogger {
	return &DefaultLogger{
		log:        log,
		Enabled:    true,
		LogConnect: true,
		LogClose:   false,
	}
}

// DefaultLogger is an implementation of the Logger interface logging to a *log.Logger from the standard library
type DefaultLogger struct {
	log StdLogger

	// Enabled sets, whether
	Enabled bool

	LogConnect bool
	LogClose   bool
}

var _ Logger = &DefaultLogger{}

// TxRollback satisfies Logger interface
func (dl *DefaultLogger) TxRollback(txID int64) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("  TX(%d) ► Rollback", txID)
}

// TxCommit satisfies Logger interface
func (dl *DefaultLogger) TxCommit(txID int64) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("  TX(%d) ► Commit", txID)
}

// RowsClose satisfies Logger interface
func (dl *DefaultLogger) RowsClose(rowsID int64) {
	if !dl.Enabled {
		return
	}
	if dl.LogClose {
		dl.log.Printf("ROWS(%d) ► Close", rowsID)
	}
}

// Connect satisfies Logger interface
func (dl *DefaultLogger) Connect(connID int64) {
	if !dl.Enabled {
		return
	}
	if !dl.LogConnect {
		return
	}
	dl.log.Printf("Connect → CONN(%d)", connID)
}

// ConnBegin satisfies Logger interface
func (dl *DefaultLogger) ConnBegin(connID, txID int64, opts driver.TxOptions) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("CONN(%d) ► Begin -> TX(%d)", connID, txID)
}

// ConnPrepare satisfies Logger interface
func (dl *DefaultLogger) ConnPrepare(connID, stmtID int64, query string) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("CONN(%d) ► Prepare(%s) → STMT(%d)", connID, query, stmtID)
}

// ConnPrepareContext satisfies Logger interface
func (dl *DefaultLogger) ConnPrepareContext(connID int64, stmtID int64, query string) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("CONN(%d) ► Prepare(%s) → STMT(%d)", connID, query, stmtID)
}

// ConnQuery satisfies Logger interface
func (dl *DefaultLogger) ConnQuery(connID, rowsID int64, query string, args []driver.Value) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("CONN(%d) ► Query(%s) → ROWS(%d)", connID, query, rowsID)
}

// ConnQueryContext satisfies Logger interface
func (dl *DefaultLogger) ConnQueryContext(connID int64, rowsID int64, query string, args []driver.NamedValue) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("CONN(%d) ► Query(%s) → ROWS(%d)", connID, query, rowsID)
}

// ConnExec satisfies Logger interface
func (dl *DefaultLogger) ConnExec(connID int64, query string, args []driver.Value) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("CONN(%d) ► Exec(%s)", connID, query)
}

// ConnExecContext satisfies Logger interface
func (dl *DefaultLogger) ConnExecContext(connID int64, query string, args []driver.NamedValue) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("CONN(%d) ► Exec(%s)", connID, query)
}

// ConnClose satisfies Logger interface
func (dl *DefaultLogger) ConnClose(connID int64) {
	if !dl.Enabled {
		return
	}
	if dl.LogClose {
		dl.log.Printf("CONN(%d) ► Close", connID)
	}
}

// StmtExec satisfies Logger interface
func (dl *DefaultLogger) StmtExec(stmtID int64, query string, args []driver.Value) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("STMT(%d) ► Exec(%s)", stmtID, query)
}

// StmtExecContext satisfies Logger interface
func (dl *DefaultLogger) StmtExecContext(stmtID int64, query string, args []driver.NamedValue) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("STMT(%d) ► Exec(%s)", stmtID, query)
}

// StmtQuery satisfies Logger interface
func (dl *DefaultLogger) StmtQuery(stmtID, rowsID int64, query string, args []driver.Value) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("STMT(%d) ► Query(%s) → ROWS(%d)", stmtID, query, rowsID)
}

// StmtQueryContext satisfies Logger interface
func (dl *DefaultLogger) StmtQueryContext(stmtID int64, rowsID int64, query string, args []driver.NamedValue) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("STMT(%d) ► Query(%s) → ROWS(%d)", stmtID, query, rowsID)
}

// StmtClose satisfies Logger interface
func (dl *DefaultLogger) StmtClose(stmtID int64) {
	if !dl.Enabled {
		return
	}
	if dl.LogClose {
		dl.log.Printf("STMT(%d) ► Close", stmtID)
	}
}
