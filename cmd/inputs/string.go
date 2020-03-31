// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package inputs

import (
	"github.com/melbahja/promptui"
)

func AskString(label string, def string) (v string, e error) {

	prompt := promptui.Prompt{
		Label: label,
	}

	v, e = prompt.Run()

	if e == nil && v == "" {
		v = def
	}

	return
}
