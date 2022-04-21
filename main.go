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

type Result struct {
	Status  int64
	Message *string
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

func loadCheckRequest(filename string) []byte {
	jsonFile := loadJSON(fmt.Sprintf("request/%s.json", filename))
	jsonData, err := model.UnmarshalRakutan(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	jsonBytes, err := jsonData.Marshal()
	if err != nil {
		log.Fatal(err)
	}
	return jsonBytes
}

func generateSignature(env Env, body []byte) string {
	hash := hmac.New(sha256.New, []byte(env.LineChannelSecret))
	hash.Write(body)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func sendRequest(resultChan chan *Result, env Env, jsonBytes []byte) {
	req, err := http.NewRequest("POST", env.KRBEndpoint, bytes.NewBuffer(jsonBytes))
	if err != nil {
		resultChan <- nil
		panic("Error")
	}

	req.Header.Set("X-Line-Signature", generateSignature(env, jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)

	resp, err := client.Do(req)
	if err != nil {
		resultChan <- nil
		panic("Error2")
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resultChan <- nil
		panic("Error")
	}

	parsedBody, _ := model.UnmarshalResponse(byteArray)
	resultChan <- &Result{Status: parsedBody.Status, Message: parsedBody.Text}
}

func main() {
	// Create channels for go routine
	resultsChan := make(chan *Result, 5)

	env := loadEnv()

	checkList := [][]byte{
		loadCheckRequest("mock_search"),
		loadCheckRequest("mock_getfav"),
	}

	for _, v := range checkList {
		go sendRequest(resultsChan, env, v)
	}

	for _, _ = range checkList {
		result := <-resultsChan
		fmt.Print(result.Status)
		if result.Message != nil {
			fmt.Println(*result.Message)
		}
	}

	// Close channels
	close(resultsChan)
}
