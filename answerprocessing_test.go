package main

import (
	"testing"
)

// Run with `go test answerprocessing_test.go answerprocessing.go answer.go textinfo.go utils.go -bench=.`

func BenchmarkGetDepartures(b *testing.B)  {
	for i := 0; i < b.N; i++ {
		getDepartures("hauptbahnhof")
	}
}

func BenchmarkGetDelays(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getDelays()
	}
}