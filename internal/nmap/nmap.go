package nmap

import (
	"bytes"
	"encoding/xml"
	"io"
	"os/exec"
	"strings"

	nmap_types "github.com/lair-framework/go-nmap"
	"github.com/rs/zerolog/log"
)

func ParseSearchInfo(searchInfo string) []string {
	args := make([]string, 0, 6)

	for _, pair := range strings.Split(searchInfo, ",") {
		kv := strings.SplitN(pair, ":", 2)
		switch kv[0] {
		case "os":
			args = append(args, "-O")
		case "tr":
			args = append(args, "--traceroute")
		case "v":
			args = append(args, "--version-intensity", kv[1])
		case "s":
			args = append(args, "--script", kv[1])
		}
	}

	return args
}

func StartNmapScan(searchInfo string, targets string) {
	cmdArgs := make([]string, 0, 9)
	cmdArgs = append(cmdArgs, "-oX", "-")
	cmdArgs = append(cmdArgs, ParseSearchInfo(searchInfo)...)
	cmdArgs = append(cmdArgs, targets)

	cmd := exec.Command("nmap", cmdArgs...)
	log.Debug().Msgf("Running nmap as: %q", cmd.String())

	// Store errors in buffer to keep after program exits
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Error().Msg("Error getting stdout pipe")
	}

	if err := cmd.Start(); err != nil {
		log.Error().Msg("Error starting nmap scan")
	}

	decoder := xml.NewDecoder(stdout)
	var hosts []nmap_types.Host

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			log.Debug().Msg("Reached end of XML output stream")
			break
		}

		if startElem, ok := token.(xml.StartElement); ok {
			if startElem.Name.Local == "host" {
				var host nmap_types.Host
				if err := decoder.DecodeElement(&host, &startElem); err != nil {
					log.Error().Err(err).Msg("Error decoding host")
					continue
				}
				hosts = append(hosts, host)
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		if strings.Contains(stderr.String(), "requires root") {
			log.Error().Msg("Sudo privileges are required for the specified nmap scan")
		}
	}

	for _, host := range hosts {
		log.Info().Msgf("Host: %+v", host.Addresses)
		for _, port := range host.Ports {
			log.Info().Msgf("Port: %d, State: %s", port.PortId, port.State.State)
		}
	}

	log.Debug().Msg("Finished nmap scan")
}
