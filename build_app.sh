#!/bin/bash
GOOS=linux GOARCH=amd64 go build -o ./build/bootstrap ./app/main.go

cd ./build
zip lambda.zip bootstrap
cd ..
