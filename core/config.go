// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package core

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/unleashable/apker/utils"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Version string `yaml:"version"`
	Name    string `yaml:"name"`
	Image   struct {
		URL    string `yaml:"url"`
		Size   string `yaml:"size"`
		Region string
		Distro string `yaml:"distro"`
	} `yaml:"image"`
	Provider struct {
		Name        string            `yaml:"name"`
		Credentials map[string]string `yaml:"credentials"`
	} `yaml:"provider"`
	Setup  []string `yaml:"setup"`
	Deploy struct {
		Steps []string `yaml:"steps"`
	} `yaml:"deploy"`
	Actions map[string]string `yaml:"actions"`
	Events  struct {
		Error string `yaml:"error"`
		Done  string `yaml:"done"`
	} `yaml:"events"`
}

func isRequired(k string) bool {

	switch k {
	case "APKER_PROVIDER", "APKER_KEY":
		return true
	}

	return false
}

func parseTpl(file string, name string) (c string, e error) {

	tpl, e := template.New(name).Funcs(template.FuncMap{
		"Env": func(key string) string {
			return utils.Env(key, isRequired(key))
		},
	}).ParseFiles(file)

	if e != nil {
		return
	}

	var buf bytes.Buffer

	e = tpl.Execute(&buf, "")

	if e != nil {
		return
	}

	c = buf.String()
	return
}

func checkSteps(steps []string) error {

	var step []string

loop:

	for k, v := range steps {

		step = strings.Split(v, " ")

		switch step[0] {

		case "run", "copy":
			continue loop

		case "reboot":
			if k != len(steps)-1 {
				return fmt.Errorf("reboot command should be the last step.")
			}
		default:
			return fmt.Errorf("Unknown deploy step: %s", v)
		}
	}

	return nil
}

func LoadConfig(projectDirectory string) (c Config, e error) {

	var data string

	if data, e = parseTpl(projectDirectory+"/apker.yaml", "apker.yaml"); e != nil {

		return

	} else if e = yaml.Unmarshal([]byte(data), &c); e != nil {

		return
	}

	e = checkSteps(c.Deploy.Steps)
	return
}
