package utils

import "math/rand"

const (
	allowedChards = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortCodeSize = 6
)

func GenerateShortCode() string {
	shortCodeinBytes := make([]byte, shortCodeSize)

	for i := range shortCodeinBytes {
		shortCodeinBytes[i] = allowedChards[rand.Intn(len(allowedChards))]
	}

	return string(shortCodeinBytes)
}
