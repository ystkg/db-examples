package main

import (
	"context"
	"database/sql"
	"time"
)

func Ex0204(ctx context.Context, db *sql.DB) error {
	now := time.Now()

	tx, _ := db.BeginTx(ctx, nil)
	tx.ExecContext(ctx, "INSERT INTO shop (name, created_at) VALUES ($1, $2)", "shop1", now)
	tx.ExecContext(ctx, "INSERT INTO shop (name, created_at) VALUES ($1, $2)", "shop2", now)

	stats(db, "before")
	tx.Rollback() // ここで返却
	stats(db, "after")

	return nil
}
