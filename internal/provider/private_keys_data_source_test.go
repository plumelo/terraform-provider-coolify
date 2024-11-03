package provider_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-coolify/internal/provider"
)

func TestAccPrivateKeysDataSource(t *testing.T) {
	resName := "data.coolify_private_keys.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Without filters
			{
				Config: `data "coolify_private_keys" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "private_keys.#", "6"),
					// Check the first private key in the list
					resource.TestCheckResourceAttr(resName, "private_keys.0.id", "0"),
					resource.TestCheckResourceAttr(resName, "private_keys.0.name", "localhost's key"),
					resource.TestCheckResourceAttrSet(resName, "private_keys.0.private_key"),
					resource.TestCheckResourceAttrSet(resName, "private_keys.0.uuid"),
				),
			},
			// Single filter by name
			{
				Config: `
				data "coolify_private_keys" "test" {
					filter {
						name = "name"
						values = ["tf-acc-test-1", "tf-acc-test-2"]
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "private_keys.#", "2"),
					resource.TestCheckResourceAttr(resName, "private_keys.0.name", "tf-acc-test-1"),
					resource.TestCheckResourceAttr(resName, "private_keys.1.name", "tf-acc-test-2"),
				),
			},
			// Multiple filters
			{
				Config: `
				data "coolify_private_keys" "test" {
					filter {
						name = "name"
						values = ["tf-acc-test-1"]
					}
					filter {
						name = "description"
						values = ["Manually created for Datasource Acceptance Test"]
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "private_keys.#", "1"),
					resource.TestCheckResourceAttr(resName, "private_keys.0.name", "tf-acc-test-1"),
					resource.TestCheckResourceAttr(resName, "private_keys.0.description", "Manually created for Datasource Acceptance Test"),
				),
			},
			// Filter by non-string fields
			{
				Config: `
				data "coolify_private_keys" "test" {
					filter {
						name = "team_id"
						values = ["0"]
					}
					filter {
						name = "is_git_related"
						values = ["false"]
					}
					filter {
						name = "name"
						values = ["tf-acc-test-2"]
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "private_keys.#", "1"),
					resource.TestCheckResourceAttr(resName, "private_keys.0.name", "tf-acc-test-2"),
					resource.TestCheckResourceAttr(resName, "private_keys.0.team_id", "0"),
					resource.TestCheckResourceAttr(resName, "private_keys.0.is_git_related", "false"),
				),
			},
		},
	})
}

func TestPrivateKeysDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	ds := provider.NewPrivateKeysDataSource()
	resp := &datasource.SchemaResponse{}
	ds.Schema(ctx, datasource.SchemaRequest{}, resp)

	// Test private_key sensitivity
	privateKeysAttr := resp.Schema.Attributes["private_keys"].(schema.SetNestedAttribute)
	privateKeyAttr := privateKeysAttr.NestedObject.Attributes["private_key"].(schema.StringAttribute)
	if !privateKeyAttr.Sensitive {
		t.Error("private_key field should be marked as sensitive in schema")
	}

	// Test filter block
	_, ok := resp.Schema.Blocks["filter"].(schema.ListNestedBlock)
	if !ok {
		t.Error("filter should be a ListNestedBlock")
	}
}
