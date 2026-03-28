package main

import (
	"os"

	"erwin/harnessruntime/internal/harness"
)

func main() {
	os.Exit(harness.Main(os.Args[1:]))
}
