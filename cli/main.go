/*
Copyright (c) 2023 Gsoc2 Inc.
*/
package main

import (
	"os"

	"github.com/Gsoc2/gsoc2-merge/packages/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	cmd.Execute()
}
