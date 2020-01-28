package appdetail

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	Version   string
	Commit    string
	BuildTime string
	Title     string
	StartTime time.Time
	Hostname  string
)

func init() {
	// If version, commit, or build time are not set, make that clear.
	if Version == "" {
		Version = "unknown"
	}
	if Commit == "" {
		Commit = "unknown"
	}
	if BuildTime == "" {
		BuildTime = "unknown"
	}
	if Title == "" {
		Title = "appdetail"
	}

	StartTime = time.Now()

	var err error
	Hostname, err = os.Hostname()
	if err != nil {
		logrus.WithError(err).Warn("Failed to set Runtime Hostname")
	}
}
