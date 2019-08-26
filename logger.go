package sqllogger

import (
	"database/sql/driver"
)

type Logger interface {
	Connect(connID int64)

	ConnBegin(connID, txID int64, opts driver.TxOptions)
	ConnPrepare(connID, stmtID int64, query string)
	ConnPrepareContext(connID int64, stmtID int64, query string)
	ConnQuery(connID, rowsID int64, query string, args []driver.Value)
	ConnQueryContext(connID int64, rowsID int64, query string, args []driver.NamedValue)
	ConnExec(connID int64, query string, args []driver.Value)
	ConnExecContext(connID int64, query string, args []driver.NamedValue)
	ConnClose(connID int64)

	StmtExec(stmtID int64, query string, args []driver.Value)
	StmtExecContext(stmtID int64, query string, args []driver.NamedValue)
	StmtQuery(stmtID int64, rowsID int64, query string, args []driver.Value)
	StmtQueryContext(stmtID int64, rowsID int64, query string, args []driver.NamedValue)
	StmtClose(stmtID int64)

	RowsClose(rowsID int64)

	TxCommit(txID int64)
	TxRollback(txID int64)
}
