// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package actions

import (
	"fmt"
	"os"
	"strings"

	"github.com/melbahja/goph"
	. "github.com/unleashable/apker/cmd/utils"
	"github.com/unleashable/apker/internal"
	"github.com/unleashable/apker/internal/utils"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh"
)

var RunFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     "addr",
		Aliases:  []string{"ip"},
		Usage:    "Set machine `ip` address.",
		Required: true,
	},
	&cli.StringFlag{
		Name:  "user",
		Value: "root",
		Usage: "Set ssh `user` name.",
	},
	&cli.BoolFlag{
		Name:  "passphrase",
		Usage: "Ask for private key passphrase for protected keys.",
	},
	&cli.BoolFlag{
		Name:  "password",
		Usage: "Ask for ssh password instead of using private keys.",
	},
	&cli.BoolFlag{
		Name:  "insecure",
		Usage: "Do not check knownhosts.",
	},
	&cli.StringSliceFlag{
		Name:    "env",
		Usage:   "Set action env variables.",
		Aliases: []string{"e"},
	},
}

func Run(c *cli.Context) (e error) {

	var (
		cmd      string = fmt.Sprintf("/usr/share/apker/bin/%s", c.Args().First())
		client   *goph.Client
		output   []byte
		callback ssh.HostKeyCallback
		project  internal.Project = internal.Project{
			Temp: utils.Temp(),
		}
	)

	//Housekeeping.
	defer os.RemoveAll(project.Temp)

	// Set auth method.
	if e = SetAuthMethod(&project, c); e != nil {
		return
	}

	// Set host key callback.
	if c.Bool("insecure") {

		callback = ssh.InsecureIgnoreHostKey()

	} else if callback, e = goph.DefaultKnownHosts(); e != nil {

		return
	}

	// Get ssh client.
	if client, e = goph.NewConn(c.String("user"), c.String("addr"), project.SSHAuth, callback); e != nil {
		return
	}

	// Get apker config file from remote machine.
	if e = client.Download("/usr/share/apker/apker.yaml", project.Temp+"/apker.yaml"); e != nil {
		return
	}

	// Load project config from apker.yaml
	if project.Config, e = internal.LoadConfig(project.Temp, []string{}); e != nil {
		return
	}

	// Override if action on apker file.
	if action := project.Config.GetAction(c.Args().First()); action != "" {
		cmd = action
	}

	// Run the action.
	if output, e = client.Run(fmt.Sprintf(`env %s bash -c '%s'`, env(c.StringSlice("env")), cmd)); e == nil {
		fmt.Println(string(output))
	}

	return e
}

func env(s []string) string {
	return strings.Join(append(s, "APKER_ACTION=yes"), " ")
}
