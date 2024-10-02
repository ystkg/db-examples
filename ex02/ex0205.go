package main

import (
	"context"
	"database/sql"
	"time"
)

func Ex0205(ctx context.Context, db *sql.DB) error {
	now := time.Now()

	conn, _ := db.Conn(ctx)
	tx, _ := conn.BeginTx(ctx, nil)
	tx.ExecContext(ctx, "INSERT INTO shop (name, created_at) VALUES ($1, $2)", "shop1", now)
	tx.ExecContext(ctx, "INSERT INTO shop (name, created_at) VALUES ($1, $2)", "shop2", now)

	stats(db, "before")
	errCommit := tx.Commit() // ここでは返却されない
	statsErr(db, "after", errCommit)

	tx.Rollback() // sql.ErrTxDone。ここでも返却されない

	conn.Close() // ここで返却

	return nil
}
