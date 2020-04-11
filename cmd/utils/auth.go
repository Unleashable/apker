package utils

import (
	"fmt"

	"github.com/melbahja/goph"
	"github.com/unleashable/apker/cmd/inputs"
	"github.com/unleashable/apker/internal"
	"github.com/unleashable/apker/internal/utils"
	"github.com/urfave/cli/v2"
)

func SetAuthMethod(project *internal.Project, c *cli.Context) (e error) {

	var pass string

	if c.Bool("password") {

		pass, e = inputs.Password("Enter ssh password", func(p string) error {

			if len(p) < 1 {
				return fmt.Errorf("Invalid password!")
			}

			return nil
		})

		if e != nil {
			return
		}

		fmt.Println("")
		project.SSHAuth = goph.Password(pass)
		return
	}

	project.PublicKey.Fingerprint,
		project.PublicKey.Path,
		project.PrivateKey.Path,
		e = utils.ResolveSSHKeys(c.String("pub"), c.String("key"))

	if e != nil {
		return
	}

	if c.Bool("passphrase") {

		pass, e = inputs.Password("Enter private key passphrase", func(p string) error {

			if utils.IsValidPassphrase(project.PrivateKey.Path, p) {
				return nil
			}

			return fmt.Errorf("Invalid passphrase!")
		})

		if e != nil {
			return
		}

		fmt.Println("")
	}

	project.SSHAuth = goph.Key(project.PrivateKey.Path, pass)
	return
}
