package errh

import (
	"fmt"
	"log"
	"os"
)

func PrintErrorAndExit(err error, exit int) {
	if os.Getenv("ENV") == "DEBUG" {
		panic(err)
	}

	fmt.Fprintln(os.Stderr, err)
	os.Exit(exit)
}

func LogErrorAndExit(err error, exit int) {

	if os.Getenv("ENV") == "DEBUG" {
		panic(err)
	}

	log.Fatal(err)
	os.Exit(exit)
}

func LogError(params ...string) {
	log.Println(params)
}
