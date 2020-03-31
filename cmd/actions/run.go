// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package actions

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var RunFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "addr",
		Aliases: []string{"ip"},
		Usage:   "Set machine `ip` address.",
	},
}

func Run(c *cli.Context) error {

	fmt.Println("TODO: run actions")
	return nil
}
