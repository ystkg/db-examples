package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"
)

func Ex01MySQL02(ctx context.Context, db *sql.DB) error {
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
	result, err := tx.ExecContext(ctx,
		"INSERT INTO movie (title, created_at, updated_at) VALUES (?, ?, ?)",
		"タイトルA", now, now,
	)
	if err != nil {
		return err
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	slog.InfoContext(ctx, "INSERT", "lastInsertId", lastInsertId, "rowsAffected", rowsAffected)

	// SELECT
	var id int32
	var title string
	var createdAt, updatedAt time.Time
	if err = tx.QueryRowContext(ctx,
		"SELECT id, title, created_at, updated_at FROM movie WHERE id = ?",
		lastInsertId,
	).Scan(
		&id, &title, &createdAt, &updatedAt,
	); err != nil {
		return err
	}
	slog.InfoContext(ctx, "SELECT", "id", id, "title", title, "created_at", createdAt, "updated_at", updatedAt)

	// コミット
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
