package provider_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-coolify/internal/provider"
)

func TestAccApplicationDataSource(t *testing.T) {
	resName := "data.coolify_application.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "coolify_application" "test" {
					uuid = "i8wcgsk"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "uuid", "i8wcgsk"),
					resource.TestCheckResourceAttr(resName, "build_pack", "dockerimage"),
					resource.TestCheckResourceAttr(resName, "created_at", "2024-04-30T09:49:56Z"),
					resource.TestCheckResourceAttr(resName, "name", "api"),
					resource.TestCheckNoResourceAttr(resName, "description"),
				),
			},
		},
	})
}

func TestApplicationDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	ds := provider.NewApplicationDataSource()
	resp := &datasource.SchemaResponse{}
	ds.Schema(ctx, datasource.SchemaRequest{}, resp)
}
