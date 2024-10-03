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

func Ex01Pg05(ctx context.Context, db *sql.DB) error {
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
	rows, err := tx.QueryContext(ctx,
		"INSERT INTO movie (title, created_at, updated_at) VALUES ($1, $2, $3), ($4, $5, $6), ($7, $8, $9) RETURNING id, title, created_at, updated_at",
		"タイトルA", now, now,
		"タイトルB", now, now,
		"タイトルC", now, now,
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
		slog.InfoContext(ctx, "INSERT", "id", id, "title", title, "created_at", createdAt, "updated_at", updatedAt)
		ids = append(ids, id)
	}

	// DELETE
	ph := make([]string, len(ids))
	for i := range ph {
		ph[i] = fmt.Sprintf("$%d", i+1) // プレースホルダ
	}
	result, err := tx.ExecContext(ctx,
		fmt.Sprintf("DELETE FROM movie WHERE id IN (%s)", strings.Join(ph, ",")),
		ids...,
	)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	slog.InfoContext(ctx, "DELETE", "rowsAffected", rowsAffected)

	// コミット
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
