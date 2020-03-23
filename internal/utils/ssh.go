// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package utils

import (
	"io/ioutil"
	"net"
	"os"
	"path"
	"time"

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

func ResolveSSHKeys(publicKeyPath string, privateKeyPath string) (pubfp string, pub string, prv string, e error) {

	if publicKeyPath == "" {
		publicKeyPath = path.Join(os.Getenv("HOME"), ".ssh", "id_rsa.pub")
	}

	if privateKeyPath == "" {
		privateKeyPath = path.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
	}

	if pubfp, e = GetFingerprint(publicKeyPath, "md5"); e != nil {

		return

	} else if _, e = os.Stat(privateKeyPath); e != nil {

		return
	}

	pub = publicKeyPath
	prv = privateKeyPath
	return
}

func IsPortOpen(addr string, tSeconds int) bool {

	_, e := net.DialTimeout("tcp", addr, time.Duration(tSeconds)*time.Second)

	return e == nil
}

func IsValidPassphrase(prv string, passphrase string) bool {

	var (
		e          error
		privateKey []byte
	)

	if privateKey, e = ioutil.ReadFile(prv); e != nil {
		return false
	}

	_, e = ssh.ParsePrivateKeyWithPassphrase(privateKey, []byte(passphrase))

	return e == nil
}
