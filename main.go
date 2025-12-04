// Copyright (c) Adfinis
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/adfinis/terraform-provider-bastion/bastion"
	"github.com/adfinis/terraform-provider-bastion/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	version string = "dev"
)

func main() {

	c, err := bastion.New(&bastion.Config{
		Host:                  "localhost",
		Port:                  2222,
		Username:              "bastionadmin",
		StrictHostKeyChecking: false,
	}, bastion.WithPrivateKeyFileAuth("ssh-keys/id_ed25519"))
	if err != nil {
		log.Fatalf("failed to create bastion client: %v", err)
	}

	fmt.Println(c.AccountInfo("bastionadmin"))
	os.Exit(0)

	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		// TODO: Update this string with the published name of your provider.
		// Also update the tfplugindocs generate command to either remove the
		// -provider-name flag or set its value to the updated provider name.
		Address: "registry.terraform.io/adfinis/bastion",
		Debug:   debug,
	}

	err = providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
