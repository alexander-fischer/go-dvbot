package main

import (
	"testing"
)

// Run with `go test dvbwebservice_test.go dvbwebservice.go utils.go -bench=.`

func TestCreateBodyForConnection(t *testing.T) {
	body := createBodyForConnection("de:14612:17", "de:14612:27", 1)
	if body == "" {
		t.Error("TestCreateBodyForConnection is null")
	}
}

func BenchmarkCreateBodyForConnection(b *testing.B) {
	for i := 0; i < b.N; i++ {
		createBodyForConnection("de:14612:17", "de:14612:27", 1)
	}
}

func TestCallDvbApi(t *testing.T) {
	body := createBodyForConnection("de:14612:17", "de:14612:27", 1)
	result := callDvbApi(body)
	if result == "" {
		t.Error("TestCallDvbApi is null")
	}
}

func BenchmarkCallDvbApi(b *testing.B) {
	body := createBodyForConnection("de:14612:17", "de:14612:27", 1)
	for i := 0; i < b.N; i++ {
		callDvbApi(body)
	}
}

func TestGetConnectionsFromDvb(t *testing.T) {
	cons := GetConnectionsFromDvb("de:14612:17", "de:14612:27", 1)

	if len(cons) != 1 {
		t.Error("TestGetConnectionsFromDvb has not 1 Connection")
	}
}

func BenchmarkGetConnectionsFromDvb(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetConnectionsFromDvb("de:14612:17", "de:14612:27", 1)
	}
}

func TestProcessRouteResults(t *testing.T) {
	body := createBodyForConnection("de:14612:17", "de:14612:27", 1)
	result := callDvbApi(body)
	cons := processRouteResults(result)

	if len(cons) != 1 {
		t.Error("TestProcessRouteResults has not 1 Connection")
	}
}

func BenchmarkProcessRouteResults(b *testing.B) {
	body := createBodyForConnection("de:14612:17", "de:14612:27", 1)
	result := callDvbApi(body)

	for i := 0; i < b.N; i++ {
		processRouteResults(result)
	}
}