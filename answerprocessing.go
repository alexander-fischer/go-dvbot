package main

import (
	"github.com/kiliankoe/dvbgo"
	"strconv"
	"net/http"
	"io/ioutil"
	"github.com/tidwall/gjson"
	"time"
	"fmt"
	"strings"
)

// Process answer.
func ProcessAnswer(info TextInfo) Answer {
	answer := Answer{}
	answerText := ""

	switch info.category {
	case DEPARTURES:
		answerText = processDepartures(info)
		break
	case CONNECTIONS:
		answerText = processConnections(info)
		break
	case DELAYS:
		answerText = processDelays(info)
		break
	default:
		answerText = "Das habe ich nicht verstanden. Versuche es erneut."
		break
	}

	if len(answerText) > 640 {
		sentences := strings.Split(answerText, "\n")

		newText := ""
		for i, sentence := range sentences {
			sentence = sentence + "\n"

			if len(newText+sentence) > 640 {
				answer.text = append(answer.text, newText)
				newText = sentence
			} else {
				newText = newText + sentence

				if i == len(sentences)-1 {
					answer.text = append(answer.text, newText)
				}
			}
		}
	} else {
		answer.text = append(answer.text, answerText)
	}
	answer.category = info.category
	return answer
}

// Process text for departures.
func processDepartures(info TextInfo) string {
	stopName := ""
	for _, v := range info.stops {
		stopName = v
	}

	answerText := ""
	deps := getDepartures(stopName)
	for _, dep := range deps {
		minutes := strconv.Itoa(dep.RelativeTime)
		if len(info.lines) > 0 {
			for _, line := range info.lines {
				if line == dep.Line {
					answerText = answerText + "Die Linie " + dep.Line + " Richtung " +
						dep.Direction + " in " + minutes + " Minuten.\n"
				}
			}
		} else {
			answerText = answerText + "Die Linie " + dep.Line + " Richtung " + dep.Direction + " in " +
				minutes + " Minuten.\n"
		}
	}
	return answerText
}

// Get departures with help of stop name.
func getDepartures(stopName string) []*dvb.Departure {
	// Ugly Albertplatz workaround
	if stopName == "albertplatz" {
		stopName = "ALP"
	}

	deps, _ := dvb.Monitor(stopName, 0, "")
	return deps
}

// Process delays.
func processDelays(info TextInfo) string {
	answerText := ""
	delays := getDelays()
	if len(info.lines) > 0 {
		for _, delay := range delays {
			for _, line := range info.lines {
				if Contains(delay.lines, line) {
					answerText = answerText + delay.text + "\n"
				}
			}
		}
	} else {
		for _, delay := range delays {
			answerText = answerText + delay.text + "\n"
		}
	}

	if answerText == "" {
		answerText = "Zurzeit sind keine Verspätungen bekannt."
	}

	return answerText
}

// Get all delays from alexfi server.
func getDelays() []Delay {
	resp, err := http.Get("http://alexfi.dubhe.uberspace.de/text.json")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	delayTexts := gjson.GetBytes(body, "#.text").Array()
	delayLines := gjson.GetBytes(body, "#.linien").Array()
	delayDates := gjson.GetBytes(body, "#.created_at").Array()

	todayDelays := []Delay{}
	for i, delay := range delayDates {
		t, _ := time.Parse(time.RubyDate, delay.String())

		delayWeekDay := t.Weekday()

		loc, _ := time.LoadLocation("Europe/Berlin")
		nowWeekDay := time.Now().In(loc).Weekday()

		if nowWeekDay == delayWeekDay {
			text := delayTexts[i].String()
			lines := []string{}
			lineResults := delayLines[i].Array()
			for _, l := range lineResults {
				lines = append(lines, l.String())
			}

			d := Delay{t, lines, text}
			todayDelays = append(todayDelays, d)
		}
	}
	return todayDelays
}

// Process connections.
func processConnections(info TextInfo) string {
	answerText := ""
	nrOfResults := 1

	stopArr := []string{}
	for k := range info.stops {
		stopArr = append(stopArr, k)
	}

	cons := GetConnectionsFromDvb(stopArr[0], stopArr[1], nrOfResults)

	if len(cons) == 0 {
		return "Zwischen diesen Haltestellen ist keine Verbindung möglich."
	}

	for i := 0; i <= (nrOfResults - 1); i++ {
		con := cons[i]
		fmt.Println(i, con)

		startPoint := con.legs[0].startPoint
		startTime := con.startTime
		startLine := con.legs[0].line
		startDir := con.legs[0].lineDirection

		answerText = answerText + "Um " + startTime + " fährt die Linie " + startLine + " in Richtung " + startDir +
			" ab " + startPoint + ".\n"

		if con.hasInterchanges {
			for i := 1; i < len(con.legs); i++ {
				leg := con.legs[i]
				answerText = answerText + "An der Haltestelle " + leg.startPoint + " steigst du in die Linie " +
					leg.line + " in Richtung " + leg.lineDirection + " um.\n"
			}
		}

		answerText = answerText + "Die Fahrt dauert " + con.duration + " Minuten und die Ankunftszeit ist " +
			con.endTime + ".\n"
	}

	return answerText
}
