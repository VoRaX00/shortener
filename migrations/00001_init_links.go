package migrations

import (
	"context"
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upLinks, downLinks)
}

func upLinks(ctx context.Context, tx *sql.Tx) error {
	query := `CREATE TABLE IF NOT EXISTS links (
    	id TEXT PRIMARY KEY NOT NULL,
    	link TEXT NOT NULL UNIQUE
	)`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}

func downLinks(ctx context.Context, tx *sql.Tx) error {
	query := `DROP TABLE IF EXISTS links`
	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}
