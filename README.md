# go-sqllogger

[![Go Report Card](https://goreportcard.com/badge/github.com/networkteam/go-sqllogger?style=flat-square)](https://goreportcard.com/report/github.com/networkteam/go-sqllogger)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/networkteam/go-sqllogger)
[![GitHub release](https://img.shields.io/github/v/release/networkteam/go-sqllogger.svg?style=flat-square&include_prereleases)](https://github.com/networkteam/go-sqllogger/releases)
[![Build Status](https://travis-ci.com/networkteam/go-sqllogger.svg?branch=master)](https://travis-ci.com/networkteam/go-sqllogger)

**Go SQL driver adapter for logging queries and other SQL operations**

* Any `driver.Connector` can be wrapped with `sqllogger.LoggingConnector(driver.Connector)`
  to log **Connect**, **Prepare**, **Exec**, **Query**, **Commit**, **Rollback** and **Close**
  on database, connection, statement, rows and transaction instances
  (see [./logger.go](logger.go) for all intercepted calls)
* The `sqllogger.Logger` interface can be implemented to log SQL to any logging library
* `sqllogger.NewDefaultLogger(*log.Logger)` offers a standard library implementation
* Zero dependencies

> Note: The adapter has been tested using `github.com/lib/pq`. Other SQL drivers might need additional work.
  As the `database/sql/driver` package offers a lot of optional interfaces, not every advanced feature might work
  as expected when using the logging connector as an adapter to the original driver.

## Example

This needs a running PostgreSQL database.

See [./example/main.go](./example/main.go):

```go
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
	// DefaultLogger logs to standard log package
	l := log.New(os.Stderr, "SQL: ", 0)

	logger := sqllogger.NewDefaultLogger(l)
	logger.LogClose = true

	// Any driver exposing a driver.Connector can be wrapped with logging
	pqConnector, err := pq.NewConnector("dbname=test sslmode=disable")
	if err != nil {
		failf("could not connect to database: %v", err)
	}
	connector := sqllogger.LoggingConnector(logger, pqConnector)

	// db is a standard *sql.DB
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
```

Running the example:

```
> createdb test # If it does not exist
> go run github.com/networkteam/go-sqllogger/example
SQL: Connect → CONN(1)
SQL: CONN(1) ► Query(SELECT 42) → ROWS(2)
The answer is: 42
SQL: ROWS(2) ► Close
```

## License

[MIT License](./LICENSE)
