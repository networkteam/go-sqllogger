package sqllogger_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/networkteam/go-sqllogger"
)

func TestLoggingConnector(t *testing.T) {
	logger := newTestLogger()
	connector := new(fakeConnector)
	loggingConnector := sqllogger.LoggingConnector(logger, connector)

	ctx := context.Background()

	db := sql.OpenDB(loggingConnector)
	stmt, err := db.PrepareContext(ctx, "CREATE|fizzbuzz|seq=int16,fizz=bool,buzz=bool")
	if err != nil {
		t.Fatalf("Unexpected error from PrepareContext: %v", err)
	}

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		t.Fatalf("Unexpected error from ExecContext: %v", err)
	}

	actualLogs := len(logger.logs)
	expectedLogs := 3
	if actualLogs != expectedLogs {
		t.Fatalf("Exepcted %d log entries, got %d", expectedLogs, actualLogs)
	}
}

type testLogger struct {
	logs []string
}

var _ sqllogger.SQLLogger = &testLogger{}

func (tl *testLogger) Connect(connID int64) {
	tl.logs = append(tl.logs, "Connect")
}

func (tl *testLogger) ConnBegin(connID, txID int64, opts driver.TxOptions) {
	panic("implement me")
}

func (tl *testLogger) ConnPrepare(connID, stmtID int64, query string) {
	panic("implement me")
}

func (tl *testLogger) ConnPrepareContext(connID int64, stmtID int64, query string) {
	tl.logs = append(tl.logs, "ConnPrepareContext")
}

func (tl *testLogger) ConnQuery(connID, rowsID int64, query string, args []driver.Value) {
	panic("implement me")
}

func (tl *testLogger) ConnQueryContext(connID int64, rowsID int64, query string, args []driver.NamedValue) {
	panic("implement me")
}

func (tl *testLogger) ConnExec(connID int64, query string, args []driver.Value) {
	panic("implement me")
}

func (tl *testLogger) ConnExecContext(connID int64, query string, args []driver.NamedValue) {
	panic("implement me")
}

func (tl *testLogger) ConnClose(connID int64) {
	panic("implement me")
}

func (tl *testLogger) StmtExec(stmtID int64, query string, args []driver.Value) {
	panic("implement me")
}

func (tl *testLogger) StmtExecContext(stmtID int64, query string, args []driver.NamedValue) {
	tl.logs = append(tl.logs, "StmtExecContext")
}

func (tl *testLogger) StmtQuery(stmtID int64, rowsID int64, query string, args []driver.Value) {
	panic("implement me")
}

func (tl *testLogger) StmtQueryContext(stmtID int64, rowsID int64, query string, args []driver.NamedValue) {
	panic("implement me")
}

func (tl *testLogger) StmtClose(stmtID int64) {
	panic("implement me")
}

func (tl *testLogger) RowsClose(rowsID int64) {
	panic("implement me")
}

func (tl *testLogger) TxCommit(txID int64) {
	panic("implement me")
}

func (tl *testLogger) TxRollback(txID int64) {
	panic("implement me")
}

func newTestLogger() *testLogger {
	return &testLogger{}
}
