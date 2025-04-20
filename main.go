package main

import (
	"github.com/rs/zerolog/log"

	"github.com/cronJohn/netry/cmd"
	_ "github.com/cronJohn/netry/pkg/logger"
)

func main() {
	log.Info().Msg("Starting netry")
	cmd.Execute()
}
