package main

import (
	"code.google.com/p/go-sqlite/go1/sqlite3"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type BuildStatus string

const (
	DateFormat = "2006-01-02 15:04:05"

	FAILED    BuildStatus = "FAILED"
	SUCCEEDED             = "SUCCEEDED"
	RUNNING               = "RUNNING"
)

// RecordedBuild adds the end time of a build and its result to a BuildId.
type RecordedBuild struct {
	*BuildId
	EndTime time.Time
	Status  BuildStatus
}

func (r RecordedBuild) Duration() time.Duration {
	if r.EndTime.IsZero() {
		return time.Since(r.DateTime)
	}
	return r.EndTime.Sub(r.DateTime)
}

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

func MarkBuildFailed(buildId BuildId) error {
	return updateBuildStatus(buildId, FAILED)
}

func MarkBuildSucceeded(buildId BuildId) error {
	return updateBuildStatus(buildId, SUCCEEDED)
}

func FindMatchingBuilds(rootDir string, project string, tag string, datetime string) ([]RecordedBuild, error) {
	query := "SELECT project, tag, started_at, finished_at, status FROM builds WHERE 1 = 1"

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

	recordedBuilds := make([]RecordedBuild, 0, 0)

	stmt, err := conn.Query(query, args...)

	if err == io.EOF {
		return recordedBuilds, nil
	} else if err != nil {
		return nil, err
	}

	for {
		recordedBuild, err := scanBuild(rootDir, stmt)
		if err != nil {
			return nil, err
		}
		recordedBuilds = append(recordedBuilds, recordedBuild)
		if err = stmt.Next(); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}

	return recordedBuilds, nil
}

func FindLatestBuild(rootDir string, project string, tag string, datetime string) (*RecordedBuild, error) {
	recordedBuilds, err := FindMatchingBuilds(rootDir, project, tag, datetime)
	if err != nil {
		return nil, err
	}
	if len(recordedBuilds) == 0 {
		return nil, nil
	}
	return &recordedBuilds[0], nil
}

func FindBuildsGreaterThanN(rootDir string, project string, n int) ([]RecordedBuild, error) {
	if n < 0 {
		return nil, fmt.Errorf("Cannot find builds greater than %d", n)
	}

	recordedBuilds, err := FindMatchingBuilds(rootDir, project, "", "")
	if err != nil {
		return recordedBuilds, err
	}

	if n > len(recordedBuilds) {
		n = len(recordedBuilds)
	}

	return recordedBuilds[n:], nil
}

func scanBuild(rootDir string, stmt *sqlite3.Stmt) (RecordedBuild, error) {
	var rowProject, rowTag, rowDatetime, rowEndTime, rowStatus string
	err := stmt.Scan(&rowProject, &rowTag, &rowDatetime, &rowEndTime, &rowStatus)
	if err != nil {
		return RecordedBuild{}, err
	}

	dateTime, err := time.Parse(DateFormat, rowDatetime)
	if err != nil {
		return RecordedBuild{}, err
	}

	var endTime time.Time
	if rowEndTime != "" {
		if endTime, err = time.Parse(DateFormat, rowEndTime); err != nil {
			return RecordedBuild{}, err
		}
	}

	buildId := BuildIdAt(rootDir, rowProject, rowTag, dateTime)
	return RecordedBuild{&buildId, endTime, BuildStatus(rowStatus)}, nil
}

func updateBuildStatus(buildId BuildId, status BuildStatus) error {
	conn, err := getConn(buildId.RootDir)
	if err != nil {
		return err
	}
	defer conn.Close()

	return conn.Exec("UPDATE builds SET status = ?, finished_at = ? WHERE project = ? AND tag = ? AND started_at = ?", string(status), time.Now().UTC().Format(DateFormat), buildId.Project, buildId.Tag, buildId.DateTime.Format(DateFormat))

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
	return conn.Exec("INSERT INTO builds (project, tag, started_at, status) VALUES (?, ?, ?, ?)", buildId.Project, buildId.Tag, buildId.DateTime.Format(DateFormat), string(RUNNING))
}
