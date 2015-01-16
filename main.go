package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	mode := ""

	if len(os.Args) > 1 {
		mode = os.Args[1]
		if len(os.Args) > 2 {
			os.Args = append(os.Args[:1], os.Args[2:]...)
		} else {
			os.Args = os.Args[0:1]
		}
	} else {
		usage()
	}

	// Subcommands may override this.
	log.SetOutput(os.Stderr)

	switch mode {
	case "build":
		DoBuildCommand()
	case "list":
		DoListCommand()
	case "print":
		DoPrintCommand()
	default:
		usage()
	}
}

func usage() {
	fmt.Printf("Usage: kerouac {build, list, print}\n")
	fmt.Printf("\n")
	fmt.Printf("Use kerouac <subcommand> -h for help.\n")
	os.Exit(1)
}
