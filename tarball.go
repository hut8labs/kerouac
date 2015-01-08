package main

import (
	"os/exec"
)

// Create a tarball of srcDir, writing it to the location indicated by
// layout.FmtTarballPath.
//
// Currently this uses the external tar exe, which must be in the path, and
// changes directory to the srcDir before tarring the current directory (".").
//
// Returns nil on success, or an error if something goes wrong.
//
// At some point it might make sense to write this using go's tar and gzip
// utilities, so it could run on systems without tar / be more self contained.
func CreateTarball(srcDir string, buildId BuildId) error {
	tarballPath := FmtTarballPath(buildId)
	cmd := exec.Command("tar", "-c", "-z", "-C", srcDir, "-f", tarballPath, ".")
	return cmd.Run()
}
