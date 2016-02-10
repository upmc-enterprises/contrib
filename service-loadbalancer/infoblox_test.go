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
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	user     = "username"
	password = "password"
	url      = "http://something/bar"
)

func setupMockService(testType string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch testType {
		case "good":
			var hosts []infoBloxHost
			hosts = GetRecordData()
			data, _ := json.Marshal(hosts)

			w.Header().Add("Content-Type", "application/json")
			w.Write(data)
		case "500":
			w.WriteHeader(http.StatusInternalServerError)
		case "badmodel":
			w.Header().Add("Content-Type", "application/json")
			data, _ := json.Marshal("{foo: 'bar'}")
			w.Write(data)
		}
	}))
}

func TestNewObject(t *testing.T) {
	ibc := newInfobloxController(user, password, url)

	if ibc.user != user {
		t.Fatalf("Expected Username: %+v, got %+v", user, ibc.user)
	}

	if ibc.password != password {
		t.Fatalf("Expected Password: %+v, got %+v", password, ibc.password)
	}

	if ibc.baseEndpoint != url {
		t.Fatalf("Expected Url: %+v, got %+v", url, ibc.baseEndpoint)
	}
}

func TestGetHostParseData(t *testing.T) {
	ms := setupMockService("good")
	defer ms.Close()

	ibc := newInfobloxController(user, password, ms.URL)

	fmt.Println("server.url: ", ms.URL)

	host, err := ibc.getHost("foo")

	if err != nil {
		t.Error(err)
	}

	if len(host) != 2 {
		t.Fatalf("Expected: %+v, got %+v", "2 host", len(host))
	}
}

func TestGetHost_ServerError(t *testing.T) {
	ms := setupMockService("500")
	defer ms.Close()

	ibc := newInfobloxController(user, password, ms.URL)

	fmt.Println("server.url: ", ms.URL)

	host, err := ibc.getHost("foo")

	if err != nil {
		t.Error(err)
	}

	if len(host) != 0 {
		t.Fatalf("Expected: %+v, got %+v", "0 hosts", len(host))
	}
}

func TestGetHost_BadModel(t *testing.T) {
	ms := setupMockService("badmodel")
	defer ms.Close()

	ibc := newInfobloxController(user, password, ms.URL)

	fmt.Println("server.url: ", ms.URL)

	_, err := ibc.getHost("foo")

	if err == nil {
		t.Error("Expected 'err' not to be nil!")
	}
}

// Sample Mock Data
func GetRecordData() []infoBloxHost {
	result := []infoBloxHost{
		{
			Ref: "host1",
		},
		{
			Ref: "host2",
		},
	}
	return result
}
