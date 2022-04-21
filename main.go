package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/das08/rakutanbot-status-check/model"
	"github.com/joho/godotenv"
	"io/ioutil"
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

func loadJSON(filename string) []byte {
	jsonFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return jsonFile
}

func generateSignature(env Env, body []byte) string {
	hash := hmac.New(sha256.New, []byte(env.LineChannelSecret))
	hash.Write(body)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func main() {
	env := loadEnv()
	jsonFile := loadJSON("request/mock_getfav.json")
	jsonData, err := model.UnmarshalRakutan(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	jsonBytes, err := jsonData.Marshal()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(generateSignature(env, jsonBytes))

}
