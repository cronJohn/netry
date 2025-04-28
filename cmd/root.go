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
)

var cfgFile string

var nmapCustomArgsFlag string

var rootCmd = &cobra.Command{
	Use:   "netry",
	Short: "A CLI tool that attempts to build a network topology visualization",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("Starting netry")
		log.Info().Msgf("Nmap args: %q", nmapCustomArgsFlag)

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
		StringVarP(&nmapCustomArgsFlag, "nargs", "n", "-sS 10.0.0.0/8", `Custom nmap arguments:
If you know nmap, you can use this to specify
custom arguments and flags to nmap.
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
	_, err := nmap.NewScanner(
		ctx,
		/* Using custom arguments instead of individual With* statements to
		avoid becoming a wrapper for a Nmap wrapper. This keeps Nmap
		functionality separate and allows users to leverage its full capabilities directly.
		*/
		nmap.WithCustomArguments(nmapCustomArgsFlag),
	)
	if err != nil {
		log.Error().Err(err).Msg("Unable to create nmap scanner")
	}

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
