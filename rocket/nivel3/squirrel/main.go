package main

import (
	"fmt"

	"github.com/Masterminds/squirrel"
)

func main() {
	filters := Filters{
		Id:       1234,
		Name:     "foo",
		email:    "bar",
		Username: "baz",
	}

	sql, args := build(filters)
	fmt.Println(sql)
	fmt.Println(args)
}

type Filters struct {
	Id       int64
	Name     string
	email    string
	Username string
}

func build(f Filters) (string, []any) {
	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Select("*").From("users")
	or := squirrel.Or{}

	if f.Id > 0 {
		or = append(or, squirrel.Eq{"id": f.Id})
	}

	if f.Name != "" {
		or = append(or, squirrel.Like{"name": "%" + f.Name + "%"})
	}

	if f.email != "" {
		or = append(or, squirrel.Eq{"email": f.email})
	}

	if f.Username != "" {
		or = append(or, squirrel.Eq{"Username": f.Username})
	}

	sql, args, err := builder.Where(or).ToSql()
	if err != nil {
		panic(err)
	}

	return sql, args
}
