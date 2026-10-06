package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	urlExample := "postgres://pg:password@localhost:5433/tests"

	// 1 conexao
	// db, err := pgx.Connect(context.Background(), urlExample)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to connect to database %v \n", err)
	// }
	// defer db.Close(context.Background())

	// mais de 1 conexão

	db, err := pgxpool.New(context.Background(), urlExample)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(context.Background()); err != nil {
		panic(err)
	}

	// query := "create table foo (id bigserial primary key, bar varchar(255));"
	// if _, err := db.Exec(context.Background(), query); err != nil {
	// 	panic(err)
	// }

	query := "insert into foo (bar) values ($1);"
	if _, err := db.Exec(context.Background(), query, "abcdef"); err != nil {
		panic(err)
	}

	query = "select * from foo limit 1;"

	type foobar struct {
		id  int64
		bar string
	}

	var baz foobar

	if err = db.QueryRow(context.Background(), query).Scan(&baz.id, &baz.bar); err != nil {
		panic(err)
	}

	fmt.Printf("%#+v\n", baz)
}
