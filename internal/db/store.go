package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type Store interface {
	Querier
	ExecTx(ctx context.Context, fn func(*Queries) error) error
	IsUniqueViolation(err error) bool
	IsForeignKeyViolation(err error) bool
	IsNoRows(err error) bool
}

type SQLStore struct {
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

func (store *SQLStore) IsUniqueViolation(err error) bool {
	pqErr, ok := err.(*pq.Error)
	return ok && pqErr.Code == "23505"
}

func (store *SQLStore) IsForeignKeyViolation(err error) bool {
	pqErr, ok := err.(*pq.Error)
	return ok && pqErr.Code == "23503"
}

func (store *SQLStore) IsNoRows(err error) bool {
	return err == sql.ErrNoRows
}
