package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/lib/pq"
	"github.com/networkteam/go-sqllogger"
)

func main() {
	l := log.New(os.Stderr, "SQL: ", 0)

	logger := sqllogger.NewDefaultLogger(l)
	logger.LogClose = true

	pqConnector, err := pq.NewConnector("dbname=test sslmode=disable")
	if err != nil {
		failf("could not connect to database: %v", err)
	}
	connector := sqllogger.LoggingConnector(logger, pqConnector)

	db := sql.OpenDB(connector)

	ctx := context.Background()

	rows, err := db.QueryContext(ctx, "SELECT 42")
	if err != nil {
		failf("could not query database: %v", err)
	}
	defer rows.Close()

	if rows.Next() {
		var answer int
		if err := rows.Scan(&answer); err != nil {
			failf("could not scan row: %v", err)
		}

		fmt.Printf("The answer is: %d\n", answer)
	}
}

func failf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
	os.Exit(1)
}
