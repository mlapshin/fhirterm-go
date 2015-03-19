package main

import (
	"flag"
	"fmt"
	"github.com/mlapshin/fhirterm/db"
	"github.com/mlapshin/fhirterm/importer"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "ftdb: Command-line FHIRterm database utility\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
	}

	var dbPath = flag.String("db", "", "Path to SQLite database")
	var action = flag.String("action", "", "Action to perform")
	var inputFile = flag.String("file", "", "Source file containing dataset to import")

	flag.Parse()

	db, _ := db.Open(*dbPath, 1)
	defer db.Close()

	var err error

	switch *action {
	case "import-loinc":
		err = importer.ImportLoinc(db, *inputFile)
	case "import-snomed":
		err = importer.ImportSnomed(db, *inputFile)
	default:
		fmt.Fprintf(os.Stderr, "Unknown action: %s\n", *action)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	}
}
