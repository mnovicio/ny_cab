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
		Use:   "ny_cab_client_grpc",
		Short: "connects to localhost:10001 to use NY CAB gRPC APIs",
		Long:  `connects to localhost:10001 to use NY CAB gRPC APIs`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("server", "s", "localhost:10001", "NY CAB gRPC host")
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
