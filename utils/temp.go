// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package utils

import (
	"io/ioutil"
	"os"

	"github.com/rs/xid"
)

func Temp() string {

	tmp, e := ioutil.TempDir(os.TempDir(), "apker-"+xid.New().String())

	if e != nil {
		panic(e)
	}

	return tmp
}
