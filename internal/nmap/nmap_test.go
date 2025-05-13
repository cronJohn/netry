package nmap

import (
	"testing"
)

const (
	searchInfo = "os,version:2,script:default,traceroute"
)

func BenchmarkParseSearchInfo(b *testing.B) {
	for b.Loop() {
		_ = ParseSearchInfo(searchInfo)
	}
}
