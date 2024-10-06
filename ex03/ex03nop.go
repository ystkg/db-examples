//go:build !deprecated

package main

import (
	"context"
	"database/sql"
	"errors"
)

var ErrNotImplemented = errors.New("not implemented")

func Ex03MySQL03(ctx context.Context, db *sql.DB) error {
	return ErrNotImplemented
}

func Ex03MySQL04(ctx context.Context, db *sql.DB) error {
	return ErrNotImplemented
}

func Ex03Pg03(ctx context.Context, db *sql.DB) error {
	return ErrNotImplemented
}

func Ex03Pg04(ctx context.Context, db *sql.DB) error {
	return ErrNotImplemented
}

func Ex03Pg05(ctx context.Context, db *sql.DB) error {
	return ErrNotImplemented
}
