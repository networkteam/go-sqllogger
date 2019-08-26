package sqllogger

import (
	"database/sql/driver"
	"log"
)

func NewDefaultLogger(log *log.Logger) Logger {
	return &DefaultLogger{
		log: log,
	}
}

type DefaultLogger struct {
	log *log.Logger
}

func (dl *DefaultLogger) TxRollback(txID int64) {
	dl.log.Printf("  TX(%d) ► Rollback", txID)
}

func (dl *DefaultLogger) TxCommit(txID int64) {
	dl.log.Printf("  TX(%d) ► Commit", txID)
}

func (dl *DefaultLogger) RowsClose(rowsID int64) {
	dl.log.Printf("ROWS(%d) ► Close", rowsID)
}

func (dl *DefaultLogger) Connect(connID int64) {
	dl.log.Printf("Connect → CONN(%d)", connID)
}

func (dl *DefaultLogger) ConnBegin(connID, txID int64, opts driver.TxOptions) {
	dl.log.Printf("CONN(%d) ► Begin -> TX(%d)", connID, txID)
}

func (dl *DefaultLogger) ConnPrepare(connID, stmtID int64, query string) {
	dl.log.Printf("CONN(%d) ► Prepare(%s) → STMT(%d)", connID, query, stmtID)
}

func (dl *DefaultLogger) ConnPrepareContext(connID int64, stmtID int64, query string) {
	dl.log.Printf("CONN(%d) ► Prepare(%s) → STMT(%d)", connID, query, stmtID)
}

func (dl *DefaultLogger) ConnQuery(connID, rowsID int64, query string, args []driver.Value) {
	dl.log.Printf("CONN(%d) ► Query(%s) → ROWS(%d)", connID, query, rowsID)
}

func (dl *DefaultLogger) ConnQueryContext(connID int64, rowsID int64, query string, args []driver.NamedValue) {
	dl.log.Printf("CONN(%d) ► Query(%s) → ROWS(%d)", connID, query, rowsID)
}

func (dl *DefaultLogger) ConnExec(connID int64, query string, args []driver.Value) {
	dl.log.Printf("CONN(%d) ► Exec(%s)", connID, query)
}

func (dl *DefaultLogger) ConnExecContext(connID int64, query string, args []driver.NamedValue) {
	dl.log.Printf("CONN(%d) ► Exec(%s)", connID, query)
}

func (dl *DefaultLogger) ConnClose(connID int64) {
	dl.log.Printf("CONN(%d) ► Close", connID)
}

func (dl *DefaultLogger) StmtExec(stmtID int64, query string, args []driver.Value) {
	dl.log.Printf("STMT(%d) ► Exec(%s)", stmtID, query)
}

func (dl *DefaultLogger) StmtExecContext(stmtID int64, query string, args []driver.NamedValue) {
	dl.log.Printf("STMT(%d) ► Exec(%s)", stmtID, query)
}

func (dl *DefaultLogger) StmtQuery(stmtID, rowsID int64, query string, args []driver.Value) {
	dl.log.Printf("STMT(%d) ► Query(%s) → ROWS(%d)", stmtID, query, rowsID)
}

func (dl *DefaultLogger) StmtQueryContext(stmtID int64, rowsID int64, query string, args []driver.NamedValue) {
	dl.log.Printf("STMT(%d) ► Query(%s) → ROWS(%d)", stmtID, query, rowsID)
}

func (dl *DefaultLogger) StmtClose(stmtID int64) {
	dl.log.Printf("STMT(%d) ► Close", stmtID)
}
