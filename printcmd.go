package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func DoPrintCommand() {
	flag.Usage = func() {
		fmt.Printf("Usage: kerouac print [options] <builddir|stdoutpath|stderrpath|kerouaclogpath|tarballpath> <kerouacRootDir> <project> <tag> [datetime]\n\n")
		fmt.Printf("Prints to stdout the build directory, stdout log path, etc. of the specified build.\n\n")
		fmt.Printf("If datetime is not specified, uses the latest build for the tag.\n")
	}

	flag.Parse()

	if len(flag.Args()) < 4 || len(flag.Args()) > 5 {
		flag.Usage()
		os.Exit(1)
	}

	path := flag.Arg(0)
	kerouacRoot := flag.Arg(1)
	project := flag.Arg(2)
	tag := flag.Arg(3)
	var datetime string
	if len(flag.Args()) == 5 {
		datetime = flag.Arg(4)
	}

	buildId, err := FindLatestBuildId(kerouacRoot, project, tag, datetime)

	if err != nil {
		log.Fatal(err)
	}

	if buildId == nil {
		os.Exit(1)
	}

	switch path {
	case "builddir":
		fmt.Print(FmtBuildDir(*buildId))
	case "stdoutpath":
		fmt.Print(FmtStdoutLogPath(*buildId))
	case "stderrpath":
		fmt.Print(FmtStderrLogPath(*buildId))
	case "kerouaclogpath":
		fmt.Print(FmtKerouacLogPath(*buildId))
	case "tarballpath":
		fmt.Print(FmtTarballPath(*buildId))
	default:
		log.Printf("Did not recognize path to print: %s\n\n", path)
		flag.Usage()
		os.Exit(1)
	}
}
