package sqllogger

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"sync/atomic"
)

// LoggingConnector wraps the given driver.Connector
// and invokes the given SQLLogger for queries and other SQL operations.
//
// Note: Due to the amount of optional interfaces in the database/sql/driver package, there might be some features
// of the original driver that are not exposed on the returned driver.Connector.
func LoggingConnector(log SQLLogger, connector driver.Connector) driver.Connector {
	return &lconnector{
		log:   log,
		cnct:  connector,
		idseq: &idseq{},
	}
}

type lconnector struct {
	cnct  driver.Connector
	log   SQLLogger
	idseq *idseq
}

func (l *lconnector) Connect(ctx context.Context) (driver.Conn, error) {
	originalConn, err := l.cnct.Connect(ctx)
	if err != nil {
		return nil, err
	}

	id := l.idseq.next()
	if logger, ok := l.log.(ConnectLogger); ok {
		logger.Connect(id)
	}
	return &lconn{id: id, log: l.log, conn: originalConn, idseq: l.idseq}, nil
}

func (l *lconnector) Driver() driver.Driver {
	origDriver := l.cnct.Driver()
	return &ld{log: l.log, drv: origDriver}
}

type lconn struct {
	id    int64
	log   SQLLogger
	conn  driver.Conn
	idseq *idseq
}

func (l *lconn) Begin() (driver.Tx, error) {
	origTx, err := l.conn.Begin()
	if err != nil {
		return nil, err
	}

	txID := l.idseq.next()
	if logger, ok := l.log.(ConnBeginLogger); ok {
		logger.ConnBegin(l.id, txID, driver.TxOptions{})
	}

	return l.wrapTx(txID, origTx), nil
}

func (l *lconn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	if connBeginTx, ok := l.conn.(driver.ConnBeginTx); ok {
		origTx, err := connBeginTx.BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}

		txID := l.idseq.next()
		if logger, ok := l.log.(ConnBeginLogger); ok {
			logger.ConnBegin(l.id, txID, opts)
		}

		return l.wrapTx(txID, origTx), nil
	}

	// Copied from driver.go to check for non-default opts if ConnBeginTx interface is not implemented by driver

	// Check the transaction level. If the transaction level is non-default
	// then return an error here as the BeginTx driver value is not supported.
	if opts.Isolation != driver.IsolationLevel(sql.LevelDefault) {
		return nil, errors.New("sql: driver does not support non-default isolation level")
	}

	// If a read-only transaction is requested return an error as the
	// BeginTx driver value is not supported.
	if opts.ReadOnly {
		return nil, errors.New("sql: driver does not support read-only transactions")
	}

	origTx, err := l.conn.Begin()
	if err != nil {
		return nil, err
	}

	txID := l.idseq.next()
	if logger, ok := l.log.(ConnBeginLogger); ok {
		logger.ConnBegin(l.id, txID, opts)
	}

	return l.wrapTx(txID, origTx), nil
}

func (l *lconn) Query(query string, args []driver.Value) (rows driver.Rows, err error) {
	var rowsID int64
	if logger, ok := l.log.(ConnQueryTracer); ok {
		defer logger.ConnQueryTrace(l.id, query, args)(rowsID, err)
	}
	if queryer, ok := l.conn.(driver.Queryer); ok {
		var origRows driver.Rows
		origRows, err = queryer.Query(query, args)
		if err != nil {
			return nil, err
		}

		rowsID = l.idseq.next()
		if logger, ok := l.log.(ConnQueryLogger); ok {
			logger.ConnQuery(l.id, rowsID, query, args)
		}

		return wrapRows(rowsID, l.log, origRows), nil
	}
	return nil, driver.ErrSkip
}

func (l *lconn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (rows driver.Rows, err error) {
	if queryerCtx, ok := l.conn.(driver.QueryerContext); ok {
		var rowsID int64
		if logger, ok := l.log.(ConnQueryContextTracer); ok {
			defer logger.ConnQueryContextTrace(ctx, l.id, query, args)(rowsID, err)
		}

		var origRows driver.Rows
		origRows, err = queryerCtx.QueryContext(ctx, query, args)
		if err != nil {
			return nil, err
		}

		rowsID = l.idseq.next()
		if logger, ok := l.log.(ConnQueryContextLogger); ok {
			logger.ConnQueryContext(l.id, rowsID, query, args)
		}

		return wrapRows(rowsID, l.log, origRows), nil
	}
	return nil, driver.ErrSkip
}

func (l *lconn) Exec(query string, args []driver.Value) (driver.Result, error) {
	if execer, ok := l.conn.(driver.Execer); ok {
		res, err := execer.Exec(query, args)
		if err != nil {
			return nil, err
		}

		if logger, ok := l.log.(ConnExecLogger); ok {
			logger.ConnExec(l.id, query, args)
		}

		return res, nil
	}
	return nil, driver.ErrSkip
}

func (l *lconn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if execerCtx, ok := l.conn.(driver.ExecerContext); ok {
		res, err := execerCtx.ExecContext(ctx, query, args)
		if err != nil {
			return nil, err
		}

		if logger, ok := l.log.(ConnExecContextLogger); ok {
			logger.ConnExecContext(l.id, query, args)
		}

		return res, nil
	}
	return nil, driver.ErrSkip
}

func (l *lconn) Prepare(query string) (driver.Stmt, error) {
	origStmt, err := l.conn.Prepare(query)
	if err != nil {
		return nil, err
	}

	stmtID := l.idseq.next()
	if logger, ok := l.log.(ConnPrepareLogger); ok {
		logger.ConnPrepare(l.id, stmtID, query)
	}

	return &lstmt{id: stmtID, log: l.log, stmt: origStmt, query: query, idseq: l.idseq}, nil
}

func (l *lconn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	if connPrepareCtx, ok := l.conn.(driver.ConnPrepareContext); ok {
		origStmt, err := connPrepareCtx.PrepareContext(ctx, query)
		if err != nil {
			return nil, err
		}

		stmtID := l.idseq.next()
		if logger, ok := l.log.(ConnPrepareContextLogger); ok {
			logger.ConnPrepareContext(l.id, stmtID, query)
		}

		return &lstmt{id: stmtID, log: l.log, stmt: origStmt, query: query, idseq: l.idseq}, nil
	}

	// Copied from ctxutil.go to handle fallback if interface is not implemented

	si, err := l.Prepare(query)
	if err == nil {
		select {
		default:
		case <-ctx.Done():
			si.Close()
			return nil, ctx.Err()
		}
	}
	return si, err
}

func (l *lconn) Close() error {
	if logger, ok := l.log.(ConnCloseLogger); ok {
		logger.ConnClose(l.id)
	}
	return l.conn.Close()
}

var _ driver.Execer = &lconn{}
var _ driver.ExecerContext = &lconn{}
var _ driver.Queryer = &lconn{}
var _ driver.QueryerContext = &lconn{}
var _ driver.ConnBeginTx = &lconn{}
var _ driver.ConnPrepareContext = &lconn{}

type lstmt struct {
	log   SQLLogger
	stmt  driver.Stmt
	query string
	id    int64
	idseq *idseq
}

func (l *lstmt) Close() error {
	err := l.stmt.Close()
	if err != nil {
		return err
	}

	if logger, ok := l.log.(StmtCloseLogger); ok {
		logger.StmtClose(l.id)
	}

	return nil
}

func (l *lstmt) NumInput() int {
	return l.stmt.NumInput()
}

func (l *lstmt) Exec(args []driver.Value) (driver.Result, error) {
	res, err := l.stmt.Exec(args)
	if err != nil {
		return nil, err
	}

	if logger, ok := l.log.(StmtExecLogger); ok {
		logger.StmtExec(l.id, l.query, args)
	}

	return res, err
}

func (l *lstmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	if stmtExecCtx, ok := l.stmt.(driver.StmtExecContext); ok {
		res, err := stmtExecCtx.ExecContext(ctx, args)
		if err != nil {
			return nil, err
		}

		if logger, ok := l.log.(StmtExecContextLogger); ok {
			logger.StmtExecContext(l.id, l.query, args)
		}

		return res, nil
	}

	// Copied from ctxutil.go for fallback handling if driver does not implement StmtExecContext

	dargs, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}

	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	return l.Exec(dargs)
}

func (l *lstmt) Query(args []driver.Value) (rows driver.Rows, err error) {
	var rowsID int64
	if logger, ok := l.log.(StmtQueryTracer); ok {
		defer logger.StmtQueryTrace(l.id, l.query, args)(rowsID, err)
	}
	var origRows driver.Rows
	origRows, err = l.stmt.Query(args)
	if err != nil {
		return nil, err
	}

	rowsID = l.idseq.next()
	if logger, ok := l.log.(StmtQueryLogger); ok {
		logger.StmtQuery(l.id, rowsID, l.query, args)
	}

	return wrapRows(rowsID, l.log, origRows), nil
}

func (l *lstmt) QueryContext(ctx context.Context, args []driver.NamedValue) (rows driver.Rows, err error) {
	if stmtQueryCtx, ok := l.stmt.(driver.StmtQueryContext); ok {
		var rowsID int64
		if logger, ok := l.log.(StmtQueryContextTracer); ok {
			defer logger.StmtQueryContextTrace(ctx, l.id, l.query, args)(rowsID, err)
		}

		rows, err = stmtQueryCtx.QueryContext(ctx, args)
		if err != nil {
			return nil, err
		}

		rowsID = l.idseq.next()
		if logger, ok := l.log.(StmtQueryContextLogger); ok {
			logger.StmtQueryContext(l.id, rowsID, l.query, args)
		}

		return wrapRows(rowsID, l.log, rows), nil
	}

	// Copied from ctxutil.go for fallback handling if driver does not implement StmtQueryContext

	dargs, err := namedValueToValue(args)
	if err != nil {
		return nil, err
	}

	select {
	default:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	return l.Query(dargs)
}

var _ driver.StmtExecContext = &lstmt{}
var _ driver.StmtQueryContext = &lstmt{}

type lrows struct {
	log  SQLLogger
	rows driver.Rows
	id   int64
}

var _ driver.RowsNextResultSet = &lrows{}

func (l *lrows) HasNextResultSet() bool {
	if nrsRows, ok := l.rows.(driver.RowsNextResultSet); ok {
		return nrsRows.HasNextResultSet()
	}
	return false
}

func (l *lrows) NextResultSet() error {
	if nrsRows, ok := l.rows.(driver.RowsNextResultSet); ok {
		return nrsRows.NextResultSet()
	}
	return io.EOF
}

func (l *lrows) Columns() []string {
	return l.rows.Columns()
}

func (l *lrows) Close() error {
	err := l.rows.Close()
	if err != nil {
		return err
	}

	if logger, ok := l.log.(RowsCloseLogger); ok {
		logger.RowsClose(l.id)
	}

	return nil
}

func (l *lrows) Next(dest []driver.Value) error {
	return l.rows.Next(dest)
}

func wrapRows(id int64, log SQLLogger, rows driver.Rows) driver.Rows {
	return &lrows{
		id:   id,
		log:  log,
		rows: rows,
	}
}

type ltx struct {
	log SQLLogger
	tx  driver.Tx
	id  int64
}

func (l *ltx) Commit() error {
	err := l.tx.Commit()
	if err != nil {
		return err
	}

	if logger, ok := l.log.(TxCommitLogger); ok {
		logger.TxCommit(l.id)
	}

	return nil
}

func (l *ltx) Rollback() error {
	err := l.tx.Rollback()
	if err != nil {
		return err
	}

	if logger, ok := l.log.(TxRollbackLogger); ok {
		logger.TxRollback(l.id)
	}

	return nil
}

func (l *lconn) wrapTx(id int64, tx driver.Tx) driver.Tx {
	return &ltx{
		id:  id,
		log: l.log,
		tx:  tx,
	}
}

type ld struct {
	log SQLLogger
	drv driver.Driver
}

func (l *ld) Open(name string) (driver.Conn, error) {
	panic("Not implemented, use sql.OpenDB(...)")
}

func namedValueToValue(named []driver.NamedValue) ([]driver.Value, error) {
	dargs := make([]driver.Value, len(named))
	for n, param := range named {
		if len(param.Name) > 0 {
			return nil, errors.New("sql: driver does not support the use of Named Parameters")
		}
		dargs[n] = param.Value
	}
	return dargs, nil
}

type idseq struct {
	id atomic.Int64
}

func (i *idseq) next() int64 {
	return i.id.Add(1)
}
