// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package inputs

import (
	"time"

	"github.com/briandowns/spinner"
)

func Spinner(msg string) *spinner.Spinner {

	s := spinner.New(spinner.CharSets[41], 100*time.Millisecond)
	s.Suffix = msg
	s.Start()
	return s
}
