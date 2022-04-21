package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Env struct {
	LineChannelSecret string
}

func loadEnv() Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return Env{
		LineChannelSecret: os.Getenv("LINE_CHANNEL_SECRET"),
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	s3Bucket := os.Getenv("S3_BUCKET")
	secretKey := os.Getenv("SECRET_KEY")
}
