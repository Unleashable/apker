// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package utils

import (
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

func GetFingerprint(keyPath string, algro string) (string, error) {

	var (
		e    error
		file []byte
	)

	if file, e = ioutil.ReadFile(keyPath); e != nil {

		return "", e

	} else if key, _, _, _, e := ssh.ParseAuthorizedKey(file); e == nil {

		switch algro {

		case "md5":
			return ssh.FingerprintLegacyMD5(key), nil

		case "sha256":
			return ssh.FingerprintSHA256(key), nil
		}
	}

	return "", e
}
