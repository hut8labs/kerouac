#!/bin/bash

set -e

go get code.google.com/p/go-sqlite/go1/sqlite3

go vet
go test
go build
