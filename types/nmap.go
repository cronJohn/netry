package types

import (
	"strconv"

	"github.com/Ullaakut/nmap/v3"
)

func WithCustomPorts(ports string) nmap.Option {
	switch portType := ports[0]; portType {
	case 't':
		topPorts, err := strconv.Atoi(ports[2:])
		if err != nil {
			return nil
		}
		foo := nmap.WithMostCommonPorts(topPorts)
		return foo

	case 'r':
		portRange := ports[2:]
		return nmap.WithPorts(portRange)

	case 'p':
		ratioBuf, err := strconv.ParseFloat(ports[2:], 32)
		if err != nil {
			return nil
		}
		// convert 64 float buffer to 32 bit float
		ratio := float32(ratioBuf)
		return nmap.WithPortRatio(ratio)

	default:
		return nil
	}
}

func WithCustomTargets(targets string) nmap.Option {
	switch targetType := targets[0]; targetType {
	case 'a', 'b', 'c':
		return nmap.WithTargets(map[byte]string{
			'a': "10.0.0.0/8",
			'b': "172.16.0.0/12",
			'c': "192.168.0.0/16",
		}[targetType])
	case 'f':
		return nmap.WithTargetInput(targets[2:])
	case 'n':
		return nmap.WithTargets(targets[2:])
	case 'r':
		targetAmount, err := strconv.Atoi(targets[2:])
		if err != nil {
			return nil
		}
		return nmap.WithRandomTargets(targetAmount)
	default:
		return nil
	}
}
