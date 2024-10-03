package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

func Ex01MySQL03(ctx context.Context, db *sql.DB) error {
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
		"INSERT INTO movie (title, created_at, updated_at) VALUES (?, ?, ?), (?, ?, ?), (?, ?, ?)",
		"タイトルA", now, now,
		"タイトルB", now, now,
		"タイトルC", now, now,
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
	rows, err := tx.QueryContext(ctx,
		"SELECT id, title, created_at, updated_at FROM movie ORDER BY id DESC LIMIT ?",
		rowsAffected,
	)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.WarnContext(ctx, "Close", "err", err)
		}
	}()

	ids := []any{} // ExecContextに渡すためint32ではなくanyにしておく
	for rows.Next() {
		var id int32
		var title string
		var createdAt, updatedAt time.Time
		if err = rows.Scan(&id, &title, &createdAt, &updatedAt); err != nil {
			return err
		}
		slog.InfoContext(ctx, "SELECT", "id", id, "title", title, "created_at", createdAt, "updated_at", updatedAt)
		ids = append(ids, id)
	}

	// DELETE
	ph := strings.Repeat(",?", len(ids)) // プレースホルダ
	result, err = tx.ExecContext(ctx,
		fmt.Sprintf("DELETE FROM movie WHERE id IN (%s)", ph[1:]),
		ids...,
	)
	if err != nil {
		return err
	}
	rowsAffected, _ = result.RowsAffected()
	slog.InfoContext(ctx, "DELETE", "rowsAffected", rowsAffected)

	// コミット
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
