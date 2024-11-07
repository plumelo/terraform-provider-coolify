package provider_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-coolify/internal/provider"
)

func TestAccServersDataSource(t *testing.T) {
	resName := "data.coolify_servers.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Without filters
			{
				Config: `data "coolify_servers" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "servers.#", "2"),
					// Check the last server in the list (expecting the first created server, order seems to be id descending)
					resource.TestCheckResourceAttr(resName, "servers.1.id", "0"),
					resource.TestCheckResourceAttr(resName, "servers.1.ip", "host.docker.internal"),
					resource.TestCheckResourceAttr(resName, "servers.1.settings.id", "1"),
				),
			},
			// Single filter by uuid
			{
				Config: `
				data "coolify_servers" "test" {
					filter {
						name = "uuid"
						values = ["rg8ks8c"]
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "servers.#", "1"),
					resource.TestCheckResourceAttr(resName, "servers.0.id", "0"),
					resource.TestCheckResourceAttr(resName, "servers.0.uuid", "rg8ks8c"),
					resource.TestCheckResourceAttr(resName, "servers.0.ip", "host.docker.internal"),
				),
			},
		},
	})
}

func TestServersDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	ds := provider.NewServersDataSource()
	resp := &datasource.SchemaResponse{}
	ds.Schema(ctx, datasource.SchemaRequest{}, resp)

	// Test filter block
	_, ok := resp.Schema.Blocks["filter"].(schema.ListNestedBlock)
	if !ok {
		t.Error("filter should be a ListNestedBlock")
	}
}
