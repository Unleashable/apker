// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package internal

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/melbahja/goph"
	"github.com/unleashable/apker/internal/utils"
)

type OutputHandler func(description string, log []byte) error

type ProgressHandler func(log string) error

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
	ProgressHandler ProgressHandler
}

func (d *Deployment) Run() (e error) {

	var (
		command string
		steps   []ExecStep = []ExecStep{}
	)

	for _, command = range d.Project.Config.Deploy.Setup {

		steps = append(steps, ExecStep{
			Done:    fmt.Sprintf("Setup: %s", command),
			Label:   fmt.Sprintf("Running: %s", command),
			Command: command,
		})
	}

	steps = append(steps, ExecStep{
		Done:    "Setup: git and rsync installed.",
		Label:   "Verifying requirements...",
		Command: "which git rsync && git --version && rsync --version",
	}, ExecStep{
		Done:    "Setup: project cloned on: /tmp/apker",
		Label:   fmt.Sprintf("Cloning project repository: %s", d.Project.Repo),
		Command: fmt.Sprintf("rm -rf /tmp/apker && git clone %s /tmp/apker/", utils.UrlAuth(d.Project.Repo, d.Project.Auth)),
	}, ExecStep{
		Done:    "Setup: apker directory created.",
		Label:   "Creating apker directory...",
		Command: "mkdir -p /usr/share/apker/bin/",
	}, ExecStep{
		Done:    "Setup: apker actions created.",
		Label:   "Creating actions...",
		Command: "chmod +x /tmp/apker_actions.sh && /tmp/apker_actions.sh && chmod +x /usr/share/apker/bin/*",
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

	var (
		env    string = envToString(d.Project.Config.Deploy.Env)
		result []byte
	)

	// Setup actions.
	if e = d.setupActions(); e != nil {
		d.StderrHandler(fmt.Sprintf("Setup actions error: %s", e.Error()), result)
		return
	}

	for _, step := range steps {

		d.ProgressHandler(step.Label)

		if result, e = d.SSH.Run(fmt.Sprintf("env %s bash -c '%s'", env, step.Command)); e != nil {

			d.StderrHandler(fmt.Sprintf("Label: %s\nCommand: %s", step.Label, step.Command), result)
			return
		}

		d.StdoutHandler(step.Done, result)
	}

	return
}

func (d Deployment) setupActions() error {

	var err error

	file, err := os.OpenFile(d.Project.Temp+"/actions.sh", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)

	if err != nil {
		return err
	}

	defer file.Close()

	file.WriteString("#!/usr/bin/sh\n")

	for name, command := range d.Project.Config.Actions {

		_, err = file.WriteString(fmt.Sprintf(`cat > /usr/share/apker/bin/%s << EOL
#!/usr/bin/sh
%s
EOL
`, name, command))

		if err != nil {
			return err
		}
	}

	return d.SSH.Upload(d.Project.Temp+"/actions.sh", "/tmp/apker_actions.sh")
}

func envToString(vars map[string]string) string {

	env := []string{}

	for i := range vars {
		env = append(env, fmt.Sprintf("%s=%s", i, strconv.Quote(vars[i])))
	}

	return strings.Join(env, " ")
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

	case "dir":
		c = fmt.Sprintf("mkdir -p %s", strings.Join(parts[1:], " "))
		break

	case "reboot":
		c = "reboot &"
		break

	default:
		e = fmt.Errorf("Unknown step: %s", step)
	}

	return
}
