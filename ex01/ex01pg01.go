package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"
)

func Ex01Pg01(ctx context.Context, db *sql.DB) error {
	// トランザクション開始
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); !errors.Is(err, sql.ErrTxDone) {
			slog.WarnContext(ctx, "Rollback", "err", err)
		}
	}()

	// INSERT
	now := time.Now()
	if _, err = tx.ExecContext(ctx,
		"INSERT INTO movie (title, created_at, updated_at) VALUES ($1, $2, $3)",
		"タイトルA", now, now,
	); err != nil {
		return err
	}

	// SELECT
	var id int32
	var title string
	var createdAt, updatedAt time.Time
	if err = tx.QueryRowContext(ctx,
		"SELECT id, title, created_at, updated_at FROM movie ORDER BY id DESC LIMIT 1",
	).Scan(
		&id, &title, &createdAt, &updatedAt,
	); err != nil {
		return err
	}
	slog.InfoContext(ctx, "SELECT", "id", id, "title", title, "created_at", createdAt, "updated_at", updatedAt)

	// DELETE
	if _, err = tx.ExecContext(ctx,
		"DELETE FROM movie WHERE id = $1",
		id,
	); err != nil {
		return err
	}

	// コミット
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
