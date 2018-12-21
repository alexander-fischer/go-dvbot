package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	facebookEndPoint = "https://graph.facebook.com/v2.6/me/messages"
)

// Main entry point.
func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/webhook", webhookHandler)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	log.Println("listen on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Server is running")
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		verifyTokenAction(w, r)
	} else if r.Method == "POST" {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(string(bodyBytes))

		messageStr, senderId := GetMessageFromRequest(bodyBytes)

		answer, _ := createBotAnswer(messageStr, senderId)
		if answer.senderId != "" {
			sendMessage(answer)
		}
	} else {
		log.Println("Request was neither GET nor POST.")
	}
}

func verifyTokenAction(w http.ResponseWriter, r *http.Request) {
	verifyToken := os.Getenv("VERIFYTOKEN")

	if r.URL.Query().Get("hub.verify_token") == verifyToken {
		log.Println("verify token success.")
		fmt.Fprintf(w, r.URL.Query().Get("hub.challenge"))
	} else {
		log.Println("Error: verify token failed.")
		fmt.Fprint(w, "Error, wrong validation token")
	}
}

func createBotAnswer(message string, senderId string) (Answer, TextInfo) {
	textInfo := TextInfo{}

	errText := []string{"Das habe ich leider nicht verstanden."}
	if senderId == "" {
		return Answer{"", 0, errText}, textInfo
	}

	if message == "" {
		return Answer{senderId, 0, errText}, textInfo
	}

	textInfo.text = message

	stopNames, stopIds, lines, isDelay := ProcessText(message)
	textInfo.stopNames = stopNames
	textInfo.stopIds = stopIds
	textInfo.lines = lines
	textInfo.delay = isDelay

	textInfo.category = Categorize(textInfo)

	answer := ProcessAnswer(textInfo)
	answer.senderId = senderId

	return answer, textInfo
}

// Sends message to FB endpoint.
func sendMessage(answer Answer) {
	for _, answerText := range answer.text {
		reqBody := Body{
			Recipient: Recipient{
				Id: answer.senderId,
			},
			Message: Message{
				Text: answerText,
			},
		}

		bodyBytes, _ := json.Marshal(reqBody)

		req, err := http.NewRequest("POST", facebookEndPoint, bytes.NewBuffer(bodyBytes))
		if err != nil {
			log.Print(err)
		}

		accessToken := os.Getenv("ACCESSTOKEN")
		values := url.Values{}
		values.Add("access_token", accessToken)

		req.URL.RawQuery = values.Encode()
		req.Header.Add("Content-Type", "application/json; charset=UTF-8")

		client := &http.Client{Timeout: time.Duration(30 * time.Second)}
		res, err := client.Do(req)
		if err != nil {
			log.Println(err)
		}

		var result map[string]interface{}
		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
		}

		if err := json.Unmarshal(resBody, &result); err != nil {
			log.Println(err)
		}

		log.Println(result)

		res.Body.Close()
	}
}

type Body struct {
	Recipient Recipient `json:"recipient"`
	Message   Message `json:"message"`
}

type Recipient struct {
	Id string `json:"id"`
}

type Message struct {
	Text string `json:"text"`
}
