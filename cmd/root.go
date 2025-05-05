package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Ullaakut/nmap/v3"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var (
	scanTargets  string
	nmapBehavior string
)

var rootCmd = &cobra.Command{
	Use:   "netry",
	Short: "A CLI tool that attempts to build a network topology visualization",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("Starting netry")
		log.Debug().Msgf("Scanning targets: %q", scanTargets)
		log.Debug().Msgf("Nmap behavior: %q", nmapBehavior)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		scanResults := startNmapScan(ctx)

		for _, host := range scanResults.Hosts {
			log.Info().Msgf("Host ports: %+v", host.Ports)
			log.Info().Msgf("Address: '%+v'", host.Addresses)
		}
	},
}

func startNmapScan(ctx context.Context) *nmap.Run {
	scanner, err := nmap.NewScanner(
		ctx,
		WithCustomNmapBehavior(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Unable to create nmap scanner")
	}

	log.Debug().Msgf("Nmap args string: %+v", scanner.Args())

	result, warnings, err := scanner.Run()
	if len(*warnings) > 0 {
		log.Warn().Msgf("Run finished with warnings: %v", *warnings)
	}
	if err != nil {
		log.Error().Err(err).Msg("Run finished with error")
	}

	return result
}

func WithCustomNmapBehavior() nmap.Option {
	options := make([]nmap.Option, 0, 2)

	switch nmapBehavior {
	// -- Define scan modes --
	case "discovery":
		options = append(options, nmap.WithPingScan())
	case "full":
		options = append(options, nmap.WithAggressiveScan())
	case "os":
		options = append(options, nmap.WithOSDetection())
	case "traceroute":
		options = append(options, nmap.WithTraceRoute())

	// -- otherwise, treat it as a custom nmap arguments string --
	default:
		options = append(options, nmap.WithCustomArguments(strings.Fields(nmapBehavior)...))
	}

	options = append(options, nmap.WithTargets(scanTargets))

	return func(scanner *nmap.Scanner) {
		for _, option := range options {
			option(scanner)
		}
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().
		StringVar(&cfgFile, "config", "", "config file (default is $HOME/.netry.yaml)")

	rootCmd.Flags().
		StringVarP(&scanTargets, "targets", "t", "localhost", `Targets to scan:
Same target specification as nmap.

Examples:
  --targets a,b,c (Scans targets a, b, and c)
  --targets 10.0.0.1/24 (Scans using CIDR notation)
  --targets 10.0.0.1-255 (Scans using range notation)
`)
	rootCmd.Flags().
		StringVarP(&nmapBehavior, "nmap", "n", "-sS", `Nmap behavior:
Can either be a specific mode or a string of custom nmap arguments and flags

Examples:
  --nmap "discover" (Performs a discovery scan, i.e. -sn)
  --nmap "full" (Performs a full scan, i.e. -A)
  --nmap "os" (Performs an OS detection scan, i.e. -O)
  --nmap "traceroute" (Performs a traceroute scan, i.e. --traceroute)
  --nmap "-sS -p 50-100" (Performs a SYN scan for ports 50-100)
`)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".netry" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".netry")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
