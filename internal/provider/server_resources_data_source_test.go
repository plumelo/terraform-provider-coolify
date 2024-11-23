package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServerDomainsDataSource(t *testing.T) {
	resName := "data.coolify_server_domains.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "coolify_server_domains" "test" {
					uuid  = "` + testAccServerUUID + `"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "uuid", testAccServerUUID),
					resource.TestCheckResourceAttrSet(resName, "server_domains.#"),
					resource.TestCheckResourceAttrSet(resName, "server_domains.0.ip"),
					resource.TestCheckResourceAttrSet(resName, "server_domains.0.domains.#"),
				),
			},
		},
	})
}
