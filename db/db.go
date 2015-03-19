package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type DB struct {
	Db *sql.DB
}

func Open(dbFile string, maxPoolSize int) (*DB, error) {
	db, err := sql.Open("sqlite3", dbFile)

	if err != nil {
		log.Fatal("Failed to open SQLite Database %s: %s", dbFile, err)
		return nil, err
	}

	db.SetMaxOpenConns(maxPoolSize)

	log.Printf("Opened SQLite Database %s", dbFile)

	return &DB{db}, nil
}

func (db *DB) Begin() (*sql.Tx, error) {
	return db.Db.Begin()
}

func (db *DB) Exec(stmt string) (sql.Result, error) {
	tx, err := db.Begin()

	if err != nil {
		return nil, err
	}

	result, err := tx.Exec(stmt)

	if err != nil {
		return nil, err
	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (db *DB) Close() (err error) {
	err = db.Db.Close()

	if err != nil {
		log.Printf("Error when closing database: %s", err)
	} else {
		log.Print("Closed database file")
	}

	return err
}
