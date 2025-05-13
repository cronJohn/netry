package nmap

import (
	"strings"
)

func ParseSearchInfo(searchInfo string) string {
	baseString := strings.Builder{}
	baseString.Grow(55)

	for _, pair := range strings.Split(searchInfo, ",") {
		kv := strings.SplitN(pair, ":", 2)
		switch kv[0] {
		case "os":
			baseString.WriteString("-O ")
		case "tr":
			baseString.WriteString("--traceroute ")
		case "v":
			baseString.WriteString("--version-intensity " + kv[1] + " ")
		case "s":
			baseString.WriteString("--script=" + kv[1] + " ")
		}
	}

	return baseString.String()
}

func StartNmapScan(targets string, searchInfo string, customArgs string) {
	var nmapCommand *strings.Builder

	_ = nmapCommand
}
