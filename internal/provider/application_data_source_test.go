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
					uuid = "` + testAccApplicationUUID + `"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "uuid", testAccApplicationUUID),
					resource.TestCheckResourceAttr(resName, "build_pack", "dockerfile"),
					resource.TestCheckResourceAttr(resName, "created_at", "2024-11-10T08:59:09Z"),
					resource.TestCheckResourceAttr(resName, "name", "dockerfile-"+testAccApplicationUUID),
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
