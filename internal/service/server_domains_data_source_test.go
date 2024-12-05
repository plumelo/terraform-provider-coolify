package service_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-coolify/internal/acctest"
)

func TestAccServerResourcesDataSource(t *testing.T) {
	resName := "data.coolify_server_resources.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "coolify_server_resources" "test" {
					uuid  = "` + acctest.ServerUUID + `"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "uuid", acctest.ServerUUID),
					resource.TestCheckResourceAttrSet(resName, "server_resources.#"),
					resource.TestCheckResourceAttrSet(resName, "server_resources.0.id"),
					resource.TestCheckResourceAttrSet(resName, "server_resources.0.name"),
				),
			},
		},
	})
}
