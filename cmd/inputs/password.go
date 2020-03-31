package inputs

import (
	"fmt"
	"github.com/melbahja/promptui"
)

func Password(label string, validator promptui.ValidateFunc) (string, error) {

	prompt := promptui.Prompt{
		Label:          label,
		Mask:           '*',
		Validate:       validator,
		LazyValidation: true,
		Templates: &promptui.PromptTemplates{
			Success: fmt.Sprintf(`%s {{faint "Password:"}} `, promptui.IconGood),
		},
	}

	return prompt.Run()
}
