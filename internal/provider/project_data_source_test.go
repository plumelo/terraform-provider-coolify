package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectDataSource(t *testing.T) {
	resName := "data.coolify_project.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "coolify_project" "test" {
					uuid = "` + testAccProjectUUID + `"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "id", "38"),
					resource.TestCheckResourceAttr(resName, "uuid", testAccProjectUUID),
					resource.TestCheckResourceAttrSet(resName, "name"),
					resource.TestCheckResourceAttrSet(resName, "environments.#"),
				),
			},
		},
	})
}
