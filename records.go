package main

import (
	"code.google.com/p/go-sqlite/go1/sqlite3"
	"os"
	"path/filepath"
	"time"
)

func CreateBuildRecord(buildId BuildId) error {
	conn, err := getConn(buildId)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err = createTablesAndIndexes(conn); err != nil {
		return err
	}

	if err = insertBuildRecord(conn, buildId); err != nil {
		return err
	}

	return nil
}

const DateFormat = "2006-01-02 15:04:05"

func MarkBuildFailed(buildId BuildId) error {
	return updateBuildStatus(buildId, "FAILED")
}

func MarkBuildSucceeded(buildId BuildId) error {
	return updateBuildStatus(buildId, "SUCCEDED")
}

func updateBuildStatus(buildId BuildId, status string) error {
	conn, err := getConn(buildId)
	if err != nil {
		return err
	}
	defer conn.Close()

	return conn.Exec("UPDATE builds SET status = ?, finished_at = ? WHERE project = ? AND tag = ? AND started_at = ?", status, time.Now().UTC().Format(DateFormat), buildId.Project, buildId.Tag, buildId.DateTime.Format(DateFormat))

}

func getConn(buildId BuildId) (*sqlite3.Conn, error) {
	buildDbPath := FmtBuildDbPath(buildId.RootDir)
	os.MkdirAll(filepath.Dir(buildDbPath), 0700)
	return sqlite3.Open(buildDbPath)
}

const createBuildsTable = "CREATE TABLE IF NOT EXISTS builds (id INTEGER PRIMARYKEY ASC, project TEXT NOT NULL, tag TEXT NOT NULL, started_at TEXT NOT NULL, finished_at TEXT, status TEXT)"

const createBuildsUniqueIdx = "CREATE UNIQUE INDEX IF NOT EXISTS builds_idx ON builds (project, tag, started_at)"

func createTablesAndIndexes(conn *sqlite3.Conn) error {
	stmts := []string{createBuildsTable, createBuildsUniqueIdx}

	for _, stmt := range stmts {
		if err := conn.Exec(stmt); err != nil {
			return err
		}
	}

	return nil
}

// This acts as the locking mechanism to make sure we don't have two builds in
// the identical folder, as well as record keeping.
func insertBuildRecord(conn *sqlite3.Conn, buildId BuildId) error {
	return conn.Exec("INSERT INTO builds (project, tag, started_at, status) VALUES (?, ?, ?, ?)", buildId.Project, buildId.Tag, buildId.DateTime.Format(DateFormat), "RUNNING")
}
