package version

import (
	"fmt"
	"runtime"
)

var (
	// These will be set during build
	Version    = "dev"
	BuildTime  = "unknown"
	CommitHash = "unknown"
)

func GetVersion() string {
	return fmt.Sprintf("Solar Boat CLI v%s (commit: %s, built: %s, %s %s/%s)",
		Version,
		CommitHash,
		BuildTime,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)
}
