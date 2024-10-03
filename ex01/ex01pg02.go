package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"
)

func Ex01Pg02(ctx context.Context, db *sql.DB) error {
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
		"INSERT INTO movie (title, created_at, updated_at) VALUES ($1, $2, $3)",
		"タイトルA", now, now,
	)
	if err != nil {
		return err
	}
	lastInsertId, errLastInsertId := result.LastInsertId()
	rowsAffected, errRowsAffected := result.RowsAffected()
	slog.InfoContext(ctx, "INSERT", "lastInsertId", lastInsertId, "errLastInsertId", errLastInsertId, "rowsAffected", rowsAffected, "errRowsAffected", errRowsAffected)

	// コミット
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
