#!/bin/bash

go build -v -o fhirterm main.go
go build -v -o ftdb cmd/ftdb/main.go
