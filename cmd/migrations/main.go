package main

import (
	"database/sql"
	"embed"

	"github.com/connor-davis/zingfibre-core/common"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	db, err := sql.Open("pgx", common.EnvString("POSTGRES_DSN", ""))

	if err != nil {
		panic(err)
	}

	defer db.Close()

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}
}
