package version

import (
	"fmt"
	"runtime"
)

var (
	GitHash   = "unknown"
	GitBranch = "unknown"
)

func Version() string {
	return fmt.Sprintf("Go Version: %s\nGit Branch: %s\nGitHash: %s", runtime.Version(), GitBranch, GitHash)
}

func PrintVersion() {
	fmt.Println(Version())
}
