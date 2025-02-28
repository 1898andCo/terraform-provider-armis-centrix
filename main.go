// Copyright (c) 1898 & Co.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"flag"
	"log"

	"your-module-path/provider" // Replace with actual path

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Dynamically injected via Goreleaser.
var version string = "0.1.0"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/provider",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
