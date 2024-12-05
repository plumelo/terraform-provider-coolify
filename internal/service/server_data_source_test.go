package service_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-coolify/internal/acctest"
)

func TestAccServerDataSource(t *testing.T) {
	resName := "data.coolify_server.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// With ID
			{
				Config: `data "coolify_server" "test" {
					uuid  = "` + acctest.ServerUUID + `"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "id", "0"),
					resource.TestCheckResourceAttr(resName, "uuid", acctest.ServerUUID),
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
