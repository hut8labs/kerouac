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
	BuildsDir           = "builds"
	LogsDir             = "logs"
	StderrLogName       = "stderr"
	StdoutLogName       = "stdout"
	KerouacLogName      = "kerouac.log"
	TarballName         = "build.tar.gz"
	BuildDbName         = "builds.db"
	BuildHTMLReportName = "builds.html"
)

func (buildId BuildId) FmtBuildDir() string {
	dateTag := buildId.DateTime.Format("2006_01_02_15_04_05")
	return filepath.Join(buildId.RootDir, BuildsDir, buildId.Project, buildId.Tag, dateTag)
}

func (buildId BuildId) FmtLogsDir() string {
	return filepath.Join(buildId.FmtBuildDir(), LogsDir)
}

func (buildId BuildId) FmtStderrLogPath() string {
	return filepath.Join(buildId.FmtLogsDir(), StderrLogName)
}

func (buildId BuildId) FmtStdoutLogPath() string {
	return filepath.Join(buildId.FmtLogsDir(), StdoutLogName)
}

func (buildId BuildId) FmtKerouacLogPath() string {
	return filepath.Join(buildId.FmtLogsDir(), KerouacLogName)
}

func (buildId BuildId) FmtTarballPath() string {
	return filepath.Join(buildId.FmtBuildDir(), TarballName)
}

func FmtBuildDbPath(rootDir string) string {
	return filepath.Join(rootDir, BuildDbName)
}

func FmtBuildHTMLReportPath(rootDir string) string {
	return filepath.Join(rootDir, BuildHTMLReportName)
}
