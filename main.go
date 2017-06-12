package main

import (
	"fmt"
	"net/http"
	"log"
)

const (
	port        = "5000"
	verifyToken = "dvb_bot_is_boss"
)

// Main entry point.
func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/webhook", webhookHandler)

	log.Println("listen on port " + port)
	http.ListenAndServe(":"+port, nil)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Server is running")
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		verifyTokenAction(w, r)
	} else if r.Method == "POST" {

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

func createBotAnswer() Answer {
	textInfo := TextInfo{}

	jsonBytes := ReadFileToBytes("json/test.json")

	messageStr := GetMessageFromRequest(jsonBytes)
	textInfo.text = messageStr

	stopMap, lines, isDelay := ProcessText(messageStr)
	textInfo.stops = stopMap
	textInfo.lines = lines
	textInfo.delay = isDelay

	textInfo.category = Categorize(textInfo)

	answer := ProcessAnswer(textInfo)

	fmt.Println(textInfo)
	fmt.Println(answer)

	return answer
}
