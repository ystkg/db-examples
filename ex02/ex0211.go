package main

import (
	"context"
	"database/sql"
)

func Ex0211(ctx context.Context, db *sql.DB) error {
	var conn [5]*sql.Conn
	db.SetMaxOpenConns(len(conn))
	db.SetMaxIdleConns(len(conn))
	for i := range conn {
		conn[i], _ = db.Conn(ctx)
	}
	conn[2].Close()
	conn[3].Close()

	stats(db, "before")
	db.Close() // Idleのコネクションがクローズされ、InUseのコネクションは残る
	stats(db, "after")

	conn[0].Close() // プールに返却されず、直接クローズされる
	stats(db, "0")

	conn[1].Close() // プールに返却されず、直接クローズされる
	stats(db, "1")

	conn[2].Close() // ErrConnDone
	stats(db, "2")

	conn[3].Close() // ErrConnDone
	stats(db, "3")

	conn[4].Close() // プールに返却されず、直接クローズされる
	stats(db, "4")

	return nil
}
