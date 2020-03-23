// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package providers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/digitalocean/godo"
	"github.com/unleashable/apker/internal"
	"golang.org/x/oauth2"
)

type Digitalocean struct {
	DropletID int
	ImageID   int
	Oauth     *http.Client
	DoClient  *godo.Client
	Project   *internal.Project
}

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {

	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func (do *Digitalocean) SetupMachine(ch chan internal.MachineStatus, attrs internal.Attributes) {

	var (
		e error
		// ctx context.Context
		image        *godo.Image
		droplet      *godo.Droplet
		dropletImage godo.DropletCreateImage = godo.DropletCreateImage{}
	)

	// Giving non 0 value to imageId or dropletId attribute it
	// means that you are sure the image or droplet is available!
	if attrs["imageId"].(int) != 0 || attrs["dropletId"].(int) != 0 {

		do.ImageID, dropletImage.ID = attrs["imageId"].(int), attrs["imageId"].(int)
		do.DropletID = attrs["dropletId"].(int)

		// without any checks, expect an error if you are lying
		ch <- internal.MachineStatus{
			Status:       "available",
			IsImageReady: true,
		}

	} else if do.Project.Config.Image.From != "" && false {

		// Create image from url
		image, e = do.CreateCustomImage(&godo.CustomImageCreateRequest{
			Url:          do.Project.Config.Image.From,
			Name:         do.Project.Name,
			Region:       do.Project.Config.Image.Region,
			Distribution: "Unknown",
			Description:  "This image created by apker",
		})

		if e != nil {

			ch <- internal.MachineStatus{
				Error: e,
			}
			return
		}

		do.ImageID = image.ID

		// Wait for image be ready in DO
		for {

			// check every 12s, normaly this takes a while
			time.Sleep(12 * time.Second)

			image, _, e = do.DoClient.Images.GetByID(context.TODO(), do.ImageID)

			ch <- internal.MachineStatus{
				Status:       image.Status,
				IsImageReady: image.Status == "available",
				Error:        e,
			}

			if e != nil {
				return
			} else if image.Status == "available" {
				break
			}
		}

		dropletImage.ID = do.ImageID

	} else {

		dropletImage.Slug = do.Project.Config.Image.From
	}

	// Create new droplet when id not provided
	if attrs["dropletId"].(int) == 0 {

		// Create droplet for the current image
		droplet, e = do.CreateDroplet(dropletImage)

		// Droplet create error!
		if e != nil {

			ch <- internal.MachineStatus{
				Error: e,
			}
			return
		}

		do.DropletID = droplet.ID
	}

	// Wait for droplet to be ready
	for {

		if attrs["dropletId"].(int) == 0 {
			time.Sleep(5 * time.Second)
		}

		droplet, _, e = do.DoClient.Droplets.Get(context.TODO(), do.DropletID)

		// First check for error
		if e != nil {

			ch <- internal.MachineStatus{
				IsImageInstalled: true,
				IsMachineReady:   false,
				Status:           "Error",
				Error:            e,
			}
			break
		}

		if len(droplet.Networks.V4) > 0 {

			do.Project.AddrV4, _ = droplet.PublicIPv4()
		}

		ch <- internal.MachineStatus{
			Status:           droplet.Status,
			IsImageInstalled: true,
			IsMachineReady:   droplet.Status == "active",
			Error:            e,
		}

		if e != nil {
			return
		} else if droplet.Status == "active" {
			break
		}

		time.Sleep(5 * time.Second)
	}
}

func (do Digitalocean) CreateCustomImage(ImageRequest *godo.CustomImageCreateRequest) (*godo.Image, error) {

	image, _, err := do.DoClient.Images.Create(context.TODO(), ImageRequest)

	if err != nil {
		return &godo.Image{}, err
	}

	return image, nil
}

func (do *Digitalocean) CreateDroplet(image godo.DropletCreateImage) (*godo.Droplet, error) {

	dropletRequest := &godo.DropletCreateRequest{
		Name:   do.Project.Name,
		Region: do.Project.Config.Image.Region,
		Size:   do.Project.Config.Image.Size,
		Image:  image,
		Tags:   []string{"apker", "api"},
	}

	// TODO: Check if key fingerprint not exists, if not add new key to DO (or not! it should be already exists)
	if do.Project.PublicKey.Fingerprint != "" {
		dropletRequest.SSHKeys = []godo.DropletCreateSSHKey{
			godo.DropletCreateSSHKey{
				Fingerprint: do.Project.PublicKey.Fingerprint,
			},
		}
	}

	droplet, _, err := do.DoClient.Droplets.Create(context.TODO(), dropletRequest)

	return droplet, err
}

func NewDigitalocean(p *internal.Project) (*Digitalocean, error) {

	if _, ok := p.Config.Provider.Credentials["API_KEY"]; !ok {

		return &Digitalocean{}, errors.New("API_KEY is required!")
	}

	tokenSource := &TokenSource{
		AccessToken: p.Config.Provider.Credentials["API_KEY"],
	}

	oauth := oauth2.NewClient(context.Background(), tokenSource)

	do := &Digitalocean{
		Project:  p,
		Oauth:    oauth,
		DoClient: godo.NewClient(oauth),
	}

	return do, nil
}
