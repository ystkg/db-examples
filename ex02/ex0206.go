package main

import (
	"context"
	"database/sql"
	"time"
)

func Ex0206(ctx context.Context, db *sql.DB) error {
	now := time.Now()

	// 1回目
	db.ExecContext(ctx, "INSERT INTO shop (name, created_at) VALUES ($1, $2)", "shop1", now)

	// 2回目
	stats(db, "before")
	db.ExecContext(ctx, "INSERT INTO shop (name, created_at) VALUES ($1, $2)", "shop1", now)
	stats(db, "after")

	return nil
}
