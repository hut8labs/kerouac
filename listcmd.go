package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func DoListCommand() {
	flag.Usage = func() {
		fmt.Printf("Usage: kerouac list [options] <kerouacRootDir> [project] [tag] [datetime]\n\n")
		fmt.Printf("Prints to stdout the list of build directories matching the supplied criteria.\n\n")
		fmt.Printf("Example: 'kerouac list' would list all builds.\n\n")
		fmt.Printf("Example: 'kerouac list myproj' would list all builds for myproj.\n")
	}

	flag.Parse()

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	kerouacRoot := flag.Arg(0)
	var project, tag, datetime string

	if len(flag.Args()) > 1 {
		project = flag.Arg(1)
	}

	if len(flag.Args()) > 2 {
		tag = flag.Arg(2)
	}

	if len(flag.Args()) > 3 {
		datetime = flag.Arg(3)
	}

	buildIds, err := FindMatchingBuildIds(kerouacRoot, project, tag, datetime)

	if err != nil {
		log.Fatalf("Error finding builds: %s", err)
	}

	for _, buildId := range buildIds {
		fmt.Printf("%s\n", FmtBuildDir(buildId))
	}

}
