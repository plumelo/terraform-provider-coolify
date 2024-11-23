package provider_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-coolify/internal/provider"
)

func TestAccPrivatekeyDataSource(t *testing.T) {
	resName := "data.coolify_private_key.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Without filters
			{
				Config: `data "coolify_private_key" "test" {
					uuid = "` + testAccPrivateKeyUUID + `"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "id", "0"),
					resource.TestCheckResourceAttr(resName, "name", "localhost's key"),
					resource.TestCheckResourceAttrSet(resName, "private_key"),
					resource.TestCheckResourceAttr(resName, "uuid", testAccPrivateKeyUUID),
				),
			},
		},
	})
}

func TestPrivatekeyDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	ds := provider.NewPrivateKeyDataSource()
	resp := &datasource.SchemaResponse{}
	ds.Schema(ctx, datasource.SchemaRequest{}, resp)

	// Test private_key sensitivity
	privateKeyAttr := resp.Schema.Attributes["private_key"].(schema.StringAttribute)
	if !privateKeyAttr.Sensitive {
		t.Error("private_key field should be marked as sensitive in schema")
	}
}
