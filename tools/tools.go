//go:build generate

package tools

import (
	_ "github.com/hashicorp/terraform-plugin-codegen-framework/cmd/tfplugingen-framework"
	_ "github.com/hashicorp/terraform-plugin-codegen-openapi/cmd/tfplugingen-openapi"
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
)

// Codegen flow: Coolify OpenAPI Spec -> TF Specification -> Go code -> OpenAPI Client -> Docs

// * Generate the provider spec file from the OpenAPI specification
//go:generate go run github.com/hashicorp/terraform-plugin-codegen-openapi/cmd/tfplugingen-openapi generate --config tfplugingen-openapi.yml --output spec_gen.json openapi.yml

// * Generate Go resource, data source, provider schema from spec file
//go:generate go run github.com/hashicorp/terraform-plugin-codegen-framework/cmd/tfplugingen-framework generate resources --input spec_gen.json --output ../internal/provider/generated
//go:generate go run github.com/hashicorp/terraform-plugin-codegen-framework/cmd/tfplugingen-framework generate data-sources --input spec_gen.json --output ../internal/provider/generated
// ? Scaffold code with: tfplugingen-framework scaffold data-source|resource --name <name> --output-dir ./internal/provider

// * Generate the OpenAPI client
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config oapi-codegen.yml openapi.yml

// * If you do not have terraform installed, you can remove the formatting command, but its suggested to ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ../examples/

// * Run the docs generation tool, check its repository for more information on how it works and how docs can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name coolify --provider-dir ../
