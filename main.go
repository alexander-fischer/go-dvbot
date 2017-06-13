package main

import (
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
	"os"
)

const (
	verifyToken = "dvb_bot_is_boss"
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

		createBotAnswer(bodyBytes)
	}
}

func verifyTokenAction(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("hub.verify_token") == verifyToken {
		log.Print("verify token success.")
		fmt.Fprintf(w, r.URL.Query().Get("hub.challenge"))
	} else {
		log.Print("Error: verify token failed.")
		fmt.Fprint(w, "Error, wrong validation token")
	}
}

func createBotAnswer(jsonBytes []byte) Answer {
	textInfo := TextInfo{}

	messageStr := GetMessageFromRequest(jsonBytes)
	textInfo.text = messageStr

	stopMap, lines, isDelay := ProcessText(messageStr)
	textInfo.stops = stopMap
	textInfo.lines = lines
	textInfo.delay = isDelay

	textInfo.category = Categorize(textInfo)

	answer := ProcessAnswer(textInfo)

	log.Print("The answer would be: ")
	log.Print(answer.text)

	return answer
}
