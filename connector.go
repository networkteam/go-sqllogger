package sqllogger

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"reflect"
	"sync"
)

// LoggingConnector wraps the given driver.Connector
// and invokes the given SQLLogger for queries and other SQL operations.
//
// Note: Due to the amount of optional interfaces in the database/sql/driver package, there might be some features
// of the original driver that are not exposed on the returned driver.Connector.
func LoggingConnector(log SQLLogger, connector driver.Connector) driver.Connector {
	return &lconnector{
		log:  log,
		cnct: connector,
	}
}

var (
	idseq   int64
	idseqMx sync.Mutex
)

type lconnector struct {
	cnct driver.Connector
	log  SQLLogger
}

var _ driver.Connector = &lconnector{}

func (l *lconnector) Connect(ctx context.Context) (driver.Conn, error) {
	originalConn, err := l.cnct.Connect(ctx)
	if err != nil {
		return nil, err
	}

	id := nextID()
	l.log.Connect(id)
	return &lconn{id: id, log: l.log, conn: originalConn}, nil
}

func (l *lconnector) Driver() driver.Driver {
	origDriver := l.cnct.Driver()
	return &ld{log: l.log, drv: origDriver}
}

type lconn struct {
	id   int64
	log  SQLLogger
	conn driver.Conn
}

var _ driver.Conn = &lconn{}
var _ driver.SessionResetter = &lconn{}
var _ driver.Execer = &lconn{}
var _ driver.ExecerContext = &lconn{}
var _ driver.Queryer = &lconn{}
var _ driver.QueryerContext = &lconn{}
var _ driver.ConnBeginTx = &lconn{}
var _ driver.ConnPrepareContext = &lconn{}
var _ driver.NamedValueChecker = &lconn{}
var _ driver.Pinger = &lconn{}
var _ driver.Validator = &lconn{}

func (l *lconn) Begin() (driver.Tx, error) {
	origTx, err := l.conn.Begin()
	if err != nil {
		return nil, err
	}

	txID := nextID()
	l.log.ConnBegin(l.id, txID, driver.TxOptions{})

	return l.wrapTx(txID, origTx), nil
}

func (l *lconn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	if connBeginTx, ok := l.conn.(driver.ConnBeginTx); ok {
		origTx, err := connBeginTx.BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}

		txID := nextID()
		l.log.ConnBegin(l.id, txID, opts)

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

	txID := nextID()
	l.log.ConnBegin(l.id, txID, opts)

	return l.wrapTx(txID, origTx), nil
}

func (l *lconn) Query(query string, args []driver.Value) (driver.Rows, error) {
	if queryer, ok := l.conn.(driver.Queryer); ok {
		origRows, err := queryer.Query(query, args)
		if err != nil {
			return nil, err
		}

		rowsID := nextID()
		l.log.ConnQuery(l.id, rowsID, query, args)

		return wrapRows(rowsID, l.log, origRows), nil
	}
	return nil, driver.ErrSkip
}

func (l *lconn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if queryerCtx, ok := l.conn.(driver.QueryerContext); ok {
		origRows, err := queryerCtx.QueryContext(ctx, query, args)
		if err != nil {
			return nil, err
		}

		rowsID := nextID()
		l.log.ConnQueryContext(l.id, rowsID, query, args)

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

		l.log.ConnExec(l.id, query, args)

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

		l.log.ConnExecContext(l.id, query, args)

		return res, nil
	}
	return nil, driver.ErrSkip
}

func (l *lconn) Prepare(query string) (driver.Stmt, error) {
	origStmt, err := l.conn.Prepare(query)
	if err != nil {
		return nil, err
	}

	stmtID := nextID()
	l.log.ConnPrepare(l.id, stmtID, query)

	return &lstmt{id: stmtID, log: l.log, stmt: origStmt, query: query}, nil
}

func (l *lconn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	if connPrepareCtx, ok := l.conn.(driver.ConnPrepareContext); ok {
		origStmt, err := connPrepareCtx.PrepareContext(ctx, query)
		if err != nil {
			return nil, err
		}

		stmtID := nextID()
		l.log.ConnPrepareContext(l.id, stmtID, query)

		return &lstmt{id: stmtID, log: l.log, stmt: origStmt, query: query}, nil
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

func (l *lconn) ResetSession(ctx context.Context) error {
	if sr, ok := l.conn.(driver.SessionResetter); ok {
		return sr.ResetSession(ctx)
	}
	return nil
}

func (l *lconn) CheckNamedValue(nv *driver.NamedValue) error {
	if checker, ok := l.conn.(driver.NamedValueChecker); ok {
		return checker.CheckNamedValue(nv)
	}
	return driver.ErrSkip
}

func (l *lconn) Ping(ctx context.Context) error {
	if pinger, ok := l.conn.(driver.Pinger); ok {
		return pinger.Ping(ctx)
	}
	return driver.ErrSkip
}

func (l *lconn) IsValid() bool {
	if validator, ok := l.conn.(driver.Validator); ok {
		return validator.IsValid()
	}
	return true // Default to assuming it's valid
}

func (l *lconn) Close() error {
	l.log.ConnClose(l.id)
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
}

var _ driver.Stmt = &lstmt{}
var _ driver.StmtExecContext = &lstmt{}
var _ driver.StmtQueryContext = &lstmt{}
var _ driver.NamedValueChecker = &lstmt{}

func (l *lstmt) Close() error {
	err := l.stmt.Close()
	if err != nil {
		return err
	}

	l.log.StmtClose(l.id)

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

	l.log.StmtExec(l.id, l.query, args)

	return res, err
}

func (l *lstmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	if stmtExecCtx, ok := l.stmt.(driver.StmtExecContext); ok {
		res, err := stmtExecCtx.ExecContext(ctx, args)
		if err != nil {
			return nil, err
		}

		l.log.StmtExecContext(l.id, l.query, args)

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

func (l *lstmt) Query(args []driver.Value) (driver.Rows, error) {
	origRows, err := l.stmt.Query(args)
	if err != nil {
		return nil, err
	}

	rowsID := nextID()
	l.log.StmtQuery(l.id, rowsID, l.query, args)

	return wrapRows(rowsID, l.log, origRows), nil
}

func (l *lstmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	if stmtQueryCtx, ok := l.stmt.(driver.StmtQueryContext); ok {
		rows, err := stmtQueryCtx.QueryContext(ctx, args)
		if err != nil {
			return nil, err
		}

		rowsID := nextID()
		l.log.StmtQueryContext(l.id, rowsID, l.query, args)

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

func (l *lstmt) CheckNamedValue(nv *driver.NamedValue) error {
	if checker, ok := l.stmt.(driver.NamedValueChecker); ok {
		return checker.CheckNamedValue(nv)
	}
	return driver.ErrSkip
}

var _ driver.StmtExecContext = &lstmt{}
var _ driver.StmtQueryContext = &lstmt{}

type lrows struct {
	log  SQLLogger
	rows driver.Rows
	id   int64
}

var _ driver.Rows = &lrows{}
var _ driver.RowsNextResultSet = &lrows{}
var _ driver.RowsColumnTypeDatabaseTypeName = &lrows{}
var _ driver.RowsColumnTypeLength = &lrows{}
var _ driver.RowsColumnTypeNullable = &lrows{}
var _ driver.RowsColumnTypePrecisionScale = &lrows{}
var _ driver.RowsColumnTypeScanType = &lrows{}

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

	l.log.RowsClose(l.id)

	return nil
}

func (l *lrows) Next(dest []driver.Value) error {
	return l.rows.Next(dest)
}

func (l *lrows) ColumnTypeDatabaseTypeName(index int) string {
	if ct, ok := l.rows.(driver.RowsColumnTypeDatabaseTypeName); ok {
		return ct.ColumnTypeDatabaseTypeName(index)
	}
	return ""
}

func (l *lrows) ColumnTypeLength(index int) (length int64, ok bool) {
	if ct, ok := l.rows.(driver.RowsColumnTypeLength); ok {
		return ct.ColumnTypeLength(index)
	}
	return 0, false
}

func (l *lrows) ColumnTypeNullable(index int) (nullable, ok bool) {
	if ct, ok := l.rows.(driver.RowsColumnTypeNullable); ok {
		return ct.ColumnTypeNullable(index)
	}
	return false, false
}

func (l *lrows) ColumnTypePrecisionScale(index int) (precision, scale int64, ok bool) {
	if ct, ok := l.rows.(driver.RowsColumnTypePrecisionScale); ok {
		return ct.ColumnTypePrecisionScale(index)
	}
	return 0, 0, false
}

func (l *lrows) ColumnTypeScanType(index int) reflect.Type {
	if ct, ok := l.rows.(driver.RowsColumnTypeScanType); ok {
		return ct.ColumnTypeScanType(index)
	}
	return nil
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

var _ driver.Tx = &ltx{}

func (l *ltx) Commit() error {
	err := l.tx.Commit()
	if err != nil {
		return err
	}

	l.log.TxCommit(l.id)

	return nil
}

func (l *ltx) Rollback() error {
	err := l.tx.Rollback()
	if err != nil {
		return err
	}

	l.log.TxRollback(l.id)

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

var _ driver.Driver = &ld{}
var _ driver.DriverContext = &ld{}

func (l *ld) Open(name string) (driver.Conn, error) {
	panic("Not implemented, use sql.OpenDB(...)")
}

func (l *ld) OpenConnector(name string) (driver.Connector, error) {
	if dc, ok := l.drv.(driver.DriverContext); ok {
		connector, err := dc.OpenConnector(name)
		if err != nil {
			return nil, err
		}
		return LoggingConnector(l.log, connector), nil
	}
	return nil, driver.ErrSkip
}

func nextID() int64 {
	idseqMx.Lock()
	defer idseqMx.Unlock()

	idseq++
	return idseq
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
