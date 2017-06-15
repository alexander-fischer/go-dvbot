package main

import (
	"testing"
)

// Run with `go test main_test.go main.go textprocessing.go utils.go textinfo.go answerprocessing.go answer.go -bench=.`

func TestCreateBotAnswer(t *testing.T)  {
	testId := "111"
	testMessage := "Albertplatz 6"
	testAnswer := createBotAnswer(testMessage, testId)

	if testAnswer.category != DEPARTURES {
		t.Error("TestCreateBotAnswer", "has false category")
		t.Error(testAnswer.text)
	}
}
