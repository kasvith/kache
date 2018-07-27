package config

import (
	"os"

	"github.com/jessevdk/go-flags"
)

// Options for the application
type Options struct {
	// Port for running the application
	Port uint `short:"p" long:"port" description:"Port for running application" default:"7869"`

	// Host for running the application
	Host string `long:"host" description:"Host for running application" default:"0.0.0.0"`
}

// ParseConfig parses configuration for the application
func ParseConfig(args []string) (options Options) {
	args, err := flags.ParseArgs(&options, args)

	if err != nil {
		os.Exit(1)
	}

	return
}
