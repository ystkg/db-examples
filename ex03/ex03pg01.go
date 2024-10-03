package main

import (
	"context"
	"database/sql"
)

func Ex03Pg01(ctx context.Context, db *sql.DB) error {
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	// PrepareContext
	stmt, err := conn.PrepareContext(ctx,
		"SELECT id, name, role FROM staff WHERE name = $1",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// find Bob
	param := "Bob"
	rows, err := stmt.QueryContext(ctx, param)
	if err != nil {
		return err
	}
	rows.Close() // 即クローズ

	// find Carol
	param = "Carol"
	rows, err = stmt.QueryContext(ctx, param)
	if err != nil {
		return err
	}
	rows.Close() // 即クローズ

	return nil
}
