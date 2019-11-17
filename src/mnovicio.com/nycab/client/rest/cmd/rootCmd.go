package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "ny_cab_client_rest",
		Short: "connects to localhost:10002 to use NY CAB REST endpoints",
		Long:  `connects to localhost:10002 to use NY CAB REST endpoints`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("server", "s", "http://localhost:10002", "NY CAB service host")
	viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath("./")
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func trackTime(start time.Time, name string) {
	elapsed := time.Since(start)

	log.Printf("%s took %s", name, elapsed)
}
