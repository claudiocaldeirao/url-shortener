package main

import (
	"fmt"
	"log"
	"os"
	"url-shortener/app/config"
	"url-shortener/app/db"
	"url-shortener/app/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func main() {
	config.Load()
	db.Init()
	db := db.GetDB()

	// @todo: should get code from GET route
	shortCode, err := GetUrl("abcedf", db)

	if err != nil {
		log.Println(err)
	}

	fmt.Println(shortCode)
}

func GetUrl(code string, db *dynamodb.DynamoDB) (*model.ShortCode, error) {
	params := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"shortCode": {
				S: aws.String(code),
			},
		},
		TableName:      aws.String(os.Getenv("AWS_DYNAMO_DB_TABLE")),
		ConsistentRead: aws.Bool(true),
	}

	resp, err := db.GetItem(params)

	if err != nil {
		return nil, err
	}

	var shortCode *model.ShortCode
	err = dynamodbattribute.UnmarshalMap(resp.Item, &shortCode)
	return shortCode, err
}
