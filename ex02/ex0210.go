package main

import (
	"context"
	"database/sql"
)

func Ex0210(ctx context.Context, db *sql.DB) error {
	var conn [5]*sql.Conn
	db.SetMaxOpenConns(len(conn))
	db.SetMaxIdleConns(len(conn))
	for i := range conn {
		conn[i], _ = db.Conn(ctx)
	}
	for _, c := range conn {
		c.Close()
	}

	stats(db, "before")
	db.Close() // プール全体のコネクションがクローズされてしまう
	stats(db, "after")

	return nil
}
