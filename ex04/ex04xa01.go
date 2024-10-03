package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgconn"
)

func Ex04Xa01(ctx context.Context, pgDB, myDB *sql.DB) error {
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

	name := "shop3rd"
	transactionId := "shop3rd2pc"
	pgPrepared, committed := false, false

	// トランザクション開始（PostgreSQL）
	if _, err = pgConn.ExecContext(ctx, "begin"); err != nil {
		return err
	}
	defer func() {
		if pgPrepared {
			return
		}
		// ロールバック（PostgreSQL）
		if _, err = pgConn.ExecContext(ctx, "rollback"); err != nil {
			slog.WarnContext(ctx, "rollback", "err", err)
		}
	}()

	// トランザクション開始（MySQL）
	if _, err = myConn.ExecContext(ctx,
		fmt.Sprintf("XA BEGIN '%s'", transactionId),
	); err != nil {
		return err
	}
	defer func() {
		if committed {
			return
		}
		// ロールバック（MySQL）
		if _, err := myConn.ExecContext(ctx,
			fmt.Sprintf("XA ROLLBACK '%s'", transactionId),
		); err != nil {
			var myerr *mysql.MySQLError
			if errors.As(err, &myerr); myerr.Number == 1397 {
				// Error 1397 (XAE04): XAER_NOTA: Unknown XID
				return
			}
			slog.WarnContext(ctx, "XA ROLLBACK", "err", err)
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

	// コミット準備（PostgreSQL）
	if _, err = pgConn.ExecContext(ctx,
		fmt.Sprintf("prepare transaction '%s'", transactionId),
	); err != nil {
		return err
	}
	pgPrepared = true
	defer func() {
		if committed {
			return
		}
		// ロールバック（PostgreSQL）
		if _, err := pgConn.ExecContext(ctx,
			fmt.Sprintf("rollback prepared '%s'", transactionId),
		); err != nil {
			var pgerr *pgconn.PgError
			if errors.As(err, &pgerr); pgerr.Code == "42704" {
				// transaction_id does not exist
				return
			}
			slog.WarnContext(ctx, "rollback prepared", "err", err)
		}
	}()

	// コミット準備（MySQL）
	if _, err = myConn.ExecContext(ctx,
		fmt.Sprintf("XA END '%s'", transactionId),
	); err != nil {
		return err
	}
	if _, err = myConn.ExecContext(ctx,
		fmt.Sprintf("XA PREPARE '%s'", transactionId),
	); err != nil {
		return err
	}

	// コミット（PostgreSQL）
	if _, err = pgConn.ExecContext(ctx,
		fmt.Sprintf("commit prepared '%s'", transactionId),
	); err != nil {
		return err
	}

	// コミット（MySQL）
	if _, err = myConn.ExecContext(ctx,
		fmt.Sprintf("XA COMMIT '%s'", transactionId),
	); err != nil {
		return err
	}
	committed = true

	return nil
}
