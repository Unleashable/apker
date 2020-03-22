// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package main

import (
	"log"
	"os"

	"github.com/unleashable/apker/cmd"
)

func main() {

	if e := cmd.Run(os.Args); e != nil {
		log.Fatal(e)
	}
}
