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

func Ex01Pg06(ctx context.Context, db *sql.DB) error {
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
		"INSERT INTO movie (title, created_at, updated_at) VALUES ($1, $2, $3), ($4, $5, $6), ($7, $8, $9) RETURNING id",
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

	ids := []any{} // QueryContextに渡すためint32ではなくanyにしておく
	for rows.Next() {
		var id int32
		if err = rows.Scan(&id); err != nil {
			return err
		}
		ids = append(ids, id)
	}
	slog.InfoContext(ctx, "INSERT", "ids", ids)

	// DELETE
	ph := make([]string, len(ids))
	for i := range ph {
		ph[i] = fmt.Sprintf("$%d", i+1) // プレースホルダ
	}
	rows, err = tx.QueryContext(ctx,
		fmt.Sprintf("DELETE FROM movie WHERE id IN (%s) RETURNING id, title, created_at, updated_at", strings.Join(ph, ",")),
		ids...,
	)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.WarnContext(ctx, "Close", "err", err)
		}
	}()

	for rows.Next() {
		var id int32
		var title string
		var createdAt, updatedAt time.Time
		if err = rows.Scan(&id, &title, &createdAt, &updatedAt); err != nil {
			return err
		}
		slog.InfoContext(ctx, "DELETE", "id", id, "title", title, "created_at", createdAt, "updated_at", updatedAt)
	}

	// コミット
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
