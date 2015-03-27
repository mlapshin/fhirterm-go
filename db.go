package fhirterm

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var globalDb *sql.DB

func GetDb() *sql.DB {
	return globalDb
}

func OpenDb(cfg *Config) error {
	if globalDb != nil {
		return fmt.Errorf("DB connection already opened")
	}

	dbFile := cfg.Databases[0]
	err := OpenDbSpecificFile(dbFile)
	if err != nil {
		return err
	}

	globalDb.SetMaxOpenConns(100)

	return nil
}

func OpenDbSpecificFile(dbFile string) error {
	var err error
	globalDb, err = sql.Open("sqlite3", dbFile)

	if err != nil {
		log.Fatalf("Failed to open SQLite Database %s: %s", dbFile, err)
		return err
	}

	log.Printf("Opened SQLite Database %s", dbFile)
	return nil
}

func CloseDb() error {
	err := globalDb.Close()

	if err != nil {
		log.Printf("Error when closing database: %s", err)
	} else {
		log.Print("Closed database file")
	}

	globalDb = nil

	return err
}
