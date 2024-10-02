package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-sql-driver/mysql"
)

func Ex04Xa02(ctx context.Context, pgDB, myDB *sql.DB) error {
	err := ex04Xa02Pg(ctx, pgDB)
	if err != nil {
		return err
	}

	err = ex04Xa02MySQL(ctx, myDB)
	if err != nil {
		return err
	}

	return nil
}

func ex04Xa02Pg(ctx context.Context, db *sql.DB) error {
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			slog.WarnContext(ctx, "Close", "err", err)
		}
	}()

	name := "shop4th"
	transactionId := "shop4th2pc"
	prepared := false

	// トランザクション開始
	if _, err = conn.ExecContext(ctx, "begin"); err != nil {
		return err
	}
	defer func() {
		if prepared {
			return
		}
		// ロールバック
		if _, err = conn.ExecContext(ctx, "rollback"); err != nil {
			slog.WarnContext(ctx, "rollback", "err", err)
		}
	}()

	// 登録
	result, err := conn.ExecContext(ctx,
		"INSERT INTO shop (name) VALUES ($1)",
		name,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	slog.InfoContext(ctx, "INSERT", "RowsAffected", rows)

	// コミット準備
	if _, err = conn.ExecContext(ctx,
		fmt.Sprintf("prepare transaction '%s'", transactionId),
	); err != nil {
		return err
	}
	prepared = true
	slog.InfoContext(ctx, "prepare transaction")

	return nil
}

func ex04Xa02MySQL(ctx context.Context, db *sql.DB) error {
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			slog.WarnContext(ctx, "Close", "err", err)
		}
	}()
	name := "shop4th"
	transactionId := "shop4th2pc"
	prepared := false

	// トランザクション開始
	if _, err = conn.ExecContext(ctx,
		fmt.Sprintf("XA BEGIN '%s'", transactionId),
	); err != nil {
		return err
	}
	defer func() {
		if prepared {
			return
		}
		// ロールバック
		if _, err := conn.ExecContext(ctx,
			fmt.Sprintf("XA ROLLBACK '%s'", transactionId),
		); err != nil {
			var myerr *mysql.MySQLError
			if errors.As(err, &myerr); myerr.Number == 1397 {
				// Error 1397 (XAE04): XAER_NOTA: Unknown XID
				return
			}
			slog.WarnContext(ctx, "ROLLBACK", "err", err)
		}
	}()

	// 削除
	result, err := conn.ExecContext(ctx,
		"DELETE FROM shop WHERE NAME = ?",
		name,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	slog.InfoContext(ctx, "DELETE", "RowsAffected", rows)

	// コミット準備
	if _, err = conn.ExecContext(ctx,
		fmt.Sprintf("XA END '%s'", transactionId),
	); err != nil {
		return err
	}
	if _, err = conn.ExecContext(ctx,
		fmt.Sprintf("XA PREPARE '%s'", transactionId),
	); err != nil {
		return err
	}
	prepared = true
	slog.InfoContext(ctx, "XA PREPARE")

	return nil
}
