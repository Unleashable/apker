// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package inputs

import (
	"context"
	"fmt"

	"github.com/melbahja/promptui"
	"github.com/unleashable/apker/internal/providers"
)

func SetDropletSize(do *providers.Digitalocean, size string) error {

	if size == "" {
		size = do.Project.Config.Image.Size
	}

	switch size {

	case "":
		return DigitaloceanSelectSize(do)

	case "small":
		do.Project.Config.Image.Size = "s-1vcpu-1gb"
		break

	default:
		do.Project.Config.Image.Size = size
	}

	return nil
}

func SetDropletRegion(do *providers.Digitalocean, region string) error {

	if region != "" {

		do.Project.Config.Image.Region = region
		return nil
	}

	return DigitaloceanSelectRegion(do)
}

func DigitaloceanSelectSize(do *providers.Digitalocean) error {

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

func DigitaloceanSelectRegion(do *providers.Digitalocean) error {

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
