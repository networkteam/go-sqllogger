package sqllogger

import (
	"database/sql/driver"
	"log"
)

func NewDefaultLogger(log *log.Logger) *DefaultLogger {
	return &DefaultLogger{
		log:        log,
		Enabled:    true,
		LogConnect: true,
		LogClose:   false,
	}
}

type DefaultLogger struct {
	log *log.Logger

	Enabled bool

	LogConnect bool
	LogClose   bool
}

var _ Logger = &DefaultLogger{}

func (dl *DefaultLogger) TxRollback(txID int64) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("  TX(%d) ► Rollback", txID)
}

func (dl *DefaultLogger) TxCommit(txID int64) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("  TX(%d) ► Commit", txID)
}

func (dl *DefaultLogger) RowsClose(rowsID int64) {
	if !dl.Enabled {
		return
	}
	if dl.LogClose {
		dl.log.Printf("ROWS(%d) ► Close", rowsID)
	}
}

func (dl *DefaultLogger) Connect(connID int64) {
	if !dl.Enabled {
		return
	}
	if !dl.LogConnect {
		return
	}
	dl.log.Printf("Connect → CONN(%d)", connID)
}

func (dl *DefaultLogger) ConnBegin(connID, txID int64, opts driver.TxOptions) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("CONN(%d) ► Begin -> TX(%d)", connID, txID)
}

func (dl *DefaultLogger) ConnPrepare(connID, stmtID int64, query string) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("CONN(%d) ► Prepare(%s) → STMT(%d)", connID, query, stmtID)
}

func (dl *DefaultLogger) ConnPrepareContext(connID int64, stmtID int64, query string) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("CONN(%d) ► Prepare(%s) → STMT(%d)", connID, query, stmtID)
}

func (dl *DefaultLogger) ConnQuery(connID, rowsID int64, query string, args []driver.Value) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("CONN(%d) ► Query(%s) → ROWS(%d)", connID, query, rowsID)
}

func (dl *DefaultLogger) ConnQueryContext(connID int64, rowsID int64, query string, args []driver.NamedValue) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("CONN(%d) ► Query(%s) → ROWS(%d)", connID, query, rowsID)
}

func (dl *DefaultLogger) ConnExec(connID int64, query string, args []driver.Value) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("CONN(%d) ► Exec(%s)", connID, query)
}

func (dl *DefaultLogger) ConnExecContext(connID int64, query string, args []driver.NamedValue) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("CONN(%d) ► Exec(%s)", connID, query)
}

func (dl *DefaultLogger) ConnClose(connID int64) {
	if !dl.Enabled {
		return
	}
	if dl.LogClose {
		dl.log.Printf("CONN(%d) ► Close", connID)
	}
}

func (dl *DefaultLogger) StmtExec(stmtID int64, query string, args []driver.Value) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("STMT(%d) ► Exec(%s)", stmtID, query)
}

func (dl *DefaultLogger) StmtExecContext(stmtID int64, query string, args []driver.NamedValue) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("STMT(%d) ► Exec(%s)", stmtID, query)
}

func (dl *DefaultLogger) StmtQuery(stmtID, rowsID int64, query string, args []driver.Value) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("STMT(%d) ► Query(%s) → ROWS(%d)", stmtID, query, rowsID)
}

func (dl *DefaultLogger) StmtQueryContext(stmtID int64, rowsID int64, query string, args []driver.NamedValue) {
	if !dl.Enabled {
		return
	}
	dl.log.Printf("STMT(%d) ► Query(%s) → ROWS(%d)", stmtID, query, rowsID)
}

func (dl *DefaultLogger) StmtClose(stmtID int64) {
	if !dl.Enabled {
		return
	}
	if dl.LogClose {
		dl.log.Printf("STMT(%d) ► Close", stmtID)
	}
}
