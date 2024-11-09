//go:build generate

package tools

import (
	_ "github.com/hashicorp/terraform-plugin-codegen-framework/cmd/tfplugingen-framework"
	_ "github.com/hashicorp/terraform-plugin-codegen-openapi/cmd/tfplugingen-openapi"
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
	_ "github.com/speakeasy-api/openapi-overlay"
)

// Codegen flow: Coolify OpenAPI Spec -> TF Specification -> Go code -> OpenAPI Client -> Docs

// * Overlay the OpenAPI file with fixes and customizations
//go:generate echo "Apply OpenAPI overlay..."
//go:generate sh -c "go run github.com/speakeasy-api/openapi-overlay apply overlay.yml openapi.yml > openapi-overlay.yml"

// * Generate the provider spec file from the OpenAPI specification
//go:generate echo "Generate provider spec file..."
//go:generate go run github.com/hashicorp/terraform-plugin-codegen-openapi/cmd/tfplugingen-openapi generate --config tfplugingen-openapi.yml --output spec_gen.json openapi-overlay.yml

// * Generate Go resource, data source, provider schema from spec file
//go:generate echo "Generate Go code from provider spec file..."
//go:generate go run github.com/hashicorp/terraform-plugin-codegen-framework/cmd/tfplugingen-framework generate resources --input spec_gen.json --output ../internal/provider/generated
//go:generate go run github.com/hashicorp/terraform-plugin-codegen-framework/cmd/tfplugingen-framework generate data-sources --input spec_gen.json --output ../internal/provider/generated
// ? Scaffold code with: tfplugingen-framework scaffold data-source|resource --name <name> --output-dir ./internal/provider

// * Generate the OpenAPI client
//go:generate echo "Generate OpenAPI client..."
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config oapi-codegen.yml openapi-overlay.yml

// * If you do not have terraform installed, you can remove the formatting command, but its suggested to ensure the documentation is formatted properly.
//go:generate echo "Format the generated code..."
//go:generate terraform fmt -recursive ../examples/

// * Run the docs generation tool, check its repository for more information on how it works and how docs can be customized.
//go:generate echo "Generate provider documentation..."
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name coolify --provider-dir ../
