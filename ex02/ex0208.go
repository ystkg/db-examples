package main

import (
	"context"
	"database/sql"
	"time"
)

func Ex0208(ctx context.Context, db *sql.DB) error {
	now := time.Now()
	db.ExecContext(ctx, "INSERT INTO shop (name, created_at) VALUES ($1, $2)", "shop1", now)
	db.ExecContext(ctx, "INSERT INTO shop (name, created_at) VALUES ($1, $2)", "shop2", now)

	// 2レコード取得
	rows, _ := db.QueryContext(ctx, "SELECT id, name FROM shop ORDER BY id LIMIT 2")

	var id int32
	var name string

	// 1レコード目
	rows.Next() // （trueなので）ここでは返却されない
	rows.Scan(&id, &name)

	// 2レコード目
	rows.Next() // （trueなので）ここでは返却されない
	rows.Scan(&id, &name)

	stats(db, "before")
	rows.Next() // （falseなので）ここで返却
	stats(db, "after")

	rows.Close() // ここではない

	return nil
}
