// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package cmd

import (
	"github.com/urfave/cli/v2"
)

// Global cli flags
var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:    "pub",
		Aliases: []string{"p"},
		Usage:   "Set ssh public key (`public_key`)",
	},
	&cli.StringFlag{
		Name:    "key",
		Aliases: []string{"i"},
		Usage:   "Set ssh private key (`identity_file`)",
	},
}
