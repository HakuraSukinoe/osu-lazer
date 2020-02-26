package main

import (
	"github.com/deissh/osu-lazer/server/cmd/risuserver/commands"
	"os"
)

func main() {
	if err := commands.Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
