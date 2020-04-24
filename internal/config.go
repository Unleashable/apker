// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package internal

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/unleashable/apker/internal/utils"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Version string `yaml:"version"`
	Name    string `yaml:"name"`
	Image   struct {
		Size   string `yaml:"size"`
		From   string `yaml:"from"`
		Region string
	} `yaml:"image"`
	Provider struct {
		Name        string            `yaml:"name"`
		Credentials map[string]string `yaml:"credentials"`
	} `yaml:"provider"`
	Deploy struct {
		Env   map[string]string `yaml:env`
		Setup []string          `yaml:"setup"`
		Steps []string          `yaml:"steps"`
	} `yaml:"deploy"`
	Actions map[string]string `yaml:"actions"`
	Events  struct {
		Failure string `yaml:"failure"`
		Success string `yaml:"success"`
	} `yaml:"events"`
}

func (c Config) Validate() error {

	if c.Image.From == "" {
		return fmt.Errorf("Image name or url is required.")
	}

	return checkSteps(c.Deploy.Steps)
}


func parseParams(params []string) map[string]string {

	paramsMap := make(map[string]string)
	paramSlice := []string{}

	for _, val := range params {

		if paramSlice = strings.Split(val, "="); len(paramSlice) == 2 {
			paramsMap[paramSlice[0]] = paramSlice[1]
			continue
		}

		panic(fmt.Sprintf("Invalid param %s", val))
	}

	return paramsMap
}

func parseTpl(file string, name string, params []string) (c string, e error) {

	paramsMap := parseParams(params)
	tpl, e := template.New(name).Funcs(template.FuncMap{
		"Env": func(key string) string {
			return utils.Env(key, false)
		},
		"Get": func(key string) string {

			if val, ok := paramsMap[key]; ok {
				return val
			}

			panic(fmt.Sprintf("Param: %s is required, add --set %s=YOUR_VALUE", key, key))
		},
		"GetOr": func(key string, def string) string {

			if val, ok := paramsMap[key]; ok {
				return val
			}

			return def
		},
		"Run": func(cmd string) string {

			val, e := utils.Run("bash", []string{"-c", cmd})

			if e != nil {
				panic(e)
			}

			return string(val)
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

// TODO: check steps len
func checkSteps(steps []string) error {

	var step []string

loop:

	for k, v := range steps {

		step = strings.Split(v, " ")

		switch step[0] {

		case "run", "copy", "dir":
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

func LoadConfig(projectDirectory string, params []string) (c *Config, e error) {

	var data string

	if data, e = parseTpl(projectDirectory+"/apker.yaml", "apker.yaml", params); e != nil {
		return
	}

	e = yaml.Unmarshal([]byte(data), &c)
	return
}
