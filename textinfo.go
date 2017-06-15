package main

const (
	DEPARTURES  = 1 << iota
	CONNECTIONS
	DELAYS
)

type TextInfo struct {
	text      string
	category  int
	delay     bool
	stopNames []string
	stopIds   []string
	lines     []string
}
