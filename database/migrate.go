package database

import (
	"context"
	"embed"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func Migrate(db *pgxpool.Pool) error {
	ctx := context.Background()

	initSQL, err := migrationFiles.ReadFile("migrations/001_init.up.sql")
	if err != nil {
		return fmt.Errorf("read migration: %w", err)
	}

	_, err = db.Exec(ctx, string(initSQL))
	if err != nil {
		return fmt.Errorf("apply migration: %w", err)
	}

	log.Println("📦 Migrations applied")
	return nil
}
