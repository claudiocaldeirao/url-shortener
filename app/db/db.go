package db

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var db *dynamodb.DynamoDB

func Init() {
	db = dynamodb.New(session.New(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewEnvCredentials(),
		Endpoint:    aws.String(os.Getenv("AWS_DYNAMO_DB_ENDPOINT")),
		DisableSSL:  aws.Bool(true),
	}))
}

func GetDB() *dynamodb.DynamoDB {
	return db
}
