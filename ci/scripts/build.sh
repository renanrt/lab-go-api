#!/bin/sh
env GOOS=linux GOARCH=arm
CGO_ENABLED=0 go build
# CGO_ENABLED=0 go run main.go migrate
