package main

import (
	"testing"
)

// Run with `go test utils_test.go utils.go -bench=.`

var filename = "json/test.json"
var stopfile = "json/transit_stops.json"
var keywordFile = "json/keywords.json"

func BenchmarkReadFileToBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ReadFileToBytes(filename)
	}
}

func TestGetMessageFromRequest(t *testing.T) {
	expectedMessage := "Verspätungen der Linie 3 und 10 ab Hauptbahnhof"

	jsonBytes := ReadFileToBytes(filename)
	resultMessage, _ := GetMessageFromRequest(jsonBytes)

	if resultMessage != expectedMessage {
		t.Error("String was not cleaned")
	}
}

func BenchmarkGetMessageFromRequest(b *testing.B) {
	jsonBytes := ReadFileToBytes(filename)
	for i := 0; i < b.N; i++ {
		GetMessageFromRequest(jsonBytes)
	}
}

func TestGetAllTransitStops(t *testing.T) {
	namesArray, resultMap := GetAllTransitStops(stopfile)

	if len(namesArray) <= 0 {
		t.Error("Transit stop map is equals or smaller as 1")
	}
	if len(resultMap) <= 0 {
		t.Error("Transit stop map is equals or smaller as 1")
	}

	expectedId := "de:14612:834"
	resultId := resultMap["steglichstraße"]
	if resultId != expectedId {
		t.Error("Stop id does not fit")
	}
}

func BenchmarkGetAllTransitStops(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetAllTransitStops(stopfile)
	}
}

func TestGetAllLines(t *testing.T) {
	lines := GetAllLines(keywordFile)

	if len(lines) <= 0 {
		t.Error("Error GetAllLines")
	}
}

func BenchmarkGetAllLines(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetAllLines(keywordFile)
	}
}

func TestGetAllDelayKeywords(t *testing.T) {
	delayWords := GetAllDelayKeywords(keywordFile)

	if len(delayWords) <= 0 {
		t.Error("Error GetAllLines")
	}
}

func BenchmarkGetAllDelayKeywords(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetAllDelayKeywords(keywordFile)
	}
}
