// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package actions

import (
	"fmt"
	"os"
	"strings"

	"github.com/melbahja/goph"
	"github.com/unleashable/apker/cmd/inputs"
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
	&cli.StringFlag{
		Name:    "key",
		Value:   os.ExpandEnv("$HOME/.ssh/id_rsa"),
		Usage:   "Set private key path.",
		Aliases: []string{"i"},
	},
	&cli.StringFlag{
		Name:  "knownhosts",
		Usage: "knownhosts `file`.",
		Value: os.ExpandEnv("$HOME/.ssh/known_hosts"),
	},
	&cli.BoolFlag{
		Name:  "passphrase",
		Usage: "Ask for private key passphrase for protected keys.",
	},
	&cli.BoolFlag{
		Name:  "agent",
		Usage: "Use ssh agent.",
	},
	&cli.BoolFlag{
		Name:  "password",
		Usage: "Ask for ssh password instead of using private keys.",
	},
	&cli.StringSliceFlag{
		Name:    "env",
		Usage:   "Set action env variables: (`NAME=VALUE`).",
		Aliases: []string{"e"},
	},
}

func Run(c *cli.Context) (e error) {

	var (
		cmd      string = fmt.Sprintf("/usr/share/apker/bin/%s", c.Args().First())
		pass     string
		auth     goph.Auth
		client   *goph.Client
		output   []byte
		callback ssh.HostKeyCallback
	)

	if callback, e = goph.KnownHosts(c.String("knownhosts")); e != nil {
		return
	}

	if c.Bool("password") {

		pass, e = inputs.Password("Enter ssh password", func(pass string) error {
			return nil
		})

		if e != nil {
			return
		}

		auth = goph.Password(pass)

	} else if c.Bool("agent") {

		auth = goph.UseAgent()

	} else {

		if c.Bool("passphrase") {

			pass, e = inputs.Password("Enter private key passphrase", func(pass string) error {
				return nil
			})
		}

		auth = goph.Key(c.String("key"), pass)
	}

	// Get ssh client.
	if client, e = goph.NewConn(c.String("user"), c.String("addr"), auth, callback); e != nil {
		return
	}

	// Run the action.
	output, e = client.Run(fmt.Sprintf(`env %s bash -c '%s'`, env(c.StringSlice("env")), cmd))

	fmt.Println("")
	fmt.Println(string(output))

	return e
}

func env(s []string) string {
	return strings.Join(append(s, "APKER_ACTION=1"), " ")
}
