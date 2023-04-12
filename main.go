package main

import "github.com/souhail-5/zeed/cmd"

var (
	version = "development" // zeed version, automatically updated by GoReleaser ldflags during build.
)

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
