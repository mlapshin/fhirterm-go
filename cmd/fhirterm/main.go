package main

import (
	"flag"
	"fmt"
	"github.com/mlapshin/fhirterm"
	"log"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "fhirterm: FHIR Terminology Server\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
	}

	var configPath = flag.String("c", "", "Path to FHIRterm config file")

	flag.Parse()

	if *configPath == "" {
		fmt.Fprintf(os.Stderr, "No config file specified.\n")
		flag.Usage()
		os.Exit(1)
	}

	config, err := fhirterm.ReadConfig(*configPath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	log.Printf("Using config file: %s", *configPath)

	err = fhirterm.OpenDb(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	err = fhirterm.InitStorage(config.Storage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing Storage: %s\n", err)
		os.Exit(1)
	}

	fhirterm.StartServer(config)

	os.Exit(0)
}
