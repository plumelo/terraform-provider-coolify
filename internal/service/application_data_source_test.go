package service_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-coolify/internal/acctest"
	"terraform-provider-coolify/internal/service"
)

func TestAccApplicationDataSource(t *testing.T) {
	resName := "data.coolify_application.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "coolify_application" "test" {
					uuid = "` + acctest.ApplicationUUID + `"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "uuid", acctest.ApplicationUUID),
					resource.TestCheckResourceAttr(resName, "build_pack", "dockerfile"),
					resource.TestCheckResourceAttr(resName, "created_at", "2024-11-10T08:59:09Z"),
					resource.TestCheckResourceAttr(resName, "name", "dockerfile-"+acctest.ApplicationUUID),
					resource.TestCheckNoResourceAttr(resName, "description"),
				),
			},
		},
	})
}

func TestApplicationDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	ds := service.NewApplicationDataSource()
	resp := &datasource.SchemaResponse{}
	ds.Schema(ctx, datasource.SchemaRequest{}, resp)
}
