package main

// * Generate the provider specification.json from the OpenAPI specification
//go:generate go run github.com/hashicorp/terraform-plugin-codegen-openapi/cmd/tfplugingen-openapi generate --config config/tfplugingen-config.yml --output specification.json config/coolify_openapi.yml

// * Generate Go resource, datasource, provider schema from specification.json
//go:generate go run github.com/hashicorp/terraform-plugin-codegen-framework/cmd/tfplugingen-framework generate resources --input specification.json --output internal/provider/generated
//go:generate go run github.com/hashicorp/terraform-plugin-codegen-framework/cmd/tfplugingen-framework generate data-sources --input specification.json --output internal/provider/generated

// * Generate the OpenAPI client
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config config/oapi-codegen-config.yml config/coolify_openapi.yml

// * If you do not have terraform installed, you can remove the formatting command, but its suggested to ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// * Run the docs generation tool, check its repository for more information on how it works and how docs can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name coolify

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"terraform-provider-coolify/internal/provider"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"

	// goreleaser can pass other information to the main package, such as the specific commit
	// https://goreleaser.com/cookbooks/using-main.version/
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "hashicorp.com/edu/coolify",
		// 		Address: "registry.terraform.io/sierrajc/coolify",
		Debug: debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
