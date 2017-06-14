package main

import (
	"io/ioutil"
	"log"
	"github.com/tidwall/gjson"
)

// Just for test purposes.
func ReadFileToBytes(filename string) []byte {
	stream, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	readString := []byte(stream)
	return readString
}

// Extract message string from message.
func GetMessageFromRequest(jsonBytes []byte) (string, string) {
	text := gjson.GetBytes(jsonBytes, "entry.#.messaging.#.message.text")
	senderId := gjson.GetBytes(jsonBytes, "entry.#.messaging.#.sender.id")
	if text.String() == "" || senderId.String() == "" {
		return "", ""
	}

	log.Println(text)

	messageStr := text.Array()[0].Array()[0].String()
	messageSId := senderId.Array()[0].Array()[0].String()

	return messageStr, messageSId
}

// Gets a map with all the transit stops and the id.
func GetAllTransitStops(file string) ([]string, map[string]string) {
	jsonBytes := ReadFileToBytes(file)
	names := gjson.GetBytes(jsonBytes, "array.#.name").Array()
	nrs := gjson.GetBytes(jsonBytes, "array.#.nr").Array()

	namesArray := []string{}
	stopMap := make(map[string]string)

	for i, name := range names {
		stopName := name.String()
		namesArray = append(namesArray, stopName)

		nr := nrs[i].String()
		stopMap[stopName] = nr
	}

	return namesArray, stopMap
}

// Get all possibly mentioned lines.
func GetAllLines(file string) []string {
	jsonBytes := ReadFileToBytes(file)
	lines := gjson.GetBytes(jsonBytes, "lines").Array()

	lineArray := []string{}
	for _, line := range lines {
		lineArray = append(lineArray, line.String())
	}
	return lineArray
}

// Get all possible delay words.
func GetAllDelayKeywords(file string) []string {
	jsonBytes := ReadFileToBytes(file)
	delayWords := gjson.GetBytes(jsonBytes, "delayWords").Array()

	wordArray := []string{}
	for _, delayWord := range delayWords {
		wordArray = append(wordArray, delayWord.String())
	}
	return wordArray
}

// Checks if array contains string.
func Contains(arr []string, str string) bool {
	isContaining := false
	for _, a := range arr {
		if a == str {
			isContaining = true
		}
	}
	return isContaining
}
