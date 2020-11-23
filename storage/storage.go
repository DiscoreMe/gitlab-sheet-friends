package storage

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sqlx.DB
}

func NewStorage(connection string) (*Storage, error) {
	db, err := sqlx.Connect("sqlite3", connection)
	if err != nil {
		return nil, err
	}
	return &Storage{
		db: db,
	}, nil
}
