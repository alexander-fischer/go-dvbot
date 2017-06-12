package main

import (
	"net/http"
	"strings"
	"fmt"
	"strconv"
	"time"
	xj "github.com/basgys/goxml2json"
	"github.com/tidwall/gjson"
	"regexp"
	"log"
)

const (
	url      string = "http://trias.vvo-online.de:9000/Middleware/Data/Trias"
	routeXml string = "xml/route.xml"
)

// Calls the TRIAS DVB backend.
func callDvbApi(body string) string {
	bodyBytes := strings.NewReader(body)

	req, err := http.NewRequest("POST", url, bodyBytes)
	if err != nil {
		fmt.Printf("http.NewRequest() error: %v\n", err)
		return ""
	}

	req.Header.Add("Cookie", "cookie")
	req.Header.Add("Content-Type", "application/xml")

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		log.Fatal("http.Do() error:", err)
		return ""
	}

	defer resp.Body.Close()

	if err != nil {
		log.Fatal("ioutil.ReadAll() error:", err)
		return ""
	}

	resJson, err := xj.Convert(resp.Body)
	if err != nil {
		log.Fatal("xj.Convert() error:", err)
		return ""
	}
	return resJson.String()
}

// Return connections between two points.
func GetConnectionsFromDvb(origin string, destination string, nrOfResults int) []Connection {
	body := createBodyForConnection(origin, destination, nrOfResults)
	resString := callDvbApi(body)

	cons := processRouteResults(resString)
	return cons
}

// Create body for route request.
func createBodyForConnection(origin string, destination string, nrOfResults int) string {
	routeXmlStr := string(ReadFileToBytes(routeXml))

	strResults := strconv.Itoa(nrOfResults)
	loc, _ := time.LoadLocation("Europe/Berlin")
	timeNow := "<DepArrTime>" + time.Now().In(loc).Format("2006-01-02T15:04:05") + "</DepArrTime>"

	body := routeXmlStr
	body = strings.Replace(body, "[[origin]]", origin, -1)
	body = strings.Replace(body, "[[destination]]", destination, 1)
	body = strings.Replace(body, "[[result_number]]", strResults, 1)
	body = strings.Replace(body, "[[destination_time]]", "", 1)
	body = strings.Replace(body, "[[origin_time]]", timeNow, 1)

	return body
}

// Processes a JSON string to an array of connections.
func processRouteResults(json string) []Connection {
	connections := []Connection{}
	results := gjson.Get(json, "Trias.ServiceDelivery.DeliveryPayload.TripResponse.TripResult")

	if !results.Exists() {
		fmt.Println("Could not find trip Results")
		return connections
	}

	resultType := gjson.Get(json, "Trias.ServiceDelivery.DeliveryPayload.TripResponse.TripResult.#")

	// Checks if is array or not.
	if resultType.Type == gjson.Number {
		results.ForEach(func(_, value gjson.Result) bool {
			connections = append(connections, parseRouteResults(value))
			return true
		})
	} else {
		connections = append(connections, parseRouteResults(results))
	}

	return connections
}

// Parses JSON to create connection.
func parseRouteResults(result gjson.Result) Connection {
	con := Connection{}
	trip := result.Get("Trip")

	// duration
	reg, err := regexp.Compile("[A-Z]+")
	if err != nil {
		log.Fatal(err)
		return con
	}
	con.duration = reg.ReplaceAllString(trip.Get("Duration").String(), "")

	// startTime
	startTime := trip.Get("StartTime").String()
	startTimeParsed, err := time.Parse("2006-01-02T15:04:05", startTime)
	if err != nil {
		log.Fatal(err)
		return con
	}
	con.startTime = startTimeParsed.Format("15:04")

	// endTime
	endTime := trip.Get("EndTime").String()
	endTimeParsed, err := time.Parse("2006-01-02T15:04:05", endTime)
	if err != nil {
		log.Fatal(err)
		return con
	}
	con.endTime = endTimeParsed.Format("15:04")

	// hasInterchanges
	legs := trip.Get("TripLeg")
	hasInterchanges := trip.Get("TripLeg.#").Type == gjson.Number
	con.hasInterchanges = hasInterchanges

	// legs
	legsArray := []Leg{}
	if hasInterchanges {
		legs.ForEach(func(_, leg gjson.Result) bool {
			legsArray = append(legsArray, processLegs(leg))
			return true
		})
	} else {
		legsArray = append(legsArray, processLegs(legs))
	}
	con.legs = legsArray

	return con
}

// Parses JSON to create Leg.
func processLegs(leg gjson.Result) Leg {
	l := Leg{}

	startPoint := leg.Get("TimedLeg.LegBoard.StopPointName.Text")
	l.startPoint = startPoint.String()

	endPoint := leg.Get("TimedLeg.LegAlight.StopPointName.Text")
	l.endPoint = endPoint.String()

	line := leg.Get("TimedLeg.Service.PublishedLineName.Text")
	l.line = line.String()

	lineDir := leg.Get("TimedLeg.Service.DestinationText.Text")
	l.lineDirection = lineDir.String()

	return l
}

type Connection struct {
	duration        string
	startTime       string
	endTime         string
	legs            []Leg
	hasInterchanges bool
}

type Leg struct {
	startPoint    string
	endPoint      string
	line          string
	lineDirection string
}
