package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
)

func Ex04Tx01(ctx context.Context, pgDB, myDB *sql.DB) error {
	// PostgreSQL
	pgConn, err := pgDB.Conn(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := pgConn.Close(); err != nil {
			slog.WarnContext(ctx, "Close", "err", err)
		}
	}()

	// MySQL
	myConn, err := myDB.Conn(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := myConn.Close(); err != nil {
			slog.WarnContext(ctx, "Close", "err", err)
		}
	}()

	name := "shop1st"

	// トランザクション開始（PostgreSQL）
	pgTx, err := pgConn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		// ロールバック（PostgreSQL）
		if err := pgTx.Rollback(); !errors.Is(err, sql.ErrTxDone) {
			slog.WarnContext(ctx, "Rollback", "err", err)
		}
	}()

	// トランザクション開始（MySQL）
	myTx, err := myConn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		// ロールバック（MySQL）
		if err := myTx.Rollback(); !errors.Is(err, sql.ErrTxDone) {
			slog.WarnContext(ctx, "Rollback", "err", err)
		}
	}()

	// 登録（PostgreSQL）
	result, err := pgTx.ExecContext(ctx,
		"INSERT INTO shop (name) VALUES ($1)",
		name,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	slog.InfoContext(ctx, "INSERT", "RowsAffected", rows)

	// 削除（MySQL）
	result, err = myTx.ExecContext(ctx,
		"DELETE FROM shop WHERE NAME = ?",
		name,
	)
	if err != nil {
		return err
	}
	rows, _ = result.RowsAffected()
	slog.InfoContext(ctx, "DELETE", "RowsAffected", rows)

	// コミット（PostgreSQL）
	if err = pgTx.Commit(); err != nil {
		return err // MySQL側はロールバックされるので両方とも未反映になる
	}

	// コミット（MySQL）
	if err = myTx.Commit(); err != nil {
		return err // MySQLだけ未反映となり不整合が生じる
	}

	return nil
}
