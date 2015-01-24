package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"
)

const REPORT_MAX_BUILDS = 100 // Eventually, this will be configurable.

type templateFields struct {
	Builds  []RecordedBuild
	CSSPath string
}

func RenderHTMLReport(reportPath string, builds []RecordedBuild) (err error) {
	file, err := os.Create(reportPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Templates can panic(), so set up a recover just in case.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error rendering HTML report: %s", r)
		}
	}()

	numBuilds := len(builds)
	if numBuilds > REPORT_MAX_BUILDS {
		numBuilds = REPORT_MAX_BUILDS
	}
	fields := &templateFields{Builds: builds[0:numBuilds]}

	tryCSSPath := filepath.Join(filepath.Dir(reportPath), "builds.css")
	if stat, err := os.Stat(tryCSSPath); err == nil && !stat.IsDir() {
		fields.CSSPath = tryCSSPath
	}

	funcMap := map[string]interface{}{
		"relative": func(path string) (string, error) {
			return filepath.Rel(filepath.Dir(reportPath), path)
		},
		"base": func(path string) string {
			return filepath.Base(path)
		},
		"friendlyDate": func(timestamp time.Time) string {
			return timestamp.Format(time.RFC1123)
		},
	}
	htmlTemplate := template.Must(template.New("HTMLReport").Funcs(funcMap).Parse(HTMLTemplate))
	return htmlTemplate.Execute(file, fields)
}

var HTMLTemplate = `<!doctype html>
<html>
<head>
  <title>Kerouac: Build Report</title>
  <style>
    table { border-collapse: collapse; }
	table, th, td { border: 1px solid black; }
    th, td { padding: 1em; text-align: center; }
  </style>
  {{ if .CSSPath }}<link rel="stylesheet" type="text/css" href="{{ .CSSPath | relative }}" />{{ end }}
</head>
<body>
<h1>Kerouac: Build Report</h1>
<table>
<thead>
<tr>
<th>Project</th>
<th>Tag</th>
<th>Start</th>
<th>End</th>
<th>Duration</th>
<th>Status</th>
<th>Logs</th>
<th>Tarball</th>
</tr>
</thead>
<tbody>
{{ range .Builds }}
<tr class="build status-{{ .Status }}">
  <td class="project">{{ .Project }}</td>
  <td class="tag">{{ .Tag }}</td>
  <td class="start">{{ .DateTime | friendlyDate }}</td>
  <td class="end">{{ if .EndTime }}{{ .EndTime | friendlyDate }}{{ end }}</td>
  <td class="duration">{{ .Duration }}</td>
  <td class="status">{{ .Status }}</td>
  <td class="logs">
	<a href="{{ .FmtStdoutLogPath | relative }}">{{ .FmtStdoutLogPath | base }}</a>
	<a href="{{ .FmtStderrLogPath | relative }}">{{ .FmtStderrLogPath | base }}</a>
	<a href="{{ .FmtKerouacLogPath | relative }}">{{ .FmtKerouacLogPath | base }}</a>
  </td>
  <td class="tarball"><a href="{{ .FmtTarballPath | relative }}">{{ .FmtTarballPath | base }}</a></td>
</tr>
{{ end }}
</tbody>
</body>
</html>`
