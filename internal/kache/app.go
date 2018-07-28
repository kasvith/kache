package kache

import (
	"os"

	"github.com/kasvith/kache/internal/kache/config"
)

// Invoke will start the application
func Invoke() {
	// Prase commandline args
	_ = config.ParseConfig(os.Args)
}
