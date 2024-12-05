package service_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-coolify/internal/acctest"
)

func TestAccServerDomainsDataSource(t *testing.T) {
	resName := "data.coolify_server_domains.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "coolify_server_domains" "test" {
					uuid  = "` + acctest.ServerUUID + `"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "uuid", acctest.ServerUUID),
					resource.TestCheckResourceAttrSet(resName, "server_domains.#"),
					resource.TestCheckResourceAttrSet(resName, "server_domains.0.ip"),
					resource.TestCheckResourceAttrSet(resName, "server_domains.0.domains.#"),
				),
			},
		},
	})
}
