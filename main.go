package main

import (
	"runtime/debug"

	"github.com/kotaicode/xrd2crd/cmd/xrd2crd"
)

func getVersionInfo() (version, commit, date string) {
	version = "dev"
	commit = "none"
	date = "unknown"

	if info, ok := debug.ReadBuildInfo(); ok {
		// Get version from module version
		version = info.Main.Version
		if version == "(devel)" {
			version = "dev"
		}

		// Get commit and date from build settings
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				commit = setting.Value
				if len(commit) > 7 {
					commit = commit[:7]
				}
			case "vcs.time":
				date = setting.Value
			}
		}
	}

	return version, commit, date
}

func main() {
	version, commit, date := getVersionInfo()
	xrd2crd.Main(version, commit, date)
}
