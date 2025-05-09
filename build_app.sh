#!/bin/bash
GOOS=linux GOARCH=amd64 go build -o ./build/bootstrap ./app/cmd/main.go

cd ./build
zip lambda.zip bootstrap
cd ..
