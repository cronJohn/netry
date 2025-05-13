package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cronJohn/netry/internal/nmap"
)

var cfgFile string

var (
	scanTargets string
	targetInfo  string
)

var rootCmd = &cobra.Command{
	Use:   "netry",
	Short: "A CLI tool that attempts to build a network topology visualization",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("Starting netry")
		log.Debug().Msgf("Scanning targets: %q", scanTargets)
		log.Debug().Msgf("Info to find: %q", targetInfo)

		nmap.StartNmapScan(targetInfo, scanTargets)
	},
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
		StringVarP(&targetInfo, "info", "i", "", `Info to find about the targets:
Contains a comma-separated list of info types to find about the targets.

Examples:
  --info os,tr (Finds OS and traceroute information)
  --info s:default (Finds nmap default script information)
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
