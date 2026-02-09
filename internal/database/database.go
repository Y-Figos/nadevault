package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type testeStruct struct {
	name string
	slug string
	id   int
}

func Connect(ctx context.Context, url string) *testeStruct {
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)
	var name string
	var slug string
	var id int
	err = conn.QueryRow(context.Background(), "select id, name, slug from maps").Scan(&id, &name, &slug)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	return &testeStruct{
		name: name,
		slug: slug,
		id:   id,
	}
}
