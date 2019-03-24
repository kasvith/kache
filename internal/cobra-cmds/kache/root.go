/*
 * MIT License
 *
 * Copyright (c) 2019 Kasun Vithanage
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
 *
 */

package kache

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	cobracmds "github.com/kasvith/kache/internal/cobra-cmds"
	"github.com/kasvith/kache/internal/config"
	"github.com/kasvith/kache/internal/klogs"
	"github.com/kasvith/kache/internal/srv"
)

var verbose bool
var cfgFile string

// RootCmd is the root command for cobra
var RootCmd = &cobra.Command{
	Use:   "kache",
	Short: "kache is a simple distributed in memory database",
	Long:  `A fast and a flexible in memory redis-compatible database`,
	Run:   runApp,
}

func init() {
	cobra.OnInitialize(initConfig)

	// Flags
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	RootCmd.PersistentFlags().BoolP("debug", "d", false, "output debug information")
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "configuration file")
	RootCmd.PersistentFlags().Bool("logging", true, "set application logs")
	RootCmd.PersistentFlags().String("logfile", "", "application log file")
	RootCmd.PersistentFlags().String("logtype", "default",
		`kache can output logs in different formats like json or logfmt. The default one is custom to kache.`)

	RootCmd.Flags().StringP("host", "", "127.0.0.1", "host for running application")
	RootCmd.Flags().IntP("port", "p", 7088, "port for running application")
	RootCmd.Flags().IntP("maxClients", "", 10000, "max connections can be handled")
	RootCmd.Flags().IntP("maxTimeout", "", 120, "max timeout for clients(in seconds)")

	// Bind the flags to config
	viper.BindPFlag("port", RootCmd.Flags().Lookup("port"))
	viper.BindPFlag("host", RootCmd.Flags().Lookup("host"))
	viper.BindPFlag("maxClients", RootCmd.Flags().Lookup("maxClients"))
	viper.BindPFlag("maxTimeout", RootCmd.Flags().Lookup("maxTimeout"))
	viper.BindPFlag("verbose", RootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("logging", RootCmd.PersistentFlags().Lookup("logging"))
	viper.BindPFlag("logfile", RootCmd.PersistentFlags().Lookup("logfile"))
	viper.BindPFlag("logtype", RootCmd.PersistentFlags().Lookup("logtype"))
	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))
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
		viper.SetConfigName("kache.default")
		viper.AddConfigPath("../config")
	}

	if err := viper.ReadInConfig(); err != nil {
		if !loadingDefault {
			fmt.Fprintf(os.Stderr, "Error reading config from %s : %s\n", viper.ConfigFileUsed(), err)
			os.Exit(1)
		}

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

// Execute will start kache
func Execute() {
	// Commands
	RootCmd.AddCommand(cobracmds.VersionCmd)

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runApp(cmd *cobra.Command, args []string) {
	var appConfig config.AppConfig
	if err := viper.Unmarshal(&appConfig); err != nil {
		klogs.PrintErrorAndExit(err, 2)
	}

	appConfig.MaxMultiBulkLength = config.DefaultMaxMultiBulkLength
	config.AppConf = appConfig

	fmt.Println(getASCIIBanner())
	fmt.Printf("Started at: %s\n", time.Now().Format(time.RFC850))
	fmt.Printf("PID: %d\n", os.Getpid())
	fmt.Printf("Port: %d\n", appConfig.Port)
	fmt.Println()

	klogs.InitLoggers(appConfig)
	srv.Start(appConfig)
}

func getASCIIBanner() string {
	val := `

 ██ ▄█▀▄▄▄       ▄████▄   ██░ ██ ▓█████ 
 ██▄█▒▒████▄    ▒██▀ ▀█  ▓██░ ██▒▓█   ▀ 
▓███▄░▒██  ▀█▄  ▒▓█    ▄ ▒██▀▀██░▒███   
▓██ █▄░██▄▄▄▄██ ▒▓▓▄ ▄██▒░▓█ ░██ ▒▓█  ▄ 
▒██▒ █▄▓█   ▓██▒▒ ▓███▀ ░░▓█▒░██▓░▒████▒
▒ ▒▒ ▓▒▒▒   ▓▒█░░ ░▒ ▒  ░ ▒ ░░▒░▒░░ ▒░ ░
░ ░▒ ▒░ ▒   ▒▒ ░  ░  ▒    ▒ ░▒░ ░ ░ ░  ░
░ ░░ ░  ░   ▒   ░         ░  ░░ ░   ░   
░  ░        ░  ░░ ░       ░  ░  ░   ░  ░
                ░
`
	return val
}
