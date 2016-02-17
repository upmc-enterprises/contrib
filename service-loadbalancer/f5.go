/*
Copyright 2015 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"github.com/golang/glog"
	f5 "github.com/openshift/origin/plugins/router/f5"
)

// const (
// 	host          = "https://ftrf5devintltm01.isd.upmc.edu"
// 	username      = "slokas_adm"
// 	password      = "foo"
// 	insecure      = true
// 	partitionPath = "tdc"
// )

type f5Controller struct {
	F5Client *f5.F5Plugin
}

func newF5Controller(host, username, password, partitionPath string, insecure bool) *f5Controller {
	config := f5.F5PluginConfig{
		Host:          host,
		Username:      username,
		Password:      password,
		Insecure:      insecure,
		PartitionPath: partitionPath,
	}
	client, err := f5.NewF5Plugin(config)

	if err != nil {
		glog.Fatalf("Could not init f5Controller: %v", err)
	}

	controller := f5Controller{
		F5Client: client,
	}

	client.ensurePoolExists("foo")

	return &controller
}

// func (f5ctl *f5Controller) checkPoolExists(pool string) error {
// 	return f5ctl.F5Client.ensurePoolExists(pool)
// }
