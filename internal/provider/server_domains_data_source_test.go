package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServerResourcesDataSource(t *testing.T) {
	resName := "data.coolify_server_resources.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "coolify_server_resources" "test" {
					uuid  = "` + testAccServerUUID + `"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "uuid", testAccServerUUID),
					resource.TestCheckResourceAttrSet(resName, "server_resources.#"),
					resource.TestCheckResourceAttrSet(resName, "server_resources.0.id"),
					resource.TestCheckResourceAttrSet(resName, "server_resources.0.name"),
				),
			},
		},
	})
}
