// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func GetContentFromUrl(url string, auth string) (content []byte, e error) {

	res, e := http.Get(UrlAuth(url, auth))

	if e != nil {
		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {

		e = fmt.Errorf("Request error: %v for %s", res.StatusCode, url)
		return
	}

	content, e = ioutil.ReadAll(res.Body)

	return
}

func IsUrl(u string) bool {

	if u == "" {
		return false
	}

	_, err := url.ParseRequestURI(u)

	return err == nil
}

func UrlAuth(url string, auth string) string {

	if auth != "" {
		return strings.Replace(url, "https://", fmt.Sprintf("https://%s@", auth), 1)
	}

	return url
}
