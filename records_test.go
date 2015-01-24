package main

import (
	"testing"
	"time"
)

func knownRecordedBuild() RecordedBuild {
	buildId := knownBuildId()
	return RecordedBuild{&buildId, buildId.DateTime.Add(1 * time.Minute), SUCCEEDED}
}

func TestRecordedBuildDuration(t *testing.T) {
	recordedBuild := knownRecordedBuild()
	expectedDuration := 1 * time.Minute
	if duration := recordedBuild.Duration(); duration != expectedDuration {
		t.Errorf("Duration() returned %s not %s", duration, expectedDuration)
	}
}
