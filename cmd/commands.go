// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package cmd

import (
	"github.com/unleashable/apker/cmd/actions"
	"github.com/urfave/cli/v2"
)

var Commands = []*cli.Command{
	{
		Name:    "deploy",
		Aliases: []string{"dep"},
		Usage:   "Deploy from remote git repository.",
		Action:  actions.Deploy,
		Flags:   actions.DeployFlags,
	},
	{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "List machines.",
		Action:  actions.List,
	},
	{
		Name:   "run",
		Usage:  "Run an action on remote machine.",
		Action: actions.Run,
		Flags:  actions.RunFlags,
	},
}
