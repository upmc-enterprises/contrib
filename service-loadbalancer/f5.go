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

// F5Config holds configuration for the f5 plugin.
type F5Config struct {
	// Host specifies the hostname or IP address of the F5 BIG-IP host.
	Host string

	// Username specifies the username with the plugin should authenticate
	// with the F5 BIG-IP host.
	Username string

	// Password specifies the password with which the plugin should
	// authenticate with F5 BIG-IP.
	Password string

	// Insecure specifies whether the F5 plugin should perform strict certificate
	// validation for connections to the F5 BIG-IP host.
	Insecure bool

	// PartitionPath specifies the F5 partition path to use. This is used
	// to create an access control boundary for users and applications.
	PartitionPath string
}

type f5Controller struct {
	Config F5Config
}

func newF5Controller(host, user, password, partition string, insecure bool) *f5Controller {
	ctrl := f5Controller{
		Config: F5Config{
			Host:          host,
			Username:      user,
			Password:      password,
			Insecure:      insecure,
			PartitionPath: partition,
		},
	}

	return &ctrl
}

func (ctrl *f5Controller) createService() error {
	return nil
}
