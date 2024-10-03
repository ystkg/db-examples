package main

import (
	"context"
	"database/sql"
	"time"
)

func Ex0202(ctx context.Context, db *sql.DB) error {
	now := time.Now()

	conn, _ := db.Conn(ctx)
	tx, _ := conn.BeginTx(ctx, nil)
	tx.ExecContext(ctx, "INSERT INTO shop (name, created_at) VALUES ($1, $2)", "shop1", now)
	tx.ExecContext(ctx, "INSERT INTO shop (name, created_at) VALUES ($1, $2)", "shop2", now)
	tx.Commit()
	tx.Rollback() // sql.ErrTxDone。Commit後にdeferでの実行を想定

	stats(db, "before")
	conn.Close() // ここでプールに返却。ここもdeferでの実行を想定
	stats(db, "after")

	return nil
}
