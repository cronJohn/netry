package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Ullaakut/nmap/v3"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cronJohn/netry/types"
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
		log.Info().Msgf("Ports: %q", portsFlag)
		log.Info().Msgf("Targets: %q", targetsFlag)
		log.Info().Msgf("Identity: %q", identityFlag)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		startNmapScan(ctx)
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
p:ratio: Scan a set of ports based on a ratio/percentage (e.g., p:.10 would scan top 10% most commonly used ports).

Examples:
- t:100: Scan the top 100 most common ports
- r:1-1024: Scan the range of ports from 1 to 1024
- p:.10: Scan top 10% most commonly used ports
`)

	rootCmd.Flags().
		StringVarP(&targetsFlag, "targets", "t", "a", `Specify the targets to scan. The format is <type>:value, with the following types:
a|b|c: Scan the private IPv4 address classes a, b, or c.
f:file: Scan the targets listed in the specified file.
n:target|range: Normal target specification.
r:rnd: Scan a random set of targets.

Examples:
- a: Scan the entire class A private IPv4 address space
- f:/path/to/targets.txt: Scan the targets listed in the file /path/to/targets.txt
- n:1.2.3.4: Scan the target 1.2.3.4
- n:10.0.0.0-255: Scan the range of targets from 10.0.0.0 to 10.0.0.255
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

func startNmapScan(ctx context.Context) {
	scanner, err := nmap.NewScanner(
		ctx,
		types.WithCustomPorts(portsFlag),
		types.WithCustomTargets(targetsFlag),
	)
	if err != nil {
		log.Error().Err(err).Msg("Unable to create nmap scanner")
	}

	log.Info().Msgf("Running nmap with args: %s", scanner.Args())

	// TODO: Test later
	// result, warnings, err := scanner.Run()
	// if len(*warnings) > 0 {
	// 	log.Warn().Msg("run finished with warnings")
	// }
	// if err != nil {
	// 	log.Error().Err(err).Msg("run finished with error")
	// }
	//
	// for _, host := range result.Hosts {
	// 	if len(host.Ports) == 0 || len(host.Addresses) == 0 {
	// 		continue
	// 	}
	//
	// 	log.Info().Msgf("Host %q:\n", host.Addresses[0])
	//
	// 	for _, port := range host.Ports {
	// 		log.Info().
	// 			Msgf("\tPort %d/%s %s %s\n", port.ID, port.Protocol, port.State, port.Service.Name)
	// 	}
	// }
}
