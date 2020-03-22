// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package actions

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	sp "github.com/briandowns/spinner"
	"github.com/unleashable/apker/core"
	"github.com/unleashable/apker/core/providers"
	"github.com/unleashable/apker/utils"
	"github.com/urfave/cli/v2"
	// "gopkg.in/src-d/go-git.v4"
)

func Deploy(c *cli.Context) (e error) {

	// Where are we?
	cwd, e := os.Getwd()

	if e != nil {
		return
	}

	// Init core.Project with the current working directory
	project := core.Project{
		Path: cwd,
		Temp: utils.Temp(),
	}

	// Housekeeping.
	defer os.RemoveAll(project.Temp)

	// Work from remote repo
	if repo := c.Args().First(); repo != "" {

		var bytes []byte

		// Get content of apker.yaml
		if bytes, e = utils.GitFile(repo, "apker.yaml"); e != nil {
			return
		}

		// Saving apker.yaml
		if e = ioutil.WriteFile(project.Temp+"/apker.yaml", bytes, 0644); e != nil {
			return
		}

		// For finding apker.yaml
		project.Path = project.Temp

		// Clone the repo to temp
		// _, e = git.PlainClone(project.Path, false, &git.CloneOptions{
		// 	URL: repo,
		// })

		// // Is everything okay in there?
		// if e != nil {
		// 	return
		// }
	}

	// Load project config from apker.yaml
	project.Config, e = core.LoadConfig(project.Path)

	// Are we okay?
	if e != nil {
		return
	}

	// Image url or distro name is required
	if project.Config.Image.URL == "" && project.Config.Image.Distro == "" {

		e = errors.New("Image destro name or url is required.")
		return
	}

	// Set ssh keys path
	if e = core.ResolveSSHKeys(&project, c.String("pub"), c.String("key")); e != nil {
		return
	}

	// Project name
	project.Name = c.String("name")

	if project.Name == "" && project.Config.Name != "" {
		project.Name = "apker-image-" + project.Config.Name
	}

	switch project.Config.Provider.Name {
	case "digitalocean":
		e = digitaloceanDeploy(&project, c)
		break
	default:
		e = errors.New("Unknown provider name: " + project.Config.Provider.Name)
	}

	return
}

func digitaloceanDeploy(project *core.Project, c *cli.Context) (e error) {

	var (
		sp          *sp.Spinner
		do          *providers.Digitalocean
		machine     core.MachineStatus
		MachineChan chan core.MachineStatus
		// installTimeout <-chan time.Time
		machineIpAddr     string
		imageAlredyExists bool = c.Int("image") != 0 || c.Int("id") != 0
	)

	if do, e = providers.NewDigitalocean(project); e != nil {
		return
	}

	if imageAlredyExists {

		goto DropletSetup

	} else if project.Name == "" {

		// ask for image name
		project.Name, e = askString("Name: ", fmt.Sprintf("apker-image-%v", time.Now().Unix()))

		if e != nil {
			return
		}

	} else {

		fmt.Println("Name: " + project.Name)
	}

	// Droplet size and region
	if e = setDropletSize(do, c.String("size")); e != nil {

		return

	} else if e = setDropletRegion(do, c.String("region")); e != nil {

		return
	}

DropletSetup:

	// Install image on digitalocean.
	sp = spinner(" Droplet...")

	// Wait for install step for 200s
	// installTimeout = time.After(10 * time.Second)

	// Installation channel
	MachineChan = make(chan core.MachineStatus)

	// Go setup Image and droplet
	go do.SetupMachine(MachineChan, core.Attributes{
		"imageId":   c.Int("image"),
		"dropletId": c.Int("id"),
	})

MachineLoop:

	for {
		select {
		case machine = <-MachineChan:

			// Handle status data
			switch true {
			case machine.Error != nil:

				e = machine.Error
				break MachineLoop

			case !machine.IsImageInstalled && machine.IsImageReady:

				sp.Stop()

				if imageAlredyExists == false {

					fmt.Println("✔ Droplet image created.")
				}

				sp = spinner(" Cheking droplet...")
				break

			case machine.IsMachineReady:

				machineIpAddr = machine.Addr
				break MachineLoop

			default:

				if machine.IsMachineReady == false && machine.IsImageReady == false {

					sp.Suffix = " Current machine status: " + machine.Status
				}
			}

			// case <-installTimeout:

			// TODO: handle timeout
			// if c.Bool("wait") == false {

			// 	sp.Stop()
			// 	fmt.Sprintf("⌛ You can [CTRL + C] now and later run: apker deply --image-id %d", do.ImageID)
			// 	sp.Start()
			// }
		}
	}

	sp.Stop()
	close(MachineChan)

	if e != nil {
		return
	}

	// Now we have a droplet ready for action
	fmt.Println("✔ Droplet now ready.")

	// Deploy:

	// TODO: wait for ssh port

	// Deploy steps spinner!
	sp = spinner(" Running deploy steps...")

	e = project.Deploy(c.String("user"), machineIpAddr, c.String("passphrase"), stdout(sp), stderr(sp))

	sp.Stop()

	if e == nil {

		fmt.Println("✔ It's 👏 Deployed 👏 successfully🚀!")
	}

	return
}

func stdout(sp *sp.Spinner) core.OutputHandler {

	return func(step string, log *bytes.Buffer) error {
		sp.Stop()
		fmt.Println("✔ " + step)
		fmt.Println(log.String())
		sp.Start()
		return nil
	}
}

func stderr(sp *sp.Spinner) core.OutputHandler {

	return func(step string, log *bytes.Buffer) error {
		sp.Stop()
		fmt.Println("✘ " + step)
		fmt.Println(log.String())
		return nil
	}
}
