// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package internal

import (
	"bytes"
	"strings"

	"github.com/melbahja/ssh"
)

type OutputHandler func(description string, log *bytes.Buffer) error

type Deployment struct {
	SSH           *ssh.Client
	Project       *Project
	StdoutHandler OutputHandler
	StderrHandler OutputHandler
}

func (d *Deployment) Run() (e error) {

	var (
		cmd = ssh.Command{
			Client: d.SSH,
		}
		clone = "run git clone " + d.Project.Repo + " /tmp/apker/"
	)

	// Install requirements
	for _, command := range d.Project.Config.Setup {

		cmd.Command = command

		if e = d.exec("Setup: "+command, &cmd); e != nil {
			return
		}
	}

	// Ok! pre deploy commands makes our sexy machine ready 💦.
	// You know what let's me check you first
	cmd.Command = "which git && git --version"

	if e = d.exec("Requirements", &cmd); e != nil {
		return
	}

	// Hmmm, let's do it 😋
	for _, step := range append([]string{clone}, d.Project.Config.Deploy.Steps...) {

		if cmd.Command, e = stepToCommand(step); e != nil {
			return
		}

		if e = d.exec(step, &cmd); e != nil {
			return
		}
	}

	return nil
}

func (d *Deployment) exec(label string, cmd *ssh.Command) (e error) {

	var result ssh.Result

	d.StdoutHandler("", bytes.NewBufferString(cmd.Command))

	if result, e = cmd.Run(); e != nil {

		d.StderrHandler("✔ "+label, &result.Stdout)
		return
	}

	d.StdoutHandler(label, &result.Stdout)
	return
}

func stepToCommand(step string) (c string, e error) {

	parts := strings.Split(step, " ")

	switch parts[0] {
	case "run":
		c = strings.Join(parts[1:], " ")
		break

	case "copy":
		// TODO: backslash (") chars in parts[n]
		c = `cd /tmp/apker && mkdir -p "` + parts[2] + `" && cp ` + parts[1] + ` "` + parts[2] + `"`
		break

	case "reboot":
		c = "reboot &"
		break
	}

	return
}
