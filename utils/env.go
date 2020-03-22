// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package utils

import "os"

func Env(k string, required bool) string {

	val := os.Getenv(k)

	if required && val == "" {

		panic("Env: " + k + " is required")
	}

	return val
}
