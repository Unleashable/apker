// Copyright 2020 Mohammed El Bahja. All rights reserved.
// Use of this source code is governed by a MIT license.

package internal

type Attributes map[string]interface{}

type MachineStatus struct {
	IsImageReady     bool
	IsImageInstalled bool
	IsMachineReady   bool
	Status           string
	Error            error
}

type Machine struct {
	ID     int
	Addr   string
	Name   string
	Region string
	Status string
}

type Provider interface {

	// Setup virtual machine on cloud provider
	SetupMachine(chan MachineStatus, Attributes)
}
