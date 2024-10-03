package main

import (
	"context"
	"database/sql"
	"time"
)

func Ex0207(ctx context.Context, db *sql.DB) error {
	db.ExecContext(ctx, "INSERT INTO shop (name, created_at) VALUES ($1, $2)", "shop1", time.Now())

	// 1レコードだけ取得
	row := db.QueryRowContext(ctx, "SELECT id, name FROM shop ORDER BY id LIMIT 1")

	var id int32
	var name string

	stats(db, "before")
	row.Scan(&id, &name) // ここで返却
	stats(db, "after")

	return nil
}
