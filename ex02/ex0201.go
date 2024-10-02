package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"
)

func Ex0201(ctx context.Context, db *sql.DB) error {
	now := time.Now()

	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer func() {
		// ここでプールに返却
		if err := conn.Close(); err != nil {
			slog.WarnContext(ctx, "Close", "err", err)
		}
	}()

	// トランザクション開始
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		// ロールバック
		if err := tx.Rollback(); !errors.Is(err, sql.ErrTxDone) {
			slog.WarnContext(ctx, "Rollback", "err", err)
		}
	}()

	// INSERT
	if _, err = tx.ExecContext(ctx,
		"INSERT INTO shop (name, created_at) VALUES ($1, $2)",
		"shop1", now,
	); err != nil {
		return err
	}

	// INSERT
	if _, err = tx.ExecContext(ctx,
		"INSERT INTO shop (name, created_at) VALUES ($1, $2)",
		"shop2", now,
	); err != nil {
		return err
	}

	// コミット
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
