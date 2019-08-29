package logrusadapter_test

import (
	"bytes"
	"database/sql/driver"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/networkteam/go-sqllogger/logrusadapter"
)

func TestNewSQLLogger(t *testing.T) {
	var out bytes.Buffer

	logger := logrus.New()
	logger.SetOutput(&out)

	sqlLogger := logrusadapter.NewSQLLogger(logger)
	sqlLogger.Connect(42)
	sqlLogger.ConnBegin(42, 43, driver.TxOptions{})
	sqlLogger.ConnQuery(42, 44, "SELECT 1", nil)
	sqlLogger.TxCommit(43)
	sqlLogger.ConnClose(42)

	actualLog := out.String()
	expectedLogLines := []string{
		`level=info msg="CONN Begin" connID=42`,
		`level=info msg="CONN Query" args="[]" connID=42 query="SELECT 1" rowsID=44`,
		`level=info msg="TX Commit" txID=43`,
	}
	for i, logLine := range expectedLogLines {
		if !strings.Contains(actualLog, logLine) {
			t.Fatalf("expected log line %d:\n%s\n, but got:\n%s\n", i, logLine, actualLog)
		}
	}

}
