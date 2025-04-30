package model

type ShortCode struct {
	shortCode string `dynamodbav:"shortCode"`
	Url       string `dynamodbav:"url"`
}
