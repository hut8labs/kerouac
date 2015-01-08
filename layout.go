package main

import (
	"path/filepath"
)

// Code to express the kerouac conventions around filesytem layout.
//
// In short, there is a root directory for kerouac results, under which we
// enforce the following layout (the [FmtXYZ] denote which functions will
// return that path):
//
// - builds.db [FmtBuildDbPath]
// - builds
//   - project_one
//     - buildtag
//       - datetag [FmtBuildDir]
//         build.tar.gz [FmtTarballPath]
//         - logs [FmtLogsDir]
//             stdout [FmtStdoutLogPath]
//             stderr [FmtStderrLogPath]
//             kerouac.log [FmtKerouacLogPath]
// - index.html
// - pages
//

const (
	BuildsDir      = "builds"
	LogsDir        = "logs"
	StderrLogName  = "stderr"
	StdoutLogName  = "stdout"
	KerouacLogName = "kerouac.log"
	TarballName    = "build.tar.gz"
	BuildDbName    = "builds.db"
)

func FmtBuildDir(buildId BuildId) string {
	dateTag := buildId.DateTime.Format("2006_01_02_15_04_05")
	return filepath.Join(buildId.RootDir, BuildsDir, buildId.Project, buildId.Tag, dateTag)
}

func FmtLogsDir(buildId BuildId) string {
	return filepath.Join(FmtBuildDir(buildId), LogsDir)
}

func FmtStderrLogPath(buildId BuildId) string {
	return filepath.Join(FmtLogsDir(buildId), StderrLogName)
}

func FmtStdoutLogPath(buildId BuildId) string {
	return filepath.Join(FmtLogsDir(buildId), StdoutLogName)
}

func FmtKerouacLogPath(buildId BuildId) string {
	return filepath.Join(FmtLogsDir(buildId), KerouacLogName)
}

func FmtTarballPath(buildId BuildId) string {
	return filepath.Join(FmtBuildDir(buildId), TarballName)
}

func FmtBuildDbPath(rootDir string) string {
	return filepath.Join(rootDir, BuildDbName)
}
