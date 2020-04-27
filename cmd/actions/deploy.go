// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package actions

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	sp "github.com/briandowns/spinner"
	"github.com/unleashable/apker/cmd/inputs"
	"github.com/unleashable/apker/cmd/outputs"
	. "github.com/unleashable/apker/cmd/utils"
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
		Name:  "passphrase",
		Usage: "Ask for private key passphrase for protected keys.",
	},
	&cli.BoolFlag{
		Name:  "password",
		Usage: "Ask for ssh password instead of using private keys.",
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
	&cli.StringSliceFlag{
		Name:    "parameter",
		Usage:   "Set deploy template parameters.",
		Aliases: []string{"set"},
	},
	&cli.BoolFlag{
		Name:    "events",
		Usage:   "Allow events to execute on local machine.",
		Aliases: []string{"with-events"},
	},
	&cli.StringFlag{
		Name:    "addr",
		Aliases: []string{"ip"},
		Usage:   "Deploy on already exists machine ip address.",
	},
}

func Deploy(c *cli.Context) (e error) {

	// Where are we?
	cwd, e := os.Getwd()

	if e != nil {
		return
	}

	// Init new project with the current working directory
	project := internal.Project{
		Path: cwd,
		Temp: utils.Temp(),
		Repo: c.String("url"),
		Name: c.String("name"),
		User: c.String("user"),
		Addr: c.String("addr"),
		Auth: os.Getenv("APKER_AUTH"),
	}

	// Housekeeping.
	defer os.RemoveAll(project.Temp)

	var tmp []byte

RemoteGetYamlFile:

	// Work from remote repo
	if project.Repo != "" {

		// Get content of apker.yaml
		if tmp, e = utils.GitFile(project.Repo, "apker.yaml", project.Auth); e != nil {
			return
		}

		// Save apker.yaml to temp file
		if e = ioutil.WriteFile(project.Temp+"/apker.yaml", tmp, 0600); e != nil {
			return
		}

	} else {

		// Get remote url
		if tmp, e = utils.Run("git", []string{"config", "--get", "remote.origin.url"}); e != nil {

			e = errors.New("Get remote repository url error: " + e.Error())
			return
		}

		project.Repo = strings.TrimSpace(string(tmp))
		goto RemoteGetYamlFile
	}

	// Override provider name if ip flag has a value.
	if c.String("addr") != "" {

		if e = os.Setenv("APKER_PROVIDER", "custom"); e != nil {
			return
		}
	}

	// Load project config from apker.yaml
	if project.Config, e = internal.LoadConfig(project.Temp, c.StringSlice("parameter")); e != nil {
		return
	}

	// Validate config.
	if e = project.Config.Validate(); e != nil {
		return
	}

	// Project name fallback
	if project.Name == "" && project.Config.Name != "" {
		project.Name = "apker-" + project.Config.Name
	}

	// Set auth method.
	if e = SetAuthMethod(&project, c); e != nil {
		return
	}

	switch project.Config.Provider.Name {
	case "digitalocean":
		e = digitaloceanDeploy(&project, c)
		break
	case "custom":
		e = customDeploy(&project, c)
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

	// Install image on digitalocean.
	sp = outputs.Spinner(" Droplet setup...")

	// Timeout for install step for
	installTimeout = time.After(c.Duration("timeout"))

	// Installation channel
	MachineChan = make(chan internal.MachineStatus)
	defer close(MachineChan)

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

				sp = outputs.Spinner(" Cheking droplet...")
				break

			case machine.IsMachineReady:

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

	// Run deployment
	e = runDeploy(project, sp, c.Bool("events"))
	return
}

func customDeploy(project *internal.Project, c *cli.Context) (e error) {

	sp := outputs.Spinner("Start...")

	if project.Addr = c.String("addr"); project.Addr == "" {

		sp.Stop()
		return errors.New("Please set machine ip address via 'addr' flag.")
	}

	return runDeploy(project, sp, c.Bool("events"))
}

func runDeploy(project *internal.Project, sp *sp.Spinner, events bool) (e error) {

	// Deploy spinner!
	sp.Suffix = " Running deploy..."

	e = project.Deploy(events, stdout(sp), stderr(sp), progress(sp))

	sp.Stop()

	if e == nil {

		outputs.Success("It's 👏 Deployed 👏 successfully🚀!", "")
	}

	return
}

func progress(sp *sp.Spinner) internal.ProgressHandler {

	return func(log string) error {

		sp.Suffix = " " + log
		return nil
	}
}

func stdout(sp *sp.Spinner) internal.OutputHandler {

	return func(label string, log []byte) error {

		sp.Stop()

		outputs.Success(label, "")

		// TODO: redirect stdout to a file by a flag. default /dev/null
		// if log := string(log); log != "" {
		// 	fmt.Println(log)
		// }

		sp.Start()
		return nil
	}
}

func stderr(sp *sp.Spinner) internal.OutputHandler {

	return func(label string, log []byte) error {

		sp.Stop()

		outputs.Error(label, "")

		if log := string(log); log != "" {
			fmt.Println(log)
		}

		return nil
	}
}
