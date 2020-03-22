// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package cmd

import (
	"github.com/unleashable/apker/cmd/actions"
	"github.com/urfave/cli/v2"
)

var commands = []*cli.Command{
	{
		Name:    "deploy",
		Aliases: []string{"dep"},
		Usage:   "Deploy from current working directory or remote git repo.",
		Action:  actions.Deploy,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Usage: "Set image name.",
			},
			&cli.StringFlag{
				Name:  "size",
				Usage: "Set image size.",
			},
			&cli.StringFlag{
				Name:    "region",
				Aliases: []string{"location"},
				Usage:   "Set image physical location.",
			},
			&cli.StringFlag{
				Name:    "pass",
				Aliases: []string{"passphrase"},
				Usage:   "SSH passphase.",
			},
			&cli.BoolFlag{
				Name:    "ssh",
				Aliases: []string{"enable-ssh"},
				Usage:   "Use ssh key instead of password.",
			},
			&cli.BoolFlag{
				Name:    "wait",
				Aliases: []string{"w"},
				Usage:   "Wait for image to be ready (without timeout).",
			},
			&cli.IntFlag{
				Name:  "id",
				Usage: "Run deploy steps on specifec droplet `id`.",
			},
			&cli.IntFlag{
				Name:  "image",
				Usage: "Create droplet from specifec image `id`.",
			},
		},
	},
	{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "List machines",
		Action:  actions.List,
	},
}
