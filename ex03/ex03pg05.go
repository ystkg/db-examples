//go:build deprecated

package main

import (
	"context"
	"database/sql"
	"fmt"
)

// Deprecated: 対比説明用
func Ex03Pg05(ctx context.Context, db *sql.DB) error {
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	param := "Bob' OR '1' = '1" // 不正なパラメータ

	// PrepareContext
	stmt, err := conn.PrepareContext(ctx,
		fmt.Sprintf("SELECT id, name, role FROM staff WHERE name = '%s'",
			param,
		))
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return err
	}
	rows.Close() // 即クローズ

	return nil
}
