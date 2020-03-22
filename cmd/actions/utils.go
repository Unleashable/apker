// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package actions

import (
	"context"
	"fmt"
	"time"

	sp "github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
	"github.com/unleashable/apker/core/providers"
)

func askString(label string, def string) (v string, e error) {

	prompt := promptui.Prompt{
		Label: label,
	}

	v, e = prompt.Run()

	if e == nil && v == "" {
		v = def
	}

	return
}

func spinner(msg string) *sp.Spinner {

	s := sp.New(sp.CharSets[41], 100*time.Millisecond)
	s.Suffix = msg
	s.Start()
	return s
}

func setDropletSize(do *providers.Digitalocean, size string) error {

	switch size {

	case "":
		return digitaloceanSelectSize(do)

	case "small":
		do.Project.Config.Image.Size = "s-1vcpu-1gb"
		break

	default:
		do.Project.Config.Image.Size = size
	}

	return nil
}

func setDropletRegion(do *providers.Digitalocean, region string) error {

	if region != "" {

		do.Project.Config.Image.Region = region
		return nil
	}

	return digitaloceanSelectRegion(do)
}

func digitaloceanSelectSize(do *providers.Digitalocean) error {

	ctx := context.TODO()

	sizes, _, e := do.DoClient.Sizes.List(ctx, nil)

	if e != nil {
		return e
	}

	var items []string

	for _, size := range sizes {
		items = append(items, fmt.Sprintf(
			"Mem: %dM, Vcpus: %d, Disk: %dG, Price: %.2f/mo",
			size.Memory,
			size.Vcpus,
			size.Disk,
			size.PriceMonthly,
		))
	}

	prompt := promptui.Select{
		Label: "• Select Size",
		Items: items,
	}

	i, _, e := prompt.Run()

	if e != nil {
		return e
	}

	do.Project.Config.Image.Size = sizes[i].Slug

	return nil
}

func digitaloceanSelectRegion(do *providers.Digitalocean) error {

	ctx := context.TODO()

	regions, _, e := do.DoClient.Regions.List(ctx, nil)

	if e != nil {
		return e
	}

	var items []string

	for _, r := range regions {

		if r.Available == false {
			continue
		}

		items = append(items, r.Name)
	}

	prompt := promptui.Select{
		Label: "• Select Region",
		Items: items,
	}

	i, _, e := prompt.Run()

	if e != nil {
		return e
	}

	do.Project.Config.Image.Region = regions[i].Slug

	return nil
}
