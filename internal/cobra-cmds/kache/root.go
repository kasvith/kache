/*
 * MIT License
 *
 * Copyright (c)  2018 Kasun Vithanage
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package kache

import (
	"fmt"
	"github.com/kasvith/kache/internal/config"
	"github.com/kasvith/kache/internal/klogs"
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
	Long:  `A fast and flexible in memory database built with go`,
	Run:   runApp,
}

func init() {
	cobra.OnInitialize(initConfig)

	// Flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "output debug information")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "", "", "configuration file")
	rootCmd.PersistentFlags().BoolP("logging", "", true, "set application logs")
	rootCmd.PersistentFlags().StringP("logfile", "", "", "application log file")

	rootCmd.Flags().StringP("host", "", "127.0.0.1", "host for running application")
	rootCmd.Flags().IntP("port", "p", 7088, "port for running application")
	rootCmd.Flags().IntP("maxClients", "", 10000, "max connections can be handled")
	rootCmd.Flags().IntP("maxTimeout", "", 120, "max timeout for clients(in seconds)")

	// Bind the flags to config
	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))
	viper.BindPFlag("host", rootCmd.Flags().Lookup("host"))
	viper.BindPFlag("maxClients", rootCmd.Flags().Lookup("maxClients"))
	viper.BindPFlag("maxTimeout", rootCmd.Flags().Lookup("maxTimeout"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("logging", rootCmd.PersistentFlags().Lookup("logging"))
	viper.BindPFlag("logfile", rootCmd.PersistentFlags().Lookup("logfile"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
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
		klogs.PrintErrorAndExit(err, 2)
	}

	klogs.InitLoggers(appConfig)
	srv.Start(appConfig)
}
