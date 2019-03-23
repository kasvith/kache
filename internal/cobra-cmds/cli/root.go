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

package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/kasvith/kache/internal/cli"
	cobracmds "github.com/kasvith/kache/internal/cobra-cmds"
)

var host string
var port int

// RootCmd of the CLI
var RootCmd = &cobra.Command{
	Use:   "kache-cli",
	Short: "kache-cli is a client to access kache server",
	Run:   runCli,
}

func init() {
	RootCmd.Flags().StringVarP(&host, "host", "", "127.0.0.1", "host of kache server")
	RootCmd.Flags().IntVarP(&port, "port", "p", 7088, "port of kache server")
}

// Execute CLI
func Execute() {
	RootCmd.AddCommand(cobracmds.VersionCmd)

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runCli(cmd *cobra.Command, args []string) {
	cli.RunCli(host, port)
}
