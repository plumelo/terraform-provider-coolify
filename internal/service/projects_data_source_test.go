package service_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-coolify/internal/acctest"
	"terraform-provider-coolify/internal/service"
)

func TestAccProjectsDataSource(t *testing.T) {
	resName := "data.coolify_projects.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Without filters
			{
				Config: `data "coolify_projects" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "projects.#"),
					// Check the last server in the list (expecting the first created server, order seems to be id descending)
					resource.TestCheckResourceAttrSet(resName, "projects.0.id"),
					resource.TestCheckResourceAttrSet(resName, "projects.0.uuid"),
					resource.TestCheckResourceAttrSet(resName, "projects.0.name"),
					resource.TestCheckNoResourceAttr(resName, "projects.0.environments"),
				),
			},
			// Single filter by name
			{
				Config: `
				data "coolify_projects" "test" {
					filter {
						name = "name"
						values = ["AccTestProj"]
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "projects.#", "1"),
					resource.TestCheckResourceAttr(resName, "projects.0.id", "38"),
					resource.TestCheckResourceAttr(resName, "projects.0.uuid", acctest.ProjectUUID),
					resource.TestCheckResourceAttr(resName, "projects.0.name", "AccTestProj"),
					resource.TestCheckNoResourceAttr(resName, "projects.0.environments"),
				),
			},
		},
	})
}

func TestProjectsDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	ds := service.NewProjectsDataSource()
	resp := &datasource.SchemaResponse{}
	ds.Schema(ctx, datasource.SchemaRequest{}, resp)

	// Test filter block
	_, ok := resp.Schema.Blocks["filter"].(schema.ListNestedBlock)
	if !ok {
		t.Error("filter should be a ListNestedBlock")
	}
}
