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
	Path string
}

type Project struct {
	*Config
	Addr       string
	User       string
	Repo       string
	Auth       string
	Name       string
	Path       string
	Temp       string
	SSHAuth    goph.Auth
	PublicKey  PublicSSHKey
	PrivateKey PrivateSSHKey
}

func (project *Project) Deploy(allowEvents bool, outHandler OutputHandler, errHandler OutputHandler, itHandler ProgressHandler) error {

	if project.User == "" {
		project.User = "root"
	}

	client, e := goph.NewUnknown(project.User, project.Addr, project.SSHAuth)

	if e != nil {
		return e
	}

	deployment := &Deployment{
		SSH:             client,
		Project:         project,
		StdoutHandler:   outHandler,
		StderrHandler:   errHandler,
		ProgressHandler: itHandler,
	}

	var out []byte

	if e = deployment.Run(); e != nil {

		if allowEvents && project.Config.Events.Failure != "" {

			// Run failure event
			out, _ = utils.Run("sh", []string{"-c", project.Config.Events.Failure})
			errHandler("Event: failure", out)
		}

	} else if allowEvents && project.Config.Events.Success != "" {

		// Run success event
		out, e = utils.Run("sh", []string{"-c", project.Config.Events.Success})
		outHandler("Event: success", out)
	}

	return e
}
