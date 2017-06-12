package main

const (
	DEPARTURES = 1 << iota
	CONNECTIONS
	DELAYS
)

type TextInfo struct {
	text     string
	category int
	delay    bool
	stops    map[string]string
	lines    []string
}
