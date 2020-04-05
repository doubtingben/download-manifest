package main

import (
	"fmt"
	"runtime"
)

var (
	Version   = "v0.0.0.0-dev"
	GoVersion = runtime.Version()
	GitCommit string
	GitBranch string
	BuildDate string
	BuildUser string
)

func printVersion() {
	fmt.Println("Version: ", Version)
	fmt.Println("from commit: ", GitCommit)
	fmt.Println("on: ", BuildDate)
	fmt.Println("by: ", BuildUser)
}
