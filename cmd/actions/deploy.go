// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package actions

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	sp "github.com/briandowns/spinner"
	"github.com/unleashable/apker/cmd/inputs"
	"github.com/unleashable/apker/internal"
	"github.com/unleashable/apker/internal/providers"
	"github.com/unleashable/apker/internal/utils"
	"github.com/urfave/cli/v2"
)

var DeployFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  "name",
		Usage: "Set machine name.",
	},
	&cli.StringFlag{
		Name:  "size",
		Usage: "Set machine size.",
	},
	&cli.StringFlag{
		Name:    "region",
		Aliases: []string{"location"},
		Usage:   "Set machine physical location.",
	},
	&cli.BoolFlag{
		Name:    "passphrase",
		Aliases: []string{"pass"},
		Usage:   "Ask for private key passphrase for protected keys.",
	},
	&cli.DurationFlag{
		Name:    "timeout",
		Aliases: []string{"t"},
		Usage:   "Set timeout for apker to wait for image (e.g: 12s).",
	},
	&cli.BoolFlag{
		Name:    "no-timeout-error",
		Aliases: []string{"nt"},
		Usage:   "Exit on timeout without error status.",
	},
	&cli.IntFlag{
		Name:  "id",
		Usage: "Run deploy steps on specific machine/droplet `id`.",
	},
	&cli.IntFlag{
		Name:  "image",
		Usage: "Create machine/droplet from specific image `id`.",
	},
	// TODO
	// &cli.StringFlag{
	// 	Name:    "addr",
	// 	Aliases: []string{"ip"},
	// 	Usage:   "Set ip address of already existing machine.",
	// },
}

func Deploy(c *cli.Context) (e error) {

	// Where are we?
	// cwd, e := os.Getwd()

	// if e != nil {
	// 	return
	// }

	// Init new project with the current working directory
	project := internal.Project{
		// Path: cwd,
		Temp: utils.Temp(),
		Repo: c.Args().First(),
	}

	// Housekeeping.
	defer os.RemoveAll(project.Temp)

	// Work from remote repo
	if project.Repo != "" {

		var bytes []byte

		// Get content of apker.yaml
		if bytes, e = utils.GitFile(project.Repo, "apker.yaml"); e != nil {
			return
		}

		// Saving apker.yaml
		if e = ioutil.WriteFile(project.Temp+"/apker.yaml", bytes, 0644); e != nil {
			return
		}

		// Change project path to temp
		project.Path = project.Temp

	} else {

		// TODO: work from cwd if it is a valid git repo
		e = errors.New("Remote git repo required")
		return
	}

	// Load project config from apker.yaml
	if project.Config, e = internal.LoadConfig(project.Path); e != nil {
		return
	}

	// Image url or distro name is required
	if project.Config.Image.From == "" {

		e = errors.New("Image name or url is required.")
		return
	}

	// Set path of ssh keys
	project.PublicKey.Fingerprint,
		project.PublicKey.Path,
		project.PrivateKey.Path,
		e = utils.ResolveSSHKeys(c.String("pub"), c.String("key"))

	if e != nil {
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

func digitaloceanDeploy(project *internal.Project, c *cli.Context) (e error) {

	var (
		sp               *sp.Spinner
		do               *providers.Digitalocean
		machine          internal.MachineStatus
		MachineChan      chan internal.MachineStatus
		installTimeout   <-chan time.Time
		skipInputPrompts bool = c.Int("image") != 0 || c.Int("id") != 0
	)

	if do, e = providers.NewDigitalocean(project); e != nil {
		return
	}

	if skipInputPrompts {

		goto DropletSetup

	} else if project.Name == "" {

		// ask for image name
		project.Name, e = inputs.AskString("Name: ", fmt.Sprintf("apker-image-%v", time.Now().Unix()))

		if e != nil {
			return
		}

	} else {

		fmt.Println("Name: " + project.Name)
	}

	// Droplet size and region
	if e = inputs.SetDropletSize(do, c.String("size")); e != nil {

		return

	} else if e = inputs.SetDropletRegion(do, c.String("region")); e != nil {

		return
	}

DropletSetup:

	// Ask for ssh key passphrase
	if c.Bool("passphrase") {

		for {

			if project.PrivateKey.Passphrase, e = inputs.Password("Private key passphrase"); e != nil {
				return
			}

			// Validate passphrase before continue
			if utils.IsValidPassphrase(project.PrivateKey.Path, project.PrivateKey.Passphrase) {
				break
			}

			fmt.Println("✘ Invalid passphrase!")
		}
	}

	// Install image on digitalocean.
	sp = inputs.Spinner(" Droplet...")

	// Timeout for install step for
	installTimeout = time.After(c.Duration("timeout"))

	// Installation channel
	MachineChan = make(chan internal.MachineStatus)

	// Go setup Image and droplet
	go do.SetupMachine(MachineChan, internal.Attributes{
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

				if skipInputPrompts == false {

					fmt.Println("✔ Droplet image created.")
				}

				sp = inputs.Spinner(" Cheking droplet...")
				break

			case machine.IsMachineReady:

				// machineIpAddr = machine.Addr
				break MachineLoop

			default:

				if machine.IsMachineReady == false && machine.IsImageReady == false {

					sp.Suffix = " Current droplet status: " + machine.Status
				}
			}

		case <-installTimeout:

			//
			// By default it's 0 no timeout!
			//

			if int(c.Duration("timeout")) != 0 {

				sp.Stop()
				fmt.Printf("⌛ You can run: '%s --image %d' when image ready.", strings.Join(os.Args, " "), do.ImageID)

				if c.Bool("no-timeout-error") == false {

					e = errors.New("Installation timeout")
				}

				return
			}
		}
	}

	sp.Stop()
	close(MachineChan)

	if e != nil {
		return
	}

	// Now we have a droplet ready for action
	fmt.Println("✔ Droplet now ready.")

	// Wait for ssh port
	sp.Suffix = " Waiting for ssh port to open..."
	sp.Start()

	for {

		time.Sleep(5 * time.Second)

		if utils.IsPortOpen(project.AddrV4+":22", 5) {
			break
		}
	}

	// Deploy steps spinner!
	sp.Suffix = " Running deploy steps..."

	e = project.Deploy(c.String("user"), stdout(sp), stderr(sp))

	sp.Stop()

	if e == nil {

		fmt.Println("✔ It's 👏 Deployed 👏 successfully🚀!")
	}

	return
}

func stdout(sp *sp.Spinner) internal.OutputHandler {

	return func(label string, log *bytes.Buffer) error {

		sp.Stop()

		if label == "" {

			sp.Suffix = " " + log.String()

		} else {

			fmt.Println(label)
			fmt.Println(log.String())
		}

		sp.Start()
		return nil
	}
}

func stderr(sp *sp.Spinner) internal.OutputHandler {

	return func(label string, log *bytes.Buffer) error {
		sp.Stop()
		fmt.Println("✘ " + label)
		fmt.Println(log.String())
		return nil
	}
}
