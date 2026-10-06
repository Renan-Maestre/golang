package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	urlExample := "postgres://pg:password@localhost:5433/tests"

	db, err := pgxpool.New(context.Background(), urlExample)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(context.Background()); err != nil {
		panic(err)
	}

	queries := New(db)
	ctx := context.Background()

	authors, err := queries.ListAuthors(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(authors)

	authorCriado, err := queries.CreateAuthor(ctx, CreateAuthorParams{
		Name: "Renan Maestre",
		Bio: pgtype.Text{
			String: "dev na faveni",
			Valid:  true,
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(authorCriado)
}
