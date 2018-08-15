package kache

import (
	"fmt"
	"github.com/kasvith/kache/internal/config"
	"github.com/kasvith/kache/internal/errh"
	"github.com/kasvith/kache/internal/srv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var verbose bool
var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "kache",
	Short: "kache is a simple distributed in memory database",
	Long:  `A fast and flexible redis alternative built with go`,
	Run:   runApp,
}

func init() {
	cobra.OnInitialize(initConfig)

	// Flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "", "", "configuration file")
	rootCmd.Flags().StringP("host", "", "127.0.0.1", "host for running application")
	rootCmd.Flags().IntP("port", "p", 6969, "port for running application")
	rootCmd.Flags().IntP("maxClients", "", 10000, "max connections can be handled")
	rootCmd.Flags().IntP("maxTimeout", "", 120, "max timeout for clients(in seconds)")

	// Bind the flags to config
	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))
	viper.BindPFlag("host", rootCmd.Flags().Lookup("host"))
	viper.BindPFlag("maxClients", rootCmd.Flags().Lookup("maxClients"))
	viper.BindPFlag("maxTimeout", rootCmd.Flags().Lookup("maxTimeout"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	loadingDefault := false
	viper.SetConfigType("toml")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		loadingDefault = true
		viper.SetConfigName(".default")
		viper.AddConfigPath("../config")
	}

	if err := viper.ReadInConfig(); err != nil {
		if !loadingDefault {
			fmt.Fprintf(os.Stderr, "Error reading config from %s : %s\n", viper.ConfigFileUsed(), err)
			os.Exit(1)
		}

		fmt.Printf("%T\n", err)

		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			fmt.Fprintf(os.Stderr, "Error reading config file from config directory\nLoading with application defaults...\n")
			break

		default:
			fmt.Fprintf(os.Stderr, "Error reading config in %s : %s\n", viper.ConfigFileUsed(), err)
			os.Exit(1)
		}
	}
}

func Execute() {
	// Commands
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runApp(cmd *cobra.Command, args []string) {
	var appConfig config.AppConfig
	if err := viper.Unmarshal(&appConfig); err != nil {
		errh.PrintErrorAndExit(err, 2)
	}

	srv.Start(appConfig)
}
