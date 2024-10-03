package main

import (
	"context"
	"database/sql"
)

func Ex03Pg02(ctx context.Context, db *sql.DB) error {
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	// find Bob
	param := "Bob"
	rows, err := conn.QueryContext(ctx,
		"SELECT id, name, role FROM staff WHERE name = $1",
		param,
	)
	if err != nil {
		return err
	}
	rows.Close() // 即クローズ

	// find Carol
	param = "Carol"
	rows, err = conn.QueryContext(ctx,
		"SELECT id, name, role FROM staff WHERE name = $1",
		param,
	)
	if err != nil {
		return err
	}
	rows.Close() // 即クローズ

	return nil
}
