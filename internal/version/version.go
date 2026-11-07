package version

import (
	"fmt"
	"runtime"
)

var (
	// SemVer The semantic version of this release
	SemVer = "develop"

	// CommitSHA The git commit SHA that this release was built against
	CommitSHA = "develop"

	// BuildDate Date/time of when this release was built
	BuildDate = "develop"
)

// String Returns a string containing SemVer, CommitSHA, BuildDate, and runtime.Version
func String() string {
	return fmt.Sprintf("%s (commit %s, built at %s, runtime %s)",
		SemVer, CommitSHA, BuildDate, runtime.Version())
}
