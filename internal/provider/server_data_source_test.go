package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServerDataSource(t *testing.T) {
	resName := "data.coolify_server.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// With ID
			{
				Config: `data "coolify_server" "test" {
					uuid = "rg8ks8c"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "id", "0"),
					resource.TestCheckResourceAttr(resName, "uuid", "rg8ks8c"),
					resource.TestCheckResourceAttr(resName, "ip", "host.docker.internal"),
					resource.TestCheckResourceAttr(resName, "port", "22"),
					resource.TestCheckResourceAttr(resName, "user", "root"),
					resource.TestCheckResourceAttr(resName, "settings.server_id", "0"),
					resource.TestCheckResourceAttr(resName, "settings.id", "1"),
				),
			},
		},
	})
}
