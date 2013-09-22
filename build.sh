#!/bin/sh

go build -o lurch.base lurch.go
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o lurch.linux lurch.go