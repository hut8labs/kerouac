package main

import (
	"path/filepath"
	"testing"
	"time"
)

// Hardcoded values for our synthetic build.
const (
	KnownRootDir    = "known_test_root"
	KnownProject    = "known_test_project"
	KnownTag        = "known_test_tag"
	KnownDateTimeS  = "2006-01-02 15:04:05"
	KnownDateTimeSU = "2006_01_02_15_04_05"
)

// Do a bunch of the same work we expect layout to do so we can test, but do it
// with filepath.Join just in case this someday works on Windows by pure
// accident.
var (
	KnownDateTime, _ = time.Parse(KnownDateTimeS, KnownDateTimeS)
	KnownBuildDir    = filepath.Join(KnownRootDir, BuildsDir, KnownProject, KnownTag, KnownDateTimeSU)
	KnownLogsDir     = filepath.Join(KnownBuildDir, LogsDir)
	KnownStderrPath  = filepath.Join(KnownLogsDir, StderrLogName)
	KnownStdoutPath  = filepath.Join(KnownLogsDir, StdoutLogName)
	KnownKerouacPath = filepath.Join(KnownLogsDir, KerouacLogName)
	KnownTarballPath = filepath.Join(KnownBuildDir, TarballName)
	KnownBuildDbPath = filepath.Join(KnownRootDir, BuildDbName)
)

// Create a known build id from constants, including the datetime, so we can
// expect deterministic results.
func knownBuildId() BuildId {
	return BuildId{RootDir: KnownRootDir, Project: KnownProject, Tag: KnownTag, DateTime: KnownDateTime}
}

func TestFmtBuildDir(t *testing.T) {
	buildId := knownBuildId()
	buildDir := FmtBuildDir(buildId)
	if buildDir != KnownBuildDir {
		t.Errorf("FmtBuildDir returned %s not %s", buildDir, KnownBuildDir)
	}
}

func TestFmtLogsDir(t *testing.T) {
	buildId := knownBuildId()
	logsDir := FmtLogsDir(buildId)
	if logsDir != KnownLogsDir {
		t.Errorf("FmtLogsDir returned %s not %s", logsDir, KnownLogsDir)
	}
}

func TestFmtStderrLogPath(t *testing.T) {
	buildId := knownBuildId()
	stderrPath := FmtStderrLogPath(buildId)
	if stderrPath != KnownStderrPath {
		t.Errorf("FmtStderrLogPath returned %s not %s", stderrPath, KnownStderrPath)
	}
}

func TestFmtStdoutLogPath(t *testing.T) {
	buildId := knownBuildId()
	stdoutPath := FmtStdoutLogPath(buildId)
	if stdoutPath != KnownStdoutPath {
		t.Errorf("FmtStdoutLogPath returned %s not %s", stdoutPath, KnownStdoutPath)
	}
}

func TestFmtKerouacLogPath(t *testing.T) {
	buildId := knownBuildId()
	stdoutPath := FmtKerouacLogPath(buildId)
	if stdoutPath != KnownKerouacPath {
		t.Errorf("FmtKerouacLogPath returned %s not %s", stdoutPath, KnownKerouacPath)
	}
}

func TestFmtTarballPath(t *testing.T) {
	buildId := knownBuildId()
	tarballPath := FmtTarballPath(buildId)
	if tarballPath != KnownTarballPath {
		t.Errorf("FmtTarballPath returned %s not %s", tarballPath, KnownTarballPath)
	}
}

func TestFmtBuildDbName(t *testing.T) {
	buildId := knownBuildId()
	buildDbPath := FmtBuildDbPath(buildId.RootDir)
	if buildDbPath != KnownBuildDbPath {
		t.Errorf("FmtBuildDbPath returned %s not %s", buildDbPath, KnownBuildDbPath)
	}
}
