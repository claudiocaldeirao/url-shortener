package model

type ShortUrl struct {
	Hash string `dynamodbav:"Hash"`
	Url  string `dynamodbav:"url"`
}
