package model

type ShortUrl struct {
	Hash string `dynamodbav:"shortCode"`
	Url  string `dynamodbav:"url"`
}
