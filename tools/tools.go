//go:build tools
// +build tools

package main

import (
	// Terraform plugin framework code generation
	_ "github.com/hashicorp/terraform-plugin-codegen-framework/cmd/tfplugingen-framework"
	// Terraform plugin code generation from openapi spec
	_ "github.com/hashicorp/terraform-plugin-codegen-openapi/cmd/tfplugingen-openapi"
	// Documentation generation
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
	// OpenAPI client generation
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
)
