package version

import (
	"fmt"
	"runtime"
)

var (
	// Version is the current version of the CLI
	Version = "0.1.2"

	// BuildTime is the time the binary was built
	BuildTime = "unknown"

	// CommitHash is the git commit hash of the build
	CommitHash = "unknown"

	// Runtime is the Go runtime version
	Runtime = fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
)

// FullVersion returns the full version information
func FullVersion() string {
	return fmt.Sprintf(`Solar Boat v%s
  Built:   %s
  Commit:  %s
  Runtime: %s`, Version, BuildTime, CommitHash, Runtime)
}
