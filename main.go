package main

import (
	"context"
	"flag"
	"log"

	"github.com/aztfmod/opentofu-provider-restful/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Generate the provider document.
//go:generate go tool tfplugindocs generate

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	ctx := context.Background()
	serveOpts := providerserver.ServeOpts{
		Debug:   debug,
		Address: "registry.opentofu.org/aztfmod/restful",
	}

	err := providerserver.Serve(ctx, provider.New, serveOpts)

	if err != nil {
		log.Fatalf("Error serving provider: %s", err)
	}
}
