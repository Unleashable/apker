// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package utils

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func GitFile(repo string, file string) (content []byte, e error) {

	url, e := url.Parse(repo)

	if e != nil {
		return
	}

	switch url.Host {
	case "github.com":
		content, e = getGithubFile(url, file)
		return
	}

	e = fmt.Errorf("Unknown host: %s", url.Host)

	return
}

func getGithubFile(url *url.URL, file string) (content []byte, e error) {

	content, e = GetContentFromUrl(fmt.Sprintf("https://api.github.com/repos%s/contents/%s", url.Path, file))

	if e != nil {
		return
	}

	var dl struct {
		URL string `json:"download_url"`
	}

	e = json.Unmarshal(content, &dl)

	if e != nil {
		return
	}

	content, e = GetContentFromUrl(dl.URL)
	return
}
