package main

import (
	"context"
	"fmt"
	"runtime/debug"
	"strconv"
	"time"
)

const dateTimeWithZone = "2006-01-02 15:04:05 MST"

// VersionCmd shows version information about the running machina executable.
type VersionCmd struct{}

// Run executes the version command.
func (cmd VersionCmd) Run(ctx context.Context) error {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("machina build information is not available")
		return nil
	}

	// Look for build settings that are of interest.
	var (
		commitTime     time.Time
		commitRevision string
		commitModified bool
	)
	for _, setting := range buildInfo.Settings {
		if setting.Key == "vcs.time" && setting.Value != "" {
			commitTime, _ = time.Parse(time.RFC3339, setting.Value)
		}
		if setting.Key == "vcs.revision" && setting.Value != "" {
			commitRevision = setting.Value
		}
		if setting.Key == "vcs.modified" && setting.Value != "" {
			commitModified, _ = strconv.ParseBool(setting.Value)
		}
	}

	// Print the main module version.
	if version := buildInfo.Main.Version; version != "" {
		fmt.Printf("%s\n", version)
	}

	// Print the commit revision.
	if commitRevision != "" {
		if commitModified {
			fmt.Printf("  machina commit revision: %s (modified)\n", commitRevision)
		} else {
			fmt.Printf("  machina commit revision: %s\n", commitRevision)
		}
	}

	// Print the commit date.
	if !commitTime.IsZero() {
		fmt.Printf("  machina commit date: %s\n", commitTime.Local().Format(dateTimeWithZone))
	}

	// Print the go version.
	if version := buildInfo.GoVersion; version != "" {
		fmt.Printf("  go version: %s\n", version)
	}

	return nil
}
