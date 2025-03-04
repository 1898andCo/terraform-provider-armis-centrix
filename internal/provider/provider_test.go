// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"os"
	"sync"
	"testing"

	"github.com/1898andCo/terraform-provider-armis-centrix/internal/armis"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"armis": providerserver.NewProtocol6WithError(New("armis")()),
}

func requireEnv(t *testing.T, name string) {
	t.Helper()
	if value := os.Getenv(name); value == "" {
		t.Fatalf("Missing required environment variable: %s", name)
	}
}

func testAccPreCheck(t *testing.T) {
	requireEnv(t, "ARMIS_API_KEY")
	requireEnv(t, "ARMIS_API_URL")
}

var (
	clientInstance *armis.Client
	clientOnce     sync.Once
)

func testClient(t *testing.T) *armis.Client {
	clientOnce.Do(func() {
		options := armis.ClientOptions{
			ApiUrl:     os.Getenv("ARMIS_API_URL"),
			ApiKey:     os.Getenv("ARMIS_API_KEY"),
			ApiVersion: "v1",
		}

		var err error
		clientInstance, err = armis.NewClient(options)
		if err != nil {
			t.Fatalf("Failed to initialize Armis client: %v", err)
		}
	})

	return clientInstance
}
