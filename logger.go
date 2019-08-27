package sqllogger

import (
	"database/sql/driver"
)

// Logger is the interface for a log adapter that is used to log SQL queries and operations in the LoggingConnector
//
// All methods are only called if the original operation returned without an error.
type Logger interface {
	// Called on DB connect with a generated connection id
	Connect(connID int64)

	// Called on a transaction begin on a connection with the connection id and a generated transaction id
	ConnBegin(connID, txID int64, opts driver.TxOptions)
	// Called on a prepare statement on a connection with the connection id and a generated statement id
	ConnPrepare(connID, stmtID int64, query string)
	// Called on a prepare statement with context on a connection with the connection id and a generated statement id
	ConnPrepareContext(connID int64, stmtID int64, query string)
	// Called on a query on a connection with the connection id and a generated rows id
	ConnQuery(connID, rowsID int64, query string, args []driver.Value)
	// Called on a query with context on a connection with the connection id and a generated rows id
	ConnQueryContext(connID int64, rowsID int64, query string, args []driver.NamedValue)
	// Called on an exec on a connection with the connection id
	ConnExec(connID int64, query string, args []driver.Value)
	// Called on an exec with context on a connection with the connection id
	ConnExecContext(connID int64, query string, args []driver.NamedValue)
	// Called on a close on a connection with the connection id
	ConnClose(connID int64)

	// Called on an exec on a statement with the statement id
	StmtExec(stmtID int64, query string, args []driver.Value)
	// Called on an exec with context on a statement with the statement id
	StmtExecContext(stmtID int64, query string, args []driver.NamedValue)
	// Called on a query on a statement with the statement id and generated rows id
	StmtQuery(stmtID int64, rowsID int64, query string, args []driver.Value)
	// Called on a query with context on a statement with the statement id and generated rows id
	StmtQueryContext(stmtID int64, rowsID int64, query string, args []driver.NamedValue)
	// Called on a close on a statement with the statement id
	StmtClose(stmtID int64)

	// Called on a close on rows with the rows id
	RowsClose(rowsID int64)

	// Called on a commit on a transaction with the transaction id
	TxCommit(txID int64)
	// Called on a rollback on a transaction with the transaction id
	TxRollback(txID int64)
}
