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
	stmt, err := db.PrepareContext(ctx, "CREATE|fizzbuzz|seq=int64,fizz=bool,buzz=bool")
	if err != nil {
		t.Fatalf("Unexpected error from PrepareContext: %v", err)
	}

	_, err = stmt.ExecContext(ctx)
	if err != nil {
		t.Fatalf("Unexpected error from ExecContext: %v", err)
	}
	err = stmt.Close()
	if err != nil {
		t.Fatalf("Unexpected error from Close: %v", err)
	}

	_, err = db.ExecContext(ctx, `INSERT|fizzbuzz|seq=?,fizz=?,buzz=?`, 1, false, false)
	if err != nil {
		t.Fatalf("Unexpected error from ExecContext: %v", err)
	}

	actualLogs := logger.logs
	expectedLogs := []string{
		`Connect`,
		`ConnPrepareContext`,
		`StmtExecContext`,
		`StmtClose`,
		`ConnPrepareContext`,
		`StmtExecContext`,
		`StmtClose`,
	}
	if len(actualLogs) != len(expectedLogs) {
		t.Fatalf("Expected %d log entries, got %d: %+v", len(expectedLogs), len(actualLogs), actualLogs)
	}
	for i, actualEntry := range actualLogs {
		if actualEntry != expectedLogs[i] {
			t.Errorf("Expected log entry %d to be %q, got %q", i, expectedLogs[i], actualEntry)
		}
	}
}

type testLogger struct {
	logs []string
}

var _ sqllogger.SQLLogger = &testLogger{}

func (tl *testLogger) Connect(ctx context.Context, connID int64) {
	tl.logs = append(tl.logs, "Connect")
}

func (tl *testLogger) ConnBegin(ctx context.Context, connID, txID int64, opts driver.TxOptions) {
	tl.logs = append(tl.logs, "ConnBegin")
}

func (tl *testLogger) ConnPrepare(connID, stmtID int64, query string) {
	tl.logs = append(tl.logs, "ConnPrepare")
}

func (tl *testLogger) ConnPrepareContext(ctx context.Context, connID int64, stmtID int64, query string) {
	tl.logs = append(tl.logs, "ConnPrepareContext")
}

func (tl *testLogger) ConnQuery(connID, rowsID int64, query string, args []driver.Value) {
	tl.logs = append(tl.logs, "ConnQuery")
}

func (tl *testLogger) ConnQueryContext(ctx context.Context, connID int64, rowsID int64, query string, args []driver.NamedValue) {
	tl.logs = append(tl.logs, "ConnQueryContext")
}

func (tl *testLogger) ConnExec(connID int64, query string, args []driver.Value) {
	tl.logs = append(tl.logs, "ConnExec")
}

func (tl *testLogger) ConnExecContext(ctx context.Context, connID int64, query string, args []driver.NamedValue) {
	tl.logs = append(tl.logs, "ConnExecContext")
}

func (tl *testLogger) ConnClose(connID int64) {
	tl.logs = append(tl.logs, "ConnClose")
}

func (tl *testLogger) StmtExec(stmtID int64, query string, args []driver.Value) {
	tl.logs = append(tl.logs, "StmtExec")
}

func (tl *testLogger) StmtExecContext(ctx context.Context, stmtID int64, query string, args []driver.NamedValue) {
	tl.logs = append(tl.logs, "StmtExecContext")
}

func (tl *testLogger) StmtQuery(stmtID int64, rowsID int64, query string, args []driver.Value) {
	tl.logs = append(tl.logs, "StmtQuery")
}

func (tl *testLogger) StmtQueryContext(ctx context.Context, stmtID int64, rowsID int64, query string, args []driver.NamedValue) {
	tl.logs = append(tl.logs, "StmtQueryContext")
}

func (tl *testLogger) StmtClose(stmtID int64) {
	tl.logs = append(tl.logs, "StmtClose")
}

func (tl *testLogger) RowsClose(rowsID int64) {
	tl.logs = append(tl.logs, "RowsClose")
}

func (tl *testLogger) TxCommit(txID int64) {
	tl.logs = append(tl.logs, "TxCommit")
}

func (tl *testLogger) TxRollback(txID int64) {
	tl.logs = append(tl.logs, "TxRollback")
}

func newTestLogger() *testLogger {
	return &testLogger{}
}
