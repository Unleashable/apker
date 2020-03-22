// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package core

import (
	"os"
	// "fmt"
	"bytes"
	"path"
	"strings"

	"github.com/melbahja/ssh"
	"github.com/unleashable/apker/utils"
)

type OutputHandler func(step string, log *bytes.Buffer) error

// TODO: Refactor this!
type Deployment struct {
	SSH           *ssh.Client
	Project       *Project
	StdoutHandler OutputHandler
	StderrHandler OutputHandler
}

func (d *Deployment) Run(preDeploy []string, steps []string) (e error) {

	cmd := ssh.Command{
		Client: d.SSH,
	}

	//
	// Install requirements
	for _, command := range preDeploy {

		cmd.Command = command

		if e = d.exec("Setup: "+command, &cmd); e != nil {
			return
		}
	}

	// Ok! preDeploy commands makes our sexy machine ready 💦.
	//
	// You know what let's me check you first
	cmd.Command = "which git"

	if e = d.exec("Requirements", &cmd); e != nil {
		return
	}

	// Hmmm, let's do it 😋
	// cmd.Command = "git clone "

	for _, step := range steps {

		if cmd.Command, e = stepToCommandLine(step); e != nil {
			return
		}

		if e = d.exec(step, &cmd); e != nil {
			return
		}
	}

	return nil
}

func ResolveSSHKeys(project *Project, publicKeyPath string, privateKeyPath string) (e error) {

	if publicKeyPath == "" {
		publicKeyPath = path.Join(os.Getenv("HOME"), ".ssh", "id_rsa.pub")
	}

	if privateKeyPath == "" {
		privateKeyPath = path.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
	}

	if project.PublicKey.Fingerprint, e = utils.GetFingerprint(publicKeyPath, "md5"); e != nil {

		return

	} else if _, e = os.Stat(privateKeyPath); e != nil {

		return
	}

	project.PublicKey.Path = publicKeyPath
	project.PrivateKey.Path = privateKeyPath
	return nil
}

func (d *Deployment) exec(label string, cmd *ssh.Command) (e error) {

	var result ssh.Result

	if result, e = cmd.Run(); e != nil {

		d.StderrHandler(label, &result.Stdout)
		return
	}

	d.StdoutHandler(label, &result.Stdout)
	return
}

func stepToCommandLine(step string) (c string, e error) {

	parts := strings.Split(step, " ")

	switch parts[0] {
	case "run":
		c = strings.Join(parts[1:], " ")
		break

	case "copy":
		// TODO: backslash (") chars in parts[1]
		//
		panic()
		// c = `mkdir -p "` + parts[1] + `" && cp`
		break

	case "reboot":
		c = "reboot &"
		break
	}

	return
}
