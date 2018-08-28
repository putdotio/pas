package main

import (
	"flag"
	"fmt"
)

var Version string

func init() {
	if Version == "" {
		Version = "v0.0.0"
	}
}

var version = flag.Bool("version", false, "version")

func main() {
	flag.Parse()
	if *version {
		fmt.Println(Version)
	}
}
