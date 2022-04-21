package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
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
	DiscordEndpoint   string
}

type Request struct {
	name      string
	jsonBytes []byte
}

type Result struct {
	Name    string
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
		DiscordEndpoint:   os.Getenv("DISCORD_WEBHOOK"),
	}
}

func loadJSON(filename string) []byte {
	jsonFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return jsonFile
}

func loadCheckRequest(text string) []byte {
	jsonFile := loadJSON("request/mock.json")
	jsonData, err := model.UnmarshalRakutan(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	jsonData.Events[0].Message.Text = text

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

func toPtr(s string) *string {
	return &s
}

func sendRequest(resultChan chan *Result, env Env, r Request) {
	req, err := http.NewRequest("POST", env.KRBEndpoint, bytes.NewBuffer(r.jsonBytes))
	if err != nil {
		resultChan <- &Result{Name: r.name, Status: 9996, Message: toPtr("Could not create request")}
		return
	}

	req.Header.Set("X-Line-Signature", generateSignature(env, r.jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)

	resp, err := client.Do(req)
	if err != nil {
		resultChan <- &Result{Name: r.name, Status: 9997, Message: toPtr("Could not connect to server")}
		return
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resultChan <- &Result{Name: r.name, Status: 9998, Message: toPtr("Could not read body")}
		return
	}

	if resp.StatusCode != 200 {
		resultChan <- &Result{Name: r.name, Status: 9999, Message: toPtr(fmt.Sprintf("Invalid status code : %d", resp.StatusCode))}
		return
	}

	parsedBody, _ := model.UnmarshalResponse(byteArray)
	resultChan <- &Result{Name: r.name, Status: parsedBody.Status, Message: parsedBody.Text}
}

func sendWebhook(whurl string, content string) {
	dw := &model.DiscordWebhook{UserName: "Status Check", Content: content}
	j, err := json.Marshal(dw)
	if err != nil {
		fmt.Println("json err:", err)
		return
	}

	req, err := http.NewRequest("POST", whurl, bytes.NewBuffer(j))
	if err != nil {
		fmt.Println("new request err:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	_, err = client.Do(req)
	if err != nil {
		fmt.Println("client err:", err)
		return
	}
}

func main() {
	// Create channels for go routine
	resultsChan := make(chan *Result, 5)

	env := loadEnv()

	checkList := []Request{
		{name: "自然地理学", jsonBytes: loadCheckRequest("自然地理学")},
		{name: "お気に入り取得", jsonBytes: loadCheckRequest("お気に入り")},
		{name: "楽単おみくじ", jsonBytes: loadCheckRequest("おみくじ")},
		{name: "楽単おみくじ", jsonBytes: loadCheckRequest("おみくじ")},
		{name: "鬼単おみくじ", jsonBytes: loadCheckRequest("鬼単おみくじ")},
		{name: "鬼単おみくじ", jsonBytes: loadCheckRequest("鬼単おみくじ")},
		{name: "#12345", jsonBytes: loadCheckRequest("#12345")},
	}

	for _, v := range checkList {
		go sendRequest(resultsChan, env, v)
	}

	for _, _ = range checkList {
		errorMsg := ""
		result := <-resultsChan
		if result.Message != nil {
			errorMsg = *result.Message
		}

		// Send Discord Webhook
		if result.Status != 2000 {
			errorMsg := fmt.Sprintf("[Error][%s] Code: %d \n Message: %s", result.Name, result.Status, errorMsg)
			sendWebhook(env.DiscordEndpoint, errorMsg)
		}
	}

	// Close channels
	close(resultsChan)
	sendWebhook(env.DiscordEndpoint, "fin")
}
