package main

import (
	"github.com/kotaicode/xrd2crd/cmd/xrd2crd"
)

// Version information (populated by GoReleaser)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	xrd2crd.Main(version, commit, date)
}
