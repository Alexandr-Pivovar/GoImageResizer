#!/bin/sh
set -e
set -x

go test -timeout=600s -race -coverprofile=cover.out ./...
go tool cover -func cover.out
