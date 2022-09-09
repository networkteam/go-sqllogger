package sqllogger

import (
	"database/sql/driver"
)

// SQLLogger is a marker interface for a specialized SQL logger that is used to log SQL queries and operations in the LoggingConnector.
//
// Implementations can choose to implement the optional *Logger interfaces for logging actual operations.
//
// With this interface, adapters can be implemented for any log framework. For the standard library log.Logger, a
// DefaultSQLLogger is provided as a default implementation implementing FullSQLLogger.
//
// All methods are only called if the original operation returned without an error.
type SQLLogger interface {
}

type FullSQLLogger interface {
	SQLLogger
	ConnectLogger
	ConnBeginLogger
	ConnPrepareLogger
	ConnPrepareContextLogger
	ConnQueryLogger
	ConnExecLogger
	ConnExecContextLogger
	ConnCloseLogger
	StmtExecLogger
	StmtExecContextLogger
	StmtQueryLogger
	StmtQueryContextLogger
	StmtCloseLogger
	RowsCloseLogger
	TxCommitLogger
	TxRollbackLogger
}

type ConnectLogger interface {
	// Connect is called on DB connect with a generated connection id
	Connect(connID int64)
}

type ConnBeginLogger interface {
	ConnBegin(connID, txID int64, opts driver.TxOptions)
}

type ConnPrepareLogger interface {
	// ConnPrepare is called on a prepare statement on a connection with the connection id and a generated statement id
	ConnPrepare(connID, stmtID int64, query string)
}

type ConnPrepareContextLogger interface {
	// ConnPrepareContext is called on a prepare statement with context on a connection with the connection id and a generated statement id
	ConnPrepareContext(connID int64, stmtID int64, query string)
}

type ConnQueryLogger interface {
	// ConnQuery is called on a query on a connection with the connection id and a generated rows id
	ConnQuery(connID, rowsID int64, query string, args []driver.Value)
}

type ConnQueryContextLogger interface {
	// ConnQueryContext is called on a query with context on a connection with the connection id and a generated rows id
	ConnQueryContext(connID int64, rowsID int64, query string, args []driver.NamedValue)
}

type ConnExecLogger interface {
	// ConnExec is called on an exec on a connection with the connection id
	ConnExec(connID int64, query string, args []driver.Value)
}

type ConnExecContextLogger interface {
	// ConnExecContext is called on an exec with context on a connection with the connection id
	ConnExecContext(connID int64, query string, args []driver.NamedValue)
}

type ConnCloseLogger interface {
	// ConnClose is called on a close on a connection with the connection id
	ConnClose(connID int64)
}

type StmtExecLogger interface {
	// StmtExec is called on an exec on a statement with the statement id
	StmtExec(stmtID int64, query string, args []driver.Value)
}

type StmtExecContextLogger interface {
	// StmtExecContext is called on an exec with context on a statement with the statement id
	StmtExecContext(stmtID int64, query string, args []driver.NamedValue)
}

type StmtQueryLogger interface {
	// StmtQuery is called on a query on a statement with the statement id and generated rows id
	StmtQuery(stmtID int64, rowsID int64, query string, args []driver.Value)
}

type StmtQueryContextLogger interface {
	// StmtQueryContext is called on a query with context on a statement with the statement id and generated rows id
	StmtQueryContext(stmtID int64, rowsID int64, query string, args []driver.NamedValue)
}

type StmtCloseLogger interface {
	// StmtClose is called on a close on a statement with the statement id
	StmtClose(stmtID int64)
}

type RowsCloseLogger interface {
	// RowsClose is called on a close on rows with the rows id
	RowsClose(rowsID int64)
}

type TxCommitLogger interface {
	// TxCommit is called on a commit on a transaction with the transaction id
	TxCommit(txID int64)
}

type TxRollbackLogger interface {
	// TxRollback is called on a rollback on a transaction with the transaction id
	TxRollback(txID int64)
}
