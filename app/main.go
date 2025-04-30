package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"url-shortener/app/db"
	"url-shortener/app/model"
	"url-shortener/app/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type ShortenRequest struct {
	Url string `json:"url"`
}

type ShortenResponse struct {
	Shortcode string `json:"shortcode"`
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	db.Init()
	db := db.GetDB()

	switch request.HTTPMethod {
	case "POST":
		// Handle POST /shorten
		var req ShortenRequest
		err := json.Unmarshal([]byte(request.Body), &req)
		if err != nil || req.Url == "" {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       `{"error":"Invalid request body"}`,
			}, nil
		}

		hash := utils.GenerateShortCode()
		_, err = PostUrl(hash, req.Url, db)

		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       fmt.Sprintf(`{"error": "%s"}`, err.Error()),
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			}, nil
		}

		resp := ShortenResponse{Shortcode: hash}
		respBody, _ := json.Marshal(resp)

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(respBody),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	case "GET":
		// Handle GET /{hash}
		pathParts := strings.Split(strings.Trim(request.Path, "/"), "/")
		if len(pathParts) != 1 {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       `{"error":"Invalid path"}`,
			}, nil
		}

		hash := pathParts[0]
		shortUrl, err := GetUrl(hash, db)

		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Body:       `{"error":"Shortcode not found"}`,
			}, nil
		}

		// Redirect to the original URL
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMovedPermanently,
			Headers: map[string]string{
				"Location": shortUrl.Url,
			},
		}, nil
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       `{"error":"Method not allowed"}`,
		}, nil
	}
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
