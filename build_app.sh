#!/bin/bash
GOOS=linux GOARCH=amd64 go build -o ./build/bootstrap ./app/main.go

zip ./build/lambda.zip ./build/bootstrap
