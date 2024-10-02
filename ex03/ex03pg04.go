package main

import (
	"context"
	"database/sql"
	"fmt"
)

func Ex03Pg04(ctx context.Context, db *sql.DB) error {
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	param := "Bob' OR '1' = '1" // 不正なパラメータ

	rows, err := conn.QueryContext(ctx,
		fmt.Sprintf("SELECT id, name, role FROM staff WHERE name = '%s'",
			param,
		))
	if err != nil {
		return err
	}
	rows.Close() // 即クローズ

	return nil
}
