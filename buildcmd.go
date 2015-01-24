package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var dryRun = flag.Bool("dry-run", false, "Print the commands that would be run.")

var removeSrcDir = flag.Bool("remove-src", false, "Remove the source dir after building.")

// We expect 5 arguments on the command line
const NumArgs = 5

func DoBuildCommand() {
	flag.Usage = func() {
		fmt.Printf("Usage: kerouac build [options] <srcDir> <configFile> <kerouacRootDir> <project> <tag>\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(flag.Args()) != NumArgs {
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

	buildSuceeded := runBuild(srcDir, config, buildId)
	createTarball(srcDir, buildId)
	maybeRemoveSrcDir(srcDir)

	if buildSuceeded {
		if err = cleanOldBuilds(buildId.RootDir, buildId.Project, config.NumBuildsToKeep); err != nil {
			log.Printf("Warning, error trying to remove old builds: %s", err)
		}

		os.Exit(0)
	} else {
		os.Exit(1)
	}

}

func cleanOldBuilds(rootDir string, project string, buildsToKeep int) error {
	if buildsToKeep < 1 {
		return fmt.Errorf("Refusing to keep < 1 build, not deleting any: %d", buildsToKeep)
	}

	buildIdsToRemove, err := FindBuildIdsGreaterThanN(rootDir, project, buildsToKeep)
	if err != nil {
		return err
	}

	for _, buildId := range buildIdsToRemove {
		buildDir := buildId.FmtBuildDir()
		log.Printf("Removing old build dir %s", buildDir)
		if err = os.RemoveAll(buildDir); err != nil {
			return err
		}
	}

	return nil
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

func maybeRemoveSrcDir(srcDir string) {
	if !*removeSrcDir {
		log.Printf("Not removing source dir.")
	} else {
		log.Printf("Removing source dir due to --remove-src %s", srcDir)
		if !*dryRun {
			os.RemoveAll(srcDir)
		}
	}
}

func createTarball(srcDir string, buildId BuildId) {
	log.Printf("Tarballing %s into %s", srcDir, buildId.FmtTarballPath())

	if !*dryRun {
		if err := CreateTarball(srcDir, buildId); err != nil {
			logAndDie(fmt.Sprintf("Error creating tarball: %s", err), buildId)
		}
	}
}

func runBuild(srcDir string, config *Config, buildId BuildId) bool {
	log.Printf("Running build in dir %s with script %s and args %s", srcDir, config.BuildScript, config.BuildScriptArgs)

	succeeded := false

	if !*dryRun {
		buildOutput, err := RunBuildScript(srcDir, config.BuildScript, config.BuildScriptArgs, config.TimeoutInSecs, buildId)

		if err != nil {
			log.Printf("Completed build with error: %s", err)
			MarkBuildFailed(buildId)
		} else {
			log.Printf("Completed build successfully.")
			MarkBuildSucceeded(buildId)
			succeeded = true
		}

		log.Printf("Build script stdout in: %s", buildOutput.StdoutPath)
		log.Printf("Build script stderr in: %s", buildOutput.StderrPath)
	}

	return succeeded
}

func configureLogging(buildId BuildId) *os.File {
	logsDir := buildId.FmtLogsDir()

	log.Printf("Creating logs dir %s with perms 0700", logsDir)

	if !*dryRun {
		// TODO: reconsider permissions
		os.MkdirAll(logsDir, 0700)
	}

	logPath := buildId.FmtKerouacLogPath()

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
	log.Printf("Build dir is %s", buildId.FmtBuildDir())
}
