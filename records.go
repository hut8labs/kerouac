package main

import (
	"code.google.com/p/go-sqlite/go1/sqlite3"
	"io"
	"os"
	"path/filepath"
	"time"
)

func CreateBuildRecord(buildId BuildId) error {
	conn, err := getConn(buildId.RootDir)
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

func FindMatchingBuildIds(rootDir string, project string, tag string, datetime string) ([]BuildId, error) {
	query := "SELECT project, tag, started_at FROM builds WHERE 1 = 1"

	args := make([]interface{}, 0, 0)

	if project != "" {
		query = query + " AND project = ?"
		args = append(args, project)
	}

	if tag != "" {
		query = query + " AND tag = ?"
		args = append(args, tag)
	}

	if datetime != "" {
		query = query + " AND started_at = ?"
		args = append(args, datetime)
	}

	query = query + " ORDER BY started_at DESC;"

	conn, err := getConn(rootDir)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	conn.Query(query, args...)

	buildIds := make([]BuildId, 0, 0)

	stmt, err := conn.Query(query, args...)

	if err == io.EOF {
		return buildIds, nil
	} else if err != nil {
		return nil, err
	}

	for {
		buildId, err := scanBuildId(rootDir, stmt)
		if err != nil {
			return nil, err
		}
		buildIds = append(buildIds, buildId)
		if err = stmt.Next(); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}

	return buildIds, nil
}

func FindLatestBuildId(rootDir string, project string, tag string, datetime string) (*BuildId, error) {
	buildIds, err := FindMatchingBuildIds(rootDir, project, tag, datetime)
	if err != nil {
		return nil, err
	}
	if len(buildIds) == 0 {
		return nil, nil
	}
	return &buildIds[0], nil
}

func scanBuildId(rootDir string, stmt *sqlite3.Stmt) (BuildId, error) {
	var rowProject, rowTag, rowDatetime string
	err := stmt.Scan(&rowProject, &rowTag, &rowDatetime)
	if err != nil {
		return BuildId{}, err
	}

	dateTime, err := time.Parse(DateFormat, rowDatetime)
	if err != nil {
		return BuildId{}, err
	}
	return BuildIdAt(rootDir, rowProject, rowTag, dateTime), nil
}

func updateBuildStatus(buildId BuildId, status string) error {
	conn, err := getConn(buildId.RootDir)
	if err != nil {
		return err
	}
	defer conn.Close()

	return conn.Exec("UPDATE builds SET status = ?, finished_at = ? WHERE project = ? AND tag = ? AND started_at = ?", status, time.Now().UTC().Format(DateFormat), buildId.Project, buildId.Tag, buildId.DateTime.Format(DateFormat))

}

func getConn(rootDir string) (*sqlite3.Conn, error) {
	buildDbPath := FmtBuildDbPath(rootDir)
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
