package importer

import (
	"database/sql"
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type unzipCallback func(extractedPath string) error

func unpackZipArchive(zipPath string, callback unzipCallback) error {
	tempPath, err := ioutil.TempDir("", "fhirterm-import")
	if err != nil {
		return err
	}

	defer os.RemoveAll(tempPath)

	unzipCmd := exec.Command("unzip", zipPath, "-d", tempPath)
	if err := unzipCmd.Run(); err != nil {
		log.Fatalf("Cannot unzip archive: %s", err)

		return err
	}

	log.Printf("Unzipped %s => %s", zipPath, tempPath)
	return callback(tempPath)
}

func importCsv(db *sql.DB, csvPath string, comma rune, fpr int, insertStmt string) (int, error) {
	file, err := os.Open(csvPath)

	if err != nil {
		log.Fatalf("Cannot open %s: %s", csvPath, err)
		return 0, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = comma
	reader.FieldsPerRecord = fpr

	rowIdx := 0
	var tx *sql.Tx
	var stmt *sql.Stmt
	stmtArgs := make([]interface{}, fpr)

	for {
		row, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			return 0, err
		}

		// insert 10000 rows per transaction
		if (rowIdx % 10000) == 0 {
			if tx != nil {
				tx.Commit()
			}

			tx, _ = db.Begin()
			stmt, err = tx.Prepare(insertStmt)
			if err != nil {
				return 0, err
			}
		}

		if rowIdx > 0 { // skip first row (useless header)
			for i, v := range row {
				stmtArgs[i] = interface{}(v)
			}

			_, err = stmt.Exec(stmtArgs...)
		}

		rowIdx++
	}
	tx.Commit()

	return rowIdx - 1, nil
}
