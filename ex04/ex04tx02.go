package main

import (
	"context"
	"database/sql"
	"log/slog"
)

func Ex04Tx02(ctx context.Context, pgDB, myDB *sql.DB) error {
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

	name := "shop2nd"
	committed := false

	// トランザクション開始（PostgreSQL）
	_, err = pgConn.ExecContext(ctx, "begin")
	if err != nil {
		return err
	}
	defer func() {
		if committed {
			return
		}
		// ロールバック
		if _, err = pgConn.ExecContext(ctx, "rollback"); err != nil {
			slog.WarnContext(ctx, "Rollback", "err", err)
		}
	}()

	// トランザクション開始（MySQL）
	_, err = myConn.ExecContext(ctx, "START TRANSACTION")
	if err != nil {
		return err
	}
	defer func() {
		if committed {
			return
		}
		// ロールバック
		if _, err = myConn.ExecContext(ctx, "ROLLBACK"); err != nil {
			slog.WarnContext(ctx, "Rollback", "err", err)
		}
	}()

	// 登録（PostgreSQL）
	result, err := pgConn.ExecContext(ctx,
		"INSERT INTO shop (name) VALUES ($1)",
		name,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	slog.InfoContext(ctx, "INSERT", "RowsAffected", rows)

	// 削除（MySQL）
	result, err = myConn.ExecContext(ctx,
		"DELETE FROM shop WHERE NAME = ?",
		name,
	)
	if err != nil {
		return err
	}
	rows, _ = result.RowsAffected()
	slog.InfoContext(ctx, "DELETE", "RowsAffected", rows)

	// コミット（PostgreSQL）
	if _, err = pgConn.ExecContext(ctx, "commit"); err != nil {
		return err // MySQL側はロールバックされるので両方とも未反映になる
	}

	// コミット（MySQL）
	if _, err = myConn.ExecContext(ctx, "COMMIT"); err != nil {
		return err // MySQLだけ未反映となり不整合が生じる
	}
	committed = true

	return nil
}
