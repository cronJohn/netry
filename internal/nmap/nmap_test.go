package nmap

import (
	"testing"
)

const (
	testSearchInfo = "os,tr,v:2,s:default"
)

func TestParseSearchInfo(t *testing.T) {
	got := ParseSearchInfo(testSearchInfo)
	want := "-O --traceroute --version-intensity 2 --script=default "

	if got != want {
		t.Errorf("Got %q, Want %q", got, want)
	}
}

func BenchmarkParseSearchInfo(b *testing.B) {
	for b.Loop() {
		_ = ParseSearchInfo(testSearchInfo)
	}
}
