package repo

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

func collectOneRow[T any](rows pgx.Rows) (*T, error) {
	bid, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByPos[T])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNoRows
	}
	return bid, err
}

func collectExactlyOneRow[T any](rows pgx.Rows) (*T, error) {
	bid, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByPos[T])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNoRows
	}
	if errors.Is(err, pgx.ErrTooManyRows) {
		return nil, ErrTooManyRows
	}
	return bid, err
}
