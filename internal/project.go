// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package internal

import (
	"github.com/melbahja/goph"
	"github.com/unleashable/apker/internal/utils"
)

type PublicSSHKey struct {
	Fingerprint string
	Path        string
}

type PrivateSSHKey struct {
	Passphrase string
	Path       string
}

type Project struct {
	Config
	Addr       string
	User       string
	Repo       string
	Auth       string
	Name       string
	Path       string
	Temp       string
	PublicKey  PublicSSHKey
	PrivateKey PrivateSSHKey
}

func (project *Project) Deploy(allowEvents bool, outHandler OutputHandler, errHandler OutputHandler, itHandler InteractiveHandler) error {

	if project.User == "" {
		project.User = "root"
	}

	client, e := goph.NewUnknown(project.User, project.Addr, goph.Key(project.PrivateKey.Path, project.PrivateKey.Passphrase))

	if e != nil {
		return e
	}

	deployment := &Deployment{
		SSH:                client,
		Project:            project,
		StdoutHandler:      outHandler,
		StderrHandler:      errHandler,
		InteractiveHandler: itHandler,
	}

	var out []byte

	if e = deployment.Run(); e != nil {

		if allowEvents && project.Config.Events.Error != "" {

			// Run error event
			out, _ = utils.Run("sh", []string{"-c", project.Config.Events.Error})
			errHandler("Event: error", out)
		}

	} else if allowEvents && project.Config.Events.Done != "" {

		// Run done event
		out, e = utils.Run("sh", []string{"-c", project.Config.Events.Done})
		outHandler("Event: done", out)
	}

	return e
}
