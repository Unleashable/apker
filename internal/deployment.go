// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/melbahja/goph"
	"github.com/unleashable/apker/internal/utils"
)

type OutputHandler func(description string, log []byte) error

type InteractiveHandler func(log string) error

type ExecStep struct {
	Done    string
	Label   string
	Command string
}

type Deployment struct {
	SSH                *goph.Client
	Project            *Project
	StdoutHandler      OutputHandler
	StderrHandler      OutputHandler
	InteractiveHandler InteractiveHandler
}

func (d *Deployment) Run() (e error) {

	var (
		command string
		steps   []ExecStep = []ExecStep{}
	)

	for _, command = range d.Project.Config.Setup {

		steps = append(steps, ExecStep{
			Done:    fmt.Sprintf("Setup: %s", command),
			Label:   fmt.Sprintf("Running: %s", command),
			Command: command,
		})
	}

	steps = append(steps, ExecStep{
		Done:    "Requirements: git and rsync installed.",
		Label:   "Verifying requirements...",
		Command: "which git rsync && git --version && rsync --version",
	}, ExecStep{
		Done:    "Project cloned successfully on: /tmp/apker",
		Label:   fmt.Sprintf("Cloning project repository: %s", d.Project.Repo),
		Command: fmt.Sprintf("git clone %s /tmp/apker/", utils.UrlAuth(d.Project.Repo, d.Project.Auth)),
	}, ExecStep{
		Done:    "Config file copied successfully.",
		Label:   "Coping config file",
		Command: "mkdir -p /usr/share/apker/ && cp /tmp/apker/apker.yaml /usr/share/apker/apker.yaml",
	})

	for _, step := range d.Project.Config.Deploy.Steps {

		if command, e = stepToCommand(step); e != nil {
			return
		}

		steps = append(steps, ExecStep{
			Done:    fmt.Sprintf("Step: %s", step),
			Label:   fmt.Sprintf("Running: %s", step),
			Command: command,
		})
	}

	return d.exec(steps)
}

func (d Deployment) exec(steps []ExecStep) (e error) {

	var result []byte

	for _, step := range steps {

		d.InteractiveHandler(step.Label)

		if result, e = d.SSH.Run(step.Command); e != nil {

			d.StderrHandler(step.Command, result)
			return
		}

		d.StdoutHandler(step.Done, result)
	}

	return
}

func stepToCommand(step string) (c string, e error) {

	parts := strings.Split(step, " ")

	switch parts[0] {
	case "run":
		c = fmt.Sprintf("cd /tmp/apker && %s", strings.Join(parts[1:], " "))
		break

	case "copy":
		c = fmt.Sprintf("cd /tmp/apker && rsync -av --quiet %s %s", strconv.Quote(parts[1]), strconv.Quote(parts[2]))
		break

	case "mkdir":
		c = fmt.Sprintf("mkdir -p %s", strings.Join(parts[1:], " "))
		break

	case "reboot":
		c = "reboot &"
		break
	}

	e = fmt.Errorf("Unknown step: %s", step)
	return
}
