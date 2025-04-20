package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var (
	portsFlag    string
	targetsFlag  string
	identityFlag string
)

var rootCmd = &cobra.Command{
	Use:   "netry",
	Short: "A CLI tool that attempts to build a network topology visualization",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("Starting netry")
		log.Info().Msgf("Ports: %s", portsFlag)
		log.Info().Msgf("Targets: %s", targetsFlag)
		log.Info().Msgf("Identity: %s", identityFlag)
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
		StringVarP(&portsFlag, "ports", "p", "t:100", `Specify the ports to scan. The format is <type>:value, with the following types:
t:num: Scan the top num most common ports.
r:start-end: Scan the range of ports from start to end (inclusive).
/:ratio: Scan a set of ports based on a ratio (e.g., /10 would scan 10% of the ports).

Examples:
- t:100: Scan the top 100 most common ports
- r:1-1024: Scan the range of ports from 1 to 1024
- /10: Scan 10% of the available ports
`)

	rootCmd.Flags().
		StringVarP(&targetsFlag, "targets", "t", "a", `Specify the targets to scan. The format is <type>:value, with the following types:
a|b|c: Scan the private IPv4 address classes a, b, or c.
f:file: Scan the targets listed in the specified file.
n:start-end: Scan the range of targets from start to end (inclusive).
r:rnd: Scan a random set of targets.

Examples:
- a: Scan the entire class A private IPv4 address space
- b:10-20: Scan the range of class B private IPv4 addresses from 10 to 20
- f:/path/to/targets.txt: Scan the targets listed in the file /path/to/targets.txt
- r:100: Scan 100 random targets
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
