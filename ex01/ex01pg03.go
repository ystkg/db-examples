package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"
)

func Ex01Pg03(ctx context.Context, db *sql.DB) error {
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
	var insertId int32
	now := time.Now()
	if err := tx.QueryRowContext(ctx,
		"INSERT INTO movie (title, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id",
		"タイトルA", now, now,
	).Scan(
		&insertId,
	); err != nil {
		return err
	}
	slog.InfoContext(ctx, "INSERT", "insertId", insertId)

	// SELECT
	var id int32
	var title string
	var createdAt, updatedAt time.Time
	if err = tx.QueryRowContext(ctx,
		"SELECT id, title, created_at, updated_at FROM movie WHERE id = $1",
		insertId,
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
