// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package internal

import (
	"fmt"
	"time"

	"github.com/melbahja/ssh"
	sh "golang.org/x/crypto/ssh"
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
	Name       string
	Path       string
	Temp       string
	PublicKey  PublicSSHKey
	PrivateKey PrivateSSHKey
}

func (project *Project) Deploy(outHandler OutputHandler, errHandler OutputHandler) error {

	if project.User == "" {
		project.User = "root"
	}

	client, e := ssh.New(ssh.Config{
		User: project.User,
		Addr: project.Addr,
		Config: &sh.ClientConfig{
			Timeout: 20 * time.Second,
		},
		Auth: ssh.Key(project.PrivateKey.Path, project.PrivateKey.Passphrase),
	})

	if e != nil {
		return e
	}

	deployment := &Deployment{
		SSH:           client,
		Project:       project,
		StdoutHandler: outHandler,
		StderrHandler: errHandler,
	}

	if e = deployment.Run(); e != nil {

		// TODO: exec events
		fmt.Println(project.Config.Events.Error)
		return e
	}

	fmt.Println(project.Config.Events.Done)

	return nil
}
