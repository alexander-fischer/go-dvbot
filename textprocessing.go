package main

import (
	"strings"
	"github.com/antzucaro/matchr"
)

var allNames, allStops = GetAllTransitStops("json/transit_stops.json")

func ProcessText(text string) ([]string, []string, []string, bool) {
	stopNames := []string{}
	stopIds := []string{}
	linesArray := []string{}
	isAboutDelay := false

	editText := strings.ToLower(text)
	editText = CleanString(editText)
	wordArray := strings.Split(editText, " ")

	// Search for stops.
	for i, word := range wordArray {
		wordIsLine := isLine(word)
		if wordIsLine {
			linesArray = append(linesArray, word)
		}

		wordIsDelayWord := isDelayWord(word)
		if wordIsDelayWord {
			isAboutDelay = true
		}

		isLast := false
		words := []string{}

		if i != 0 {
			words = append(words, wordArray[i-1])
		}
		words = append(words, word)

		// If is last word
		if i != len(wordArray)-1 {
			words = append(words, wordArray[i+1])
		} else {
			isLast = true
		}

		stopName, stopNr := FindStops(words, isLast)
		if stopName != "" && stopNr != "" {
			if !Contains(stopIds, stopNr) {
				stopNames = append(stopNames, stopName)
				stopIds = append(stopIds, stopNr)
			}
		}
	}

	return stopNames, stopIds, linesArray, isAboutDelay
}

// Clean the word from useless chars.
func CleanString(text string) string {
	text = strings.Replace(text, ".", "", -1)
	text = strings.Replace(text, ",", "", -1)
	text = strings.Replace(text, ";", "", -1)
	text = strings.Replace(text, ":", "", -1)
	text = strings.Replace(text, "?", "", -1)
	text = strings.Replace(text, "!", "", -1)

	return text
}

// Processed different word possibilities.
func FindStops(words []string, isLast bool) (string, string) {
	stopName, stopNr := "", ""

	switch len(words) {
	case 3:
		// Search for all 3 words.
		searchedString := words[0] + " " + words[1] + " " + words[2]
		stopName, stopNr = search(searchedString)
		if stopName != "" {
			return stopName, stopNr
		}

		// Search for first 2 words.
		searchedString = words[0] + " " + words[1]
		stopName, stopNr = search(searchedString)
		if stopName != "" {
			return stopName, stopNr
		}

		// Search for last 2 words.
		searchedString = words[1] + " " + words[2]
		stopName, stopNr = search(searchedString)
		if stopName != "" {
			return stopName, stopNr
		}

		// Search original word.
		searchedString = words[1]
		stopName, stopNr = search(searchedString)
		if stopName != "" {
			return stopName, stopNr
		}
	case 2:
		// Searches all 2 words.
		searchedString := words[0] + " " + words[1]
		stopName, stopNr = search(searchedString)
		if stopName != "" {
			return stopName, stopNr
		}

		if !isLast {
			// Searches for the first word.
			stopName, stopNr = search(words[0])
			if stopName != "" {
				return stopName, stopNr
			}
		} else {
			// Searches for the last word.
			stopName, stopNr = search(words[1])
			if stopName != "" {
				return stopName, stopNr
			}
		}
	case 1:
		stopName, stopNr = search(words[0])
		if stopName != "" {
			return stopName, stopNr
		}
	default:
		return stopName, stopNr

	}

	return stopName, stopNr
}

// Searches the stop of the list with help of Damerau-Levenshtein algorithm.
func search(wordGroup string) (string, string) {
	threshold := 0.93

	matchedNames := []Rank{}

	for _, name := range allNames {
		dist := matchr.DamerauLevenshtein(wordGroup, name)
		if dist <= 5 {
			d := float64(dist)
			l := float64(len(wordGroup + name))
			ratio := 1 - float64(d/l)
			if ratio >= threshold {
				rank := Rank{ratio, name}
				matchedNames = append(matchedNames, rank)
			}
		}
	}

	if len(matchedNames) > 0 {
		bestRank := findBestRank(matchedNames)
		stopName := bestRank.Name
		stopNr := allStops[stopName]

		return stopName, stopNr
	} else {
		return "", ""
	}
}

// Categorize the text for answer creation.
// TODO tests needed
func Categorize(info TextInfo) int {
	if info.delay {
		return DELAYS
	}

	switch len(info.stopIds) {
	case 0:
		break
	case 1:
		return DEPARTURES
	case 2:
		return CONNECTIONS
	default:
		break
	}

	return 0
}

// Checks if word is line.
func isLine(word string) bool {
	lines := GetAllLines("json/keywords.json")
	for i := 0; i < len(lines); i++ {
		if lines[i] == word {
			return true
		}
	}
	return false
}

// Checks if is delay word.
func isDelayWord(word string) bool {
	delayWords := GetAllDelayKeywords("json/keywords.json")
	for i := 0; i < len(delayWords); i++ {
		if delayWords[i] == word {
			return true
		}
	}
	return false
}

type Rank struct {
	RelDist float64
	Name    string
}

// Find best rank.
func findBestRank(ranks []Rank) Rank {
	bestDist := 0.0
	bestRank := Rank{}

	for _, rank := range ranks {
		rankDist := rank.RelDist
		if rankDist > bestDist {
			bestRank = rank
		}
	}

	return bestRank
}
