package main

import (
	"testing"
)

// Run with `go test main_test.go main.go textprocessing.go utils.go textinfo.go answerprocessing.go answer.go -bench=.`

func TestCreateBotAnswer(t *testing.T)  {
	testId := "111"
	testMessage := "Caspar-David-Friedrich Straße"
	testAnswer, textInfo := createBotAnswer(testMessage, testId)

	if testAnswer.category != 0 {
		t.Error("TestCreateBotAnswer", "has false category")
		t.Error("Answer:", testAnswer)
		t.Error("TextInfo", textInfo)
	}
}
