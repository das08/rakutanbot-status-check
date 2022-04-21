package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/das08/rakutanbot-status-check/model"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Env struct {
	LineChannelSecret string
	KRBEndpoint       string
}

func loadEnv() Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return Env{
		LineChannelSecret: os.Getenv("LINE_CHANNEL_SECRET"),
		KRBEndpoint:       os.Getenv("KRB_ENDPOINT"),
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

	endpoint := env.KRBEndpoint
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBytes))
	if err != nil {
		panic("Error")
	}

	req.Header.Set("X-Line-Signature", generateSignature(env, jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)

	resp, err := client.Do(req)
	if err != nil {
		panic("Error2")
	}
	defer resp.Body.Close()

	fmt.Println("status: ", resp.Status)

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("Error")
	}

	parsedBody, _ := model.UnmarshalResponse(byteArray)
	fmt.Printf("%#v, %s", parsedBody.Status, *parsedBody.Text)

}
