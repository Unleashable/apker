// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package utils

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func GitFile(repo string, file string, auth string) (content []byte, e error) {

	url, e := url.Parse(repo)

	if e != nil {
		return
	}

	switch url.Host {
	case "github.com":
		content, e = getGithubFile(url, file, auth)
		return
	case "bitbucket.org":
		content, e = getBitbucketFile(url, file, auth)
		return
	}

	e = fmt.Errorf("Unknown host: %s", url.Host)

	return
}

func getGithubFile(url *url.URL, file string, auth string) (content []byte, e error) {

	content, e = GetContentFromUrl(fmt.Sprintf("https://api.github.com/repos%s/contents/%s", url.Path, file), auth)

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

	content, e = GetContentFromUrl(dl.URL, auth)
	return
}

func getBitbucketFile(url *url.URL, file string, auth string) ([]byte, error) {

	// IDK how to get default branch, using master as default!
	return GetContentFromUrl(fmt.Sprintf("https://api.bitbucket.org/2.0/repositories%s/src/master/%s", url.Path, file), auth)
}
