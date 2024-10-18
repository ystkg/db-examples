//go:build deprecated

package main

import (
	"context"
	"database/sql"
	"fmt"
)

// Deprecated: 対比説明用
func Ex03MySQL03(ctx context.Context, db *sql.DB) error {
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	// find Bob
	param := "Bob"
	rows, err := conn.QueryContext(ctx,
		fmt.Sprintf("SELECT id, name, role FROM staff WHERE name = '%s'",
			param,
		))
	if err != nil {
		return err
	}
	rows.Close() // 即クローズ

	// find Carol
	param = "Carol"
	rows, err = conn.QueryContext(ctx,
		fmt.Sprintf("SELECT id, name, role FROM staff WHERE name = '%s'",
			param,
		))
	if err != nil {
		return err
	}
	rows.Close() // 即クローズ

	return nil
}
