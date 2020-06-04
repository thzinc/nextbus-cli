package main

import (
	"github.com/spf13/cobra/doc"
	"github.com/thzinc/nextbus-cli/internal/cmd/nextbus"
)

func main() {
	err := doc.GenMarkdownTree(nextbus.RootCmd, "docs/nextbus")
	if err != nil {
		panic(err)
	}
}
