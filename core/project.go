// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package core

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
	Name       string
	Path       string
	Temp       string
	Config     Config
	PublicKey  PublicSSHKey
	PrivateKey PrivateSSHKey
}

func (project *Project) Deploy(user string, addr string, passphrase string, outHandler OutputHandler, errHandler OutputHandler) error {

	if user == "" {
		user = "root"
	}

	client, e := ssh.New(ssh.Config{
		User: user,
		Addr: addr,
		Config: &sh.ClientConfig{
			Timeout: 20 * time.Second,
		},
		Auth: ssh.Key(project.PrivateKey.Path, passphrase),
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

	if e = deployment.Run(project.Config.Setup, project.Config.Deploy.Steps); e != nil {

		fmt.Println(project.Config.Events.Error)
		return e
	}

	fmt.Println(project.Config.Events.Done)

	return nil
}
