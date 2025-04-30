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
	shortUrl, err := GetUrl("abcedf", db)

	if err != nil {
		log.Println(err)
	}

	fmt.Println(shortUrl)
}

func GetUrl(code string, db *dynamodb.DynamoDB) (*model.ShortUrl, error) {
	params := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"shortUrl": {
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

	var shortUrl *model.ShortUrl
	err = dynamodbattribute.UnmarshalMap(resp.Item, &shortUrl)
	return shortUrl, err
}

func PostUrl(hash string, url string, db *dynamodb.DynamoDB) (*model.ShortUrl, error) {
	shortUrl := model.ShortUrl{
		Hash: hash,
		Url:  url,
	}

	serializedShortUrl, err := dynamodbattribute.MarshalMap(shortUrl)

	if err != nil {
		return nil, err
	}

	params := &dynamodb.PutItemInput{
		Item:      serializedShortUrl,
		TableName: aws.String(os.Getenv("AWS_DYNAMO_DB_TABLE")),
	}

	if _, err := db.PutItem(params); err != nil {
		return nil, err
	}
	return &shortUrl, nil
}
