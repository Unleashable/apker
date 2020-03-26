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
	"github.com/unleashable/apker/cmd/outputs"
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
	&cli.StringFlag{
		Name:    "url",
		Aliases: []string{"repo"},
		Usage:   "Set project git repository url.",
	},
	&cli.BoolFlag{
		Name:    "passphrase",
		Usage:   "Ask for private key passphrase for protected keys.",
	},
	&cli.DurationFlag{
		Name:    "timeout",
		Aliases: []string{"t"},
		Usage:   "Set timeout duration (e.g: 30s) to wait for image to be ready.",
	},
	&cli.BoolFlag{
		Name:    "no-timeout-error",
		Aliases: []string{"nt"},
		Usage:   "When timeout exit without error code.",
	},
	&cli.IntFlag{
		Name:  "id",
		Usage: "Run deploy steps on specific machine/droplet `id`.",
	},
	&cli.IntFlag{
		Name:  "image",
		Usage: "Create machine/droplet from specific image `id`.",
	},
	// &cli.StringFlag{
	// 	Name:    "addr",
	// 	Aliases: []string{"ip"},
	// 	Usage:   "Deploy to already existing machine.",
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
		Repo: c.String("url"),
		Auth: os.Getenv("APKER_AUTH"),
	}

	// Housekeeping.
	defer os.RemoveAll(project.Temp)

	// Work from remote repo
	if project.Repo != "" {

		var bytes []byte

		// Get content of apker.yaml
		if bytes, e = utils.GitFile(project.Repo, "apker.yaml", project.Auth); e != nil {
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

	// TODO: add step here to validate all config
	// Image url or distro name is required
	if project.Config.Image.From == "" {

		e = errors.New("Image name or url is required.")
		return
	}

	// TODO: instead of this, add a func in project to Set and resolve defaults
	// TODO: add ability to connect with password
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
	// SSH user
	project.User = c.String("user")

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

		outputs.Success("Name: "+project.Name, "")
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

		project.PrivateKey.Passphrase, e = inputs.Password("Private key passphrase", func(p string) error {

			if utils.IsValidPassphrase(project.PrivateKey.Path, p) {
				return nil
			}

			return errors.New("✘ Invalid passphrase!")
		})

		if e != nil {
			return
		}
	}

	// Install image on digitalocean.
	sp = inputs.Spinner(" Droplet setup...")

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

					outputs.Success("Droplet image created.", "")
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

			// By default it's 0 no timeout!
			if int(c.Duration("timeout")) != 0 {

				sp.Stop()
				fmt.Printf("⌛ You can run: '%s --image %d' when image is ready.", strings.Join(os.Args, " "), do.ImageID)

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
	outputs.Success("Droplet now ready.", "")

	// Wait for ssh port
	sp.Suffix = " Waiting for ssh port..."
	sp.Start()

	for {

		time.Sleep(5 * time.Second)

		if utils.IsPortOpen(project.Addr+":22", 5) {
			break
		}
	}

	// Deploy steps spinner!
	sp.Suffix = " Running deploy steps..."

	e = project.Deploy(stdout(sp), stderr(sp))

	sp.Stop()

	if e == nil {

		outputs.Success("It's 👏 Deployed 👏 successfully🚀!", "")
	}

	return
}

func stdout(sp *sp.Spinner) internal.OutputHandler {

	return func(label string, log *bytes.Buffer) error {

		sp.Stop()

		if label == "" {

			sp.Suffix = " " + log.String()

		} else {

			outputs.Success(label, "")

			if log := log.String(); log != "" {
				fmt.Println(log)
			}
		}

		sp.Start()
		return nil
	}
}

func stderr(sp *sp.Spinner) internal.OutputHandler {

	return func(label string, log *bytes.Buffer) error {
		sp.Stop()
		outputs.Error(label, "")

		if log := log.String(); log != "" {
			fmt.Println(log)
		}

		return nil
	}
}
