package errh

import (
	"fmt"
	"os"
)

func LogErrorAndExit(err error, exit int) {
	if os.Getenv("ENV") == "DEBUG" {
		panic(err)
	}

	fmt.Fprintln(os.Stderr, err)
	os.Exit(exit)
}
