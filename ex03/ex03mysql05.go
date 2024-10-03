package main

import (
	"context"
	"database/sql"
)

func Ex03MySQL05(ctx context.Context, db *sql.DB) error {
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	param := "Bob' OR '1' = '1" // 不正なパラメータ

	rows, err := conn.QueryContext(ctx,
		"SELECT id, name, role FROM staff WHERE name = ?",
		param,
	)
	if err != nil {
		return err
	}
	rows.Close() // 即クローズ

	return nil
}
