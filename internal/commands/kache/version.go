package kache

import (
	"fmt"
	"github.com/spf13/cobra"
)

const APPVER = "1.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display application version",
	Long:  `Display application version on the screen`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO automated build version to be added
		fmt.Println("Version", APPVER)
	},
}
