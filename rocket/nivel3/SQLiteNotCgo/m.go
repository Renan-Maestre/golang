package main

import (
	"fmt"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

func main() {
	conn, err := sqlite.OpenConn("./baz.db", sqlite.OpenReadWrite|sqlite.OpenCreate)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Execute a query.
	err = sqlitex.ExecuteTransient(conn, "SELECT 'hello, world';", &sqlitex.ExecOptions{
		ResultFunc: func(stmt *sqlite.Stmt) error {
			fmt.Println(stmt.ColumnText(0))
			return nil
		},
	})
	if err != nil {
		panic(err)
	}
}
