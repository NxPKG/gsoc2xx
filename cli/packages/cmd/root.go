/*
Copyright (c) 2023 Gsoc2 Inc.
*/
package cmd

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/Gsoc2/gsoc2-merge/packages/config"
	"github.com/Gsoc2/gsoc2-merge/packages/telemetry"
	"github.com/Gsoc2/gsoc2-merge/packages/util"
)

var Telemetry *telemetry.Telemetry

var rootCmd = &cobra.Command{
	Use:               "gsoc2",
	Short:             "Gsoc2 CLI is used to inject environment variables into any process",
	Long:              `Gsoc2 is a simple, end-to-end encrypted service that enables teams to sync and manage their environment variables across their development life cycle.`,
	CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
	Version:           util.CLI_VERSION,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLog)
	rootCmd.PersistentFlags().StringP("log-level", "l", "info", "log level (trace, debug, info, warn, error, fatal)")
	rootCmd.PersistentFlags().Bool("telemetry", true, "Gsoc2 collects non-sensitive telemetry data to enhance features and improve user experience. Participation is voluntary")
	rootCmd.PersistentFlags().StringVar(&config.GSOC2_URL, "domain", util.GSOC2_DEFAULT_API_URL, "Point the CLI to your own backend [can also set via environment variable name: GSOC2_API_URL]")
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if !util.IsRunningInDocker() {
			util.CheckForUpdate()
		}
	}

	// if config.GSOC2_URL is set to the default value, check if GSOC2_URL is set in the environment
	// this is used to allow overrides of the default value
	if !rootCmd.Flag("domain").Changed {
		if envGsoc2BackendUrl, ok := os.LookupEnv("GSOC2_API_URL"); ok {
			config.GSOC2_URL = envGsoc2BackendUrl
		}
	}

	isTelemetryOn, _ := rootCmd.PersistentFlags().GetBool("telemetry")
	Telemetry = telemetry.NewTelemetry(isTelemetryOn)
}

func initLog() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	ll, err := rootCmd.Flags().GetString("log-level")
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	switch strings.ToLower(ll) {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "err", "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
