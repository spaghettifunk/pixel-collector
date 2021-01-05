package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	logDebugFlag = "debug"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pixel",
	Short: "",
	Long:  ``,
}

// Execute will start the application
func Execute() {
	cobra.OnInitialize(initConfig)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}

// initConfig sets AutomaticEnv in viper to true.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match

	r := rootCmd.PersistentFlags()
	r.Bool(logDebugFlag, false, "sets log level to debug")

	viper.BindEnv(logDebugFlag, "LOG_DEBUG")
}
