package storage

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type Storage struct {
	db *sqlx.DB
}

func NewStorage(connection string) (*Storage, error) {
	var isNewFile bool
	if _, err := os.Stat(connection); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		isNewFile = true
	}

	db, err := sqlx.Connect("sqlite3", connection)
	if err != nil {
		return nil, err
	}

	if isNewFile {
		if err := migrate(db); err != nil {
			return nil, err
		}
	}

	return &Storage{
		db: db,
	}, nil
}

func migrate(db *sqlx.DB) error {
	const sql = `
CREATE TABLE "issues" (
	"id"	INTEGER NOT NULL,
	"git"	TEXT NOT NULL,
	"project_id"	INTEGER NOT NULL,
	"issue_id"	INTEGER NOT NULL,
	"created_at"	DATETIME NOT NULL,
	"list_id"	INTEGER NOT NULL,
	"is_closed"	INTEGER NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT)
);

CREATE TABLE "lists" (
	"id"	INTEGER,
	"name"	TEXT NOT NULL,
	"created_at"	DATETIME NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT)
);
`
	if _, err := db.Exec(sql); err != nil {
		return err
	}

	return nil
}
