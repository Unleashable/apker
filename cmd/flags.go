// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package cmd

import (
	"github.com/urfave/cli/v2"
	"os"
)

// Global cli flags
var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:    "pub",
		Aliases: []string{"p"},
		Value:   os.ExpandEnv("$HOME/.ssh/id_rsa.pub"),
		Usage:   "Set ssh public key `path`",
	},
	&cli.StringFlag{
		Name:    "key",
		Aliases: []string{"i"},
		Value:   os.ExpandEnv("$HOME/.ssh/id_rsa"),
		Usage:   "Set ssh private key `path`",
	},
}
