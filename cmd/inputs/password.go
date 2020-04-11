package inputs

import (
	"fmt"

	"github.com/unleashable/apker/cmd/outputs"
	"golang.org/x/crypto/ssh/terminal"
)

func Password(label string, validator func(string) error) (string, error) {

	var (
		err   error
		bpass []byte
		spass string
	)

	for {

		outputs.Input(label, "")

		if bpass, err = terminal.ReadPassword(0); err != nil {
			return "", err
		}

		spass = string(bpass)

		if err = validator(spass); err == nil {
			break
		}

		fmt.Println("")
		outputs.Error(err.Error(), "")
		err = nil
	}

	return spass, nil
}
