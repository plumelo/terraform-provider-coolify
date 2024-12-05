package service_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-coolify/internal/acctest"
)

func TestAccServiceDataSource(t *testing.T) {
	resName := "data.coolify_service.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "coolify_service" "test" {
					uuid = "` + acctest.ServiceUUID + `"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "uuid", acctest.ServiceUUID),
					resource.TestCheckResourceAttr(resName, "created_at", "2024-11-18T06:03:18.000000Z"),
					resource.TestCheckResourceAttrSet(resName, "docker_compose"),
					resource.TestCheckResourceAttr(resName, "name", "service-"+acctest.ServiceUUID),
					resource.TestCheckNoResourceAttr(resName, "description"),
				),
			},
		},
	})
}
