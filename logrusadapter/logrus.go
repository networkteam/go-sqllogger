package logrusadapter

import (
	"database/sql/driver"

	"github.com/networkteam/go-sqllogger"
	"github.com/sirupsen/logrus"
)

func NewLogger(l *logrus.Logger, opts ...Opts) *Logger {
	var o Opts
	switch len(opts) {
	case 0:
		o = DefaultOpts()
	case 1:
		o = opts[0]
	default:
		panic("NewLogger expects zero or one opts")
	}
	return &Logger{
		logrusLogger: l,
		opts:         o,
	}
}

type Logger struct {
	logrusLogger *logrus.Logger

	opts Opts
}

var _ sqllogger.Logger = Logger{}

type Opts struct {
	ConnectLevel logrus.Level
	PrepareLevel logrus.Level
	QueryLevel   logrus.Level
	ExecLevel    logrus.Level
	CloseLevel   logrus.Level
	TxLevel      logrus.Level
}

func DefaultOpts() Opts {
	return Opts{
		ConnectLevel: logrus.DebugLevel,
		PrepareLevel: logrus.DebugLevel,
		QueryLevel:   logrus.InfoLevel,
		ExecLevel:    logrus.InfoLevel,
		CloseLevel:   logrus.DebugLevel,
		TxLevel:      logrus.InfoLevel,
	}
}

func (l Logger) Connect(connID int64) {
	l.logrusLogger.
		WithField("connID", connID).
		Log(l.opts.ConnectLevel, "DB Connect")
}

func (l Logger) ConnBegin(connID, txID int64, opts driver.TxOptions) {
	l.logrusLogger.
		WithField("connID", connID).
		Log(l.opts.TxLevel, "CONN Begin")
}

func (l Logger) ConnPrepare(connID, stmtID int64, query string) {
	l.logrusLogger.
		WithField("connID", connID).
		WithField("query", query).
		WithField("stmtID", stmtID).
		Log(l.opts.PrepareLevel, "CONN Prepare")
}

func (l Logger) ConnPrepareContext(connID int64, stmtID int64, query string) {
	l.logrusLogger.
		WithField("connID", connID).
		WithField("query", query).
		WithField("stmtID", stmtID).
		Log(l.opts.PrepareLevel, "CONN Prepare")
}

func (l Logger) ConnQuery(connID, rowsID int64, query string, args []driver.Value) {
	l.logrusLogger.
		WithField("connID", connID).
		WithField("query", query).
		WithField("args", args).
		WithField("rowsID", rowsID).
		Log(l.opts.QueryLevel, "CONN Query")
}

func (l Logger) ConnQueryContext(connID int64, rowsID int64, query string, args []driver.NamedValue) {
	l.logrusLogger.
		WithField("connID", connID).
		WithField("query", query).
		WithField("args", args).
		WithField("rowsID", rowsID).
		Log(l.opts.QueryLevel, "CONN Query")
}

func (l Logger) ConnExec(connID int64, query string, args []driver.Value) {
	l.logrusLogger.
		WithField("connID", connID).
		WithField("query", query).
		Log(l.opts.ExecLevel, "CONN Exec")
}

func (l Logger) ConnExecContext(connID int64, query string, args []driver.NamedValue) {
	l.logrusLogger.
		WithField("connID", connID).
		WithField("query", query).
		WithField("args", args).
		Log(l.opts.ExecLevel, "CONN Exec")
}

func (l Logger) ConnClose(connID int64) {
	l.logrusLogger.
		WithField("connID", connID).
		Log(l.opts.CloseLevel, "CONN Close")
}

func (l Logger) StmtExec(stmtID int64, query string, args []driver.Value) {
	l.logrusLogger.
		WithField("stmtID", stmtID).
		WithField("query", query).
		WithField("args", args).
		Log(l.opts.ExecLevel, "STMT Exec")
}

func (l Logger) StmtExecContext(stmtID int64, query string, args []driver.NamedValue) {
	l.logrusLogger.
		WithField("stmtID", stmtID).
		WithField("query", query).
		WithField("args", args).
		Log(l.opts.ExecLevel, "STMT Exec")
}

func (l Logger) StmtQuery(stmtID int64, rowsID int64, query string, args []driver.Value) {
	l.logrusLogger.
		WithField("stmtID", stmtID).
		WithField("query", query).
		WithField("args", args).
		WithField("rowsID", rowsID).
		Log(l.opts.QueryLevel, "STMT Query")
}

func (l Logger) StmtQueryContext(stmtID int64, rowsID int64, query string, args []driver.NamedValue) {
	l.logrusLogger.
		WithField("stmtID", stmtID).
		WithField("query", query).
		WithField("args", args).
		WithField("rowsID", rowsID).
		Log(l.opts.QueryLevel, "STMT Query")
}

func (l Logger) StmtClose(stmtID int64) {
	l.logrusLogger.
		WithField("stmtID", stmtID).
		Log(l.opts.CloseLevel, "STMT Close")
}

func (l Logger) RowsClose(rowsID int64) {
	l.logrusLogger.
		WithField("rowsID", rowsID).
		Log(l.opts.CloseLevel, "ROWS Close")
}

func (l Logger) TxCommit(txID int64) {
	l.logrusLogger.
		WithField("txID", txID).
		Log(l.opts.TxLevel, "TX Commit")
}

func (l Logger) TxRollback(txID int64) {
	l.logrusLogger.
		WithField("txID", txID).
		Log(l.opts.TxLevel, "TX Rollback")
}
