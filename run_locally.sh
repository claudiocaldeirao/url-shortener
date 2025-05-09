#!/bin/bash
GOOS=linux GOARCH=amd64 go build -o ./build/bootstrap ./app/cmd/main.go

# 2. Start the API locally
sam local start-api
