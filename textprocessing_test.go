package main

import (
	"testing"
	"strings"
)

// Run with `go test textprocessing_test.go textprocessing.go utils.go textinfo.go -bench=.`
var testStr = "Abfahrten. der, Linie? 3! ab Hauptbahnhof;"
var expectedStr = "Abfahrten der Linie 3 ab Hauptbahnhof"

func TestCleanString(t *testing.T) {
	result := CleanString(testStr)
	if result != expectedStr {
		t.Error("String was not cleaned")
	}
}

func BenchmarkCleanWord(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CleanString(testStr)
	}
}

func TestFindStops(t *testing.T) {
	expectedStr = strings.ToLower(expectedStr)
	wordArray := strings.Split(expectedStr, " ")
	stopName, stopNr := "", ""

	for id, word := range wordArray {
		isLast := false
		words := []string{}

		if id != 0 {
			words = append(words, wordArray[id-1])
		}
		words = append(words, word)

		if id != len(wordArray)-1 {
			words = append(words, wordArray[id+1])
		} else {
			isLast = true
		}
		stopName, stopNr = FindStops(words, isLast)
	}

	if stopName == "" && stopNr == "" {
		t.Error("Can't find stop")
	}
}

func BenchmarkFindStops(b *testing.B) {
	for i := 0; i < b.N; i++ {
		expectedStr = strings.ToLower(expectedStr)
		wordArray := strings.Split(expectedStr, " ")
		for id, word := range wordArray {
			isLast := false
			words := []string{}

			if id != 0 {
				words = append(words, wordArray[id-1])
			}
			words = append(words, word)

			if id != len(wordArray)-1 {
				words = append(words, wordArray[id+1])
			} else {
				isLast = true
			}
			FindStops(words, isLast)
		}
	}
}

func TestFuzzySearch(t *testing.T) {
	wordGroup := "hauptbahnho"
	stopName, stopNr := fuzzySearch(wordGroup)

	if stopName != "hauptbahnhof" || stopNr != "de:14612:28" {
		t.Error("Can't find stop")
	}
}

func BenchmarkFuzzySearch(b *testing.B) {
	wordGroup := "hauptbahnho"
	for i := 0; i < b.N; i++ {
		fuzzySearch(wordGroup)
	}
}

func TestIsLine(t *testing.T) {
	word1 := "3"
	word2 := "kebab"

	if !isLine(word1) {
		t.Error("TestIsLine: Should be true.")
	}

	if isLine(word2) {
		t.Error("TestIsLine: Should be false.")
	}
}

func BenchmarkIsLine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isLine("61")
	}
}

func TestIsDelayWord(t *testing.T) {
	word1 := "verspätungen"
	word2 := "kebab"

	if !isDelayWord(word1) {
		t.Error("TestIsDelayWord: Should be true.")
	}

	if isDelayWord(word2) {
		t.Error("TestIsDelayWord: Should be false.")
	}
}

func BenchmarkIsDelayWord(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isLine("verspätungen")
	}
}
