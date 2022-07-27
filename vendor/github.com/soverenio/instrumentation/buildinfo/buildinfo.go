package buildinfo

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
	"time"
)

var (
	Version   = "unset" // Version is release semantic version.
	GitCommit = "unset" // GitCommit is short git commit hash.
	TimeBuilt = "0"     // TimeBuilt is build timestamp.

	TimeStarted time.Time // TimeStarted is a Time when service started.
)

func init() {
	TimeStarted = time.Now()
}

type Info struct {
	Version     string
	GitCommit   string
	GoVersion   string
	GoCompiler  string
	TimeBuilt   time.Time
	OSArch      string
	TimeStarted time.Time
	Uptime      time.Duration
}

// NewInfo creates a struct with information about service build and runtime.
func NewInfo() Info {
	tb, err := strconv.ParseInt(TimeBuilt, 10, 64)
	if err != nil {
		tb = 0
	}
	return Info{
		Version:     Version,
		GitCommit:   GitCommit,
		GoVersion:   runtime.Version(),
		GoCompiler:  runtime.Compiler,
		TimeBuilt:   time.Unix(tb, 0),
		OSArch:      fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		TimeStarted: TimeStarted,
		Uptime:      time.Since(TimeStarted),
	}
}

// String returns multi line text representation of Info.
func (i Info) String() string {
	return fmt.Sprintf(
		` Version     : %s
 Git commit  : %s
 Go version  : %s
 Go compiler : %s
 Built       : %s
 OS/Arch     : %s
 Started     : %s
 Uptime      : %s`,
		i.Version,
		i.GitCommit,
		i.GoVersion,
		i.GoCompiler,
		i.TimeBuilt.Format(time.RFC1123),
		i.OSArch,
		i.TimeStarted.Format(time.RFC1123),
		i.Uptime,
	)
}

// MarshalJSON converts Info to pretty printed JSON.
func (i Info) MarshalJSON() ([]byte, error) {
	strinfo := struct {
		Version, GitCommit, GoVersion, GoCompiler, TimeBuilt, OSArch, TimeStarted, Uptime string
	}{
		i.Version, i.GitCommit, i.GoVersion, i.GoCompiler, i.TimeBuilt.Format(time.RFC3339),
		i.OSArch, i.TimeStarted.Format(time.RFC3339), i.Uptime.Round(time.Second).String(),
	}
	return json.Marshal(strinfo)
}
