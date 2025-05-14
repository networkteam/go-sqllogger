package sqllogger

import (
	"context"
	"database/sql/driver"
)

// SQLLogger is the interface for a specialized SQL logger that is used to log SQL queries and operations in the LoggingConnector
//
// With this interface, adapters can be implemented for any log framework. For the standard library log.Logger, a
// DefaultSQLLogger is provided as a default implementation.
//
// All methods are only called if the original operation returned without an error.
type SQLLogger interface {
	// Connect is called on DB connect with a generated connection id.
	Connect(ctx context.Context, connID int64)
	// ConnBegin is called on a transaction begin on a connection with the connection id and a generated transaction id.
	ConnBegin(ctx context.Context, connID, txID int64, opts driver.TxOptions)
	// ConnPrepare is called on a prepare statement on a connection with the connection id and a generated statement id.
	// Note: ctx is only for sqllogger metadata since ConnPrepare does not receive a context.
	ConnPrepare(ctx context.Context, connID, stmtID int64, query string)
	// ConnPrepareContext is called on a prepare statement with context on a connection with the connection id and a generated statement id.
	ConnPrepareContext(ctx context.Context, connID int64, stmtID int64, query string)
	// ConnQuery is called on a query on a connection with the connection id and a generated rows id.
	// Note: ctx is only for sqllogger metadata since ConnQuery does not receive a context.
	ConnQuery(ctx context.Context, connID, rowsID int64, query string, args []driver.Value)
	// ConnQueryContext is called on a query with context on a connection with the connection id and a generated rows id.
	ConnQueryContext(ctx context.Context, connID int64, rowsID int64, query string, args []driver.NamedValue)
	// ConnExec is called on an exec on a connection with the connection id.
	// Note: ctx is only for sqllogger metadata since ConnExec does not receive a context.
	ConnExec(ctx context.Context, connID int64, query string, args []driver.Value)
	// ConnExecContext is called on an exec with context on a connection with the connection id.
	ConnExecContext(ctx context.Context, connID int64, query string, args []driver.NamedValue)
	// ConnClose is called on a close on a connection with the connection id.
	// Note: ctx is only for sqllogger metadata since ConnClose does not receive a context.
	ConnClose(ctx context.Context, connID int64)

	// StmtExec is called on an exec on a statement with the statement id.
	// Note: ctx is only for sqllogger metadata since StmtExec does not receive a context.
	StmtExec(ctx context.Context, stmtID int64, query string, args []driver.Value)
	// StmtExecContext is called on an exec with context on a statement with the statement id.
	StmtExecContext(ctx context.Context, stmtID int64, query string, args []driver.NamedValue)
	// StmtQuery is called on a query on a statement with the statement id and generated rows id.
	// Note: ctx is only for sqllogger metadata since StmtQuery does not receive a context.
	StmtQuery(ctx context.Context, stmtID int64, rowsID int64, query string, args []driver.Value)
	// StmtQueryContext is called on a query with context on a statement with the statement id and generated rows id.
	StmtQueryContext(ctx context.Context, stmtID int64, rowsID int64, query string, args []driver.NamedValue)
	// StmtClose is called on a close on a statement with the statement id.
	// Note: ctx is only for sqllogger metadata since StmtClose does not receive a context.
	StmtClose(ctx context.Context, stmtID int64)

	// RowsClose is called on a close on rows with the rows id.
	// Note: ctx is only for sqllogger metadata since RowsClose does not receive a context.
	RowsClose(ctx context.Context, rowsID int64)

	// TxCommit is called on a commit on a transaction with the transaction id.
	// Note: ctx is only for sqllogger metadata since TxCommit does not receive a context.
	TxCommit(ctx context.Context, txID int64)
	// TxRollback is called on a rollback on a transaction with the transaction id.
	// Note: ctx is only for sqllogger metadata since TxRollback does not receive a context.
	TxRollback(ctx context.Context, txID int64)
}
