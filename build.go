package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

// Represents the full identification for a build, including the kerouac root,
// the project, the build tag, and the datetime.
type BuildId struct {
	RootDir  string
	Project  string
	Tag      string
	DateTime time.Time
}

func BuildIdAtNow(rootDir string, project string, tag string) BuildId {
	dateTime := time.Now().UTC()
	return BuildId{RootDir: rootDir, Project: project, Tag: tag, DateTime: dateTime}
}

// Contains paths to files containing stdout and stderr from the build process.
type BuildOutput struct {
	StdoutPath string
	StderrPath string
}

// Run the supplied build script, after changing directory to buildDir.
//
// The script's stdout and stderr will be captured and written to the
// appropriate directory under kerouacResultsRootDir (see dirs.go for more).
func RunBuildScript(buildDir string, buildScript string, buildScriptArgs []string, timeoutInSecs int, buildId BuildId) (*BuildOutput, error) {
	cmd := exec.Command(buildScript, buildScriptArgs...)
	cmd.Dir = buildDir

	stdoutPath := FmtStdoutLogPath(buildId)
	stderrPath := FmtStderrLogPath(buildId)

	buildOutput := &BuildOutput{StdoutPath: stdoutPath, StderrPath: stderrPath}

	stdoutFile, err := os.Create(stdoutPath)
	if err != nil {
		return nil, err
	}
	defer stdoutFile.Close()

	stderrFile, err := os.Create(stderrPath)
	if err != nil {
		return nil, err
	}
	defer stderrFile.Close()

	cmd.Stdout = stdoutFile
	cmd.Stderr = stderrFile

	cmdDone := make(chan error)

	go runCmd(cmd, cmdDone)

	err = waitForCmd(cmd, timeoutInSecs, cmdDone)

	return buildOutput, err
}

func runCmd(cmd *exec.Cmd, cmdDone chan<- error) {
	cmdDone <- cmd.Run()
}

func waitForCmd(cmd *exec.Cmd, timeoutInSecs int, cmdDone <-chan error) error {
	var err error

	select {
	case result := <-cmdDone:
		err = result
	case <-time.After(time.Second * time.Duration(timeoutInSecs)):
		err = fmt.Errorf("Execution of build timed out after %d seconds", timeoutInSecs)
		log.Printf("Attempting to kill long-running build ...")
		if perr := cmd.Process.Kill(); perr != nil {
			// TODO: try -9'ing with signal
			log.Printf("Could not kill process, aborting in dirty state: %s", perr)
			err = perr
		} else {
			log.Printf("Long-running build killed.")
		}
		<-cmdDone
	}
	return err
}
