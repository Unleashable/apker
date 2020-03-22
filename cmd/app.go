// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package cmd

import "github.com/urfave/cli/v2"

var app *cli.App

func init() {

	app = &cli.App{
		Name:     "apker",
		Usage:    "deploy custom images in seconds.",
		Version:  "v0.0.1",
		Flags:    flags,
		Commands: commands,
	}

}

func Run(args []string) error {
	return app.Run(args)
}
