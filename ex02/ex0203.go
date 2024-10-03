package main

import (
	"context"
	"database/sql"
	"time"
)

func Ex0203(ctx context.Context, db *sql.DB) error {
	now := time.Now()

	tx, _ := db.BeginTx(ctx, nil)
	tx.ExecContext(ctx, "INSERT INTO shop (name, created_at) VALUES ($1, $2)", "shop1", now)
	tx.ExecContext(ctx, "INSERT INTO shop (name, created_at) VALUES ($1, $2)", "shop2", now)

	stats(db, "before")
	tx.Commit() // ここで返却
	stats(db, "after")

	tx.Rollback() // sql.ErrTxDone。Commit後にdeferでの実行を想定

	return nil
}
