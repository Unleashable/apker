package outputs

import (
	"fmt"
	"github.com/melbahja/promptui"
)

func Success(msg string, icon string) {

	if icon == "" {
		icon = "✔"
	}

	fmt.Println(Render(icon, msg, true))
}

func Input(msg string, icon string) {

	if icon == "" {
		icon = "🔑"
	}

	fmt.Printf("%s: ", Render(icon, msg, true))
}

func Error(msg string, icon string) {

	if icon == "" {
		icon = "✗"
	}

	fmt.Println(Render(icon, msg, false))
}

func Render(icon string, msg string, isSuccess bool) string {

	if isSuccess {
		return promptui.Styler(promptui.FGGreen)(icon) + " " + promptui.Styler(promptui.FGFaint)(msg)
	}

	return promptui.Styler(promptui.FGRed)(icon) + " " + promptui.Styler(promptui.FGFaint, promptui.FGBold)(msg)
}
