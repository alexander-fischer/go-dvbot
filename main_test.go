package main

import (
	"testing"
)

// Run with `go test main_test.go main.go textprocessing.go utils.go textinfo.go answerprocessing.go answer.go -bench=.`


func TestMain(m *testing.M) {
	main()
}

func BenchmarkMain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		main()
	}
}
