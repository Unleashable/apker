package inputs

import (
	"github.com/manifoldco/promptui"
)

func Password(label string) (string, error) {

	prompt := promptui.Prompt{
		Label: label,
		Mask:  '*',
	}

	return prompt.Run()
}
