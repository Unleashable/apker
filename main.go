// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package main

import (
	"log"
	"os"

	"github.com/unleashable/apker/cmd"
	"github.com/urfave/cli/v2"
)

var (
	app     *cli.App
	version string
)

func init() {

	app = &cli.App{
		Name:                 "apker",
		Usage:                "deploy custom images in seconds.",
		Version:              version,
		Flags:                cmd.Flags,
		Commands:             cmd.Commands,
		EnableBashCompletion: true,
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Mohamed Elbahja",
				Email: "bm9qdW5r@gmail.com",
			},
		},
	}

}

func main() {

	if e := app.Run(os.Args); e != nil {
		log.Fatal(e)
	}
}
