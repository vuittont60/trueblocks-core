package main

import (
	"os"

	"github.com/TrueBlocks/trueblocks-core/src/go-apps/findSig/cmd"
)

func main() {
	cmd.Execute()
	os.Exit(0)
}
