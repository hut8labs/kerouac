package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var dryRun = flag.Bool("dry-run", false, "Print the commands that would be run.")

var noRemoveSrcDir = flag.Bool("no-remove-src", false, "Do not remove the source dir after building.")

// We expect 5 arguments on the command line
const NumArgs = 5

func main() {
	flag.Parse()

	if len(flag.Args()) != NumArgs {
		// TODO: need to make the usage describe the args
		flag.Usage()
		os.Exit(1)
	}

	srcDir := flag.Arg(0)
	configFile := flag.Arg(1)
	rootDir := flag.Arg(2)
	project := flag.Arg(3)
	tag := flag.Arg(4)

	buildId := BuildIdAtNow(rootDir, project, tag)

	log.SetOutput(os.Stderr)
	if *dryRun {
		log.Printf("Dry run, will print actions but not take them.")
	}

	createBuildRecord(buildId)

	logFile := configureLogging(buildId)
	defer logFile.Close()

	logStart(buildId)

	config, err := ParseConfigFile(configFile)
	if err != nil {
		logAndDie(fmt.Sprintf("Error parsing config file: %s", err), buildId)
	}

	runBuild(srcDir, config, buildId)
	createTarball(srcDir, buildId)
	removeSrcDir(srcDir)

	// TODO clean up old builds unless told not to

	// TODO update the html
}

func logAndDie(msg string, buildId BuildId) {
	if err := MarkBuildFailed(buildId); err != nil {
		log.Printf("Could not mark build failed in db: %s", err)
	}
	log.Fatalf(msg)
}

func createBuildRecord(buildId BuildId) {
	log.Printf("Creating db record for build.")

	if !*dryRun {
		if err := CreateBuildRecord(buildId); err != nil {
			log.Fatalf("Could not create build record: %s", err)
		}
	}
}

func removeSrcDir(srcDir string) {
	if *noRemoveSrcDir {
		log.Printf("Not removing source dir due to --no-remove-src.")
	} else {
		log.Printf("Removing source dir %s", srcDir)
		if !*dryRun {
			os.RemoveAll(srcDir)
		}
	}
}

func createTarball(srcDir string, buildId BuildId) {
	log.Printf("Tarballing %s into %s", srcDir, FmtTarballPath(buildId))

	if !*dryRun {
		if err := CreateTarball(srcDir, buildId); err != nil {
			logAndDie(fmt.Sprintf("Error creating tarball: %s", err), buildId)
		}
	}
}

func runBuild(srcDir string, config *Config, buildId BuildId) {
	log.Printf("Running build in dir %s with script %s and args %s", srcDir, config.BuildScript, config.BuildScriptArgs)

	if !*dryRun {
		buildOutput, err := RunBuildScript(srcDir, config.BuildScript, config.BuildScriptArgs, buildId)

		if err != nil {
			log.Printf("Completed build with error: %s", err)
			MarkBuildFailed(buildId)
		} else {
			log.Printf("Completed build successfully.")
			MarkBuildSucceeded(buildId)
		}

		log.Printf("Build script stdout in: %s", buildOutput.StdoutPath)
		log.Printf("Build scrip stderr in: %s", buildOutput.StderrPath)
	}
}

func configureLogging(buildId BuildId) *os.File {
	logsDir := FmtLogsDir(buildId)

	log.Printf("Creating logs dir %s with perms 0700", logsDir)

	if !*dryRun {
		// TODO: reconsider permissions
		os.MkdirAll(logsDir, 0700)
	}

	logPath := FmtKerouacLogPath(buildId)

	log.Printf("Creating log at %s", logPath)

	var logFile *os.File
	var err error

	if !*dryRun {
		logFile, err = os.Create(logPath)
		if err != nil {
			log.Fatalf("Logging ironic error, could not configure logging: %s", err)
		}
	}

	log.Printf("Teeing log output to %s and stderr", logPath)

	if !*dryRun {
		writer := io.MultiWriter(os.Stderr, logFile)
		log.SetOutput(writer)
	}

	return logFile
}

func logStart(buildId BuildId) {
	log.Printf("Starting build of %s with tag %s at %s", buildId.Project, buildId.Tag, buildId.DateTime.Format("2006-01-02 15:04:05 (MST)"))
	log.Printf("Build dir is %s", FmtBuildDir(buildId))
}
