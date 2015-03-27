package main

import (
	"flag"
	"fmt"
	"github.com/mlapshin/fhirterm"
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
	var err error

	flag.Parse()

	err = fhirterm.OpenDbSpecificFile(*dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening SQLite database: %s\n", err)
		os.Exit(1)
	}
	defer fhirterm.CloseDb()

	switch *action {
	case "import-loinc":
		err = importer.ImportLoinc(fhirterm.GetDb(), *inputFile)
	case "import-snomed":
		err = importer.ImportSnomed(fhirterm.GetDb(), *inputFile)
	default:
		fmt.Fprintf(os.Stderr, "Unknown action: %s\n", *action)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}
