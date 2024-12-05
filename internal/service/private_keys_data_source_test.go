package service_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-coolify/internal/acctest"
)

func TestAccPrivateKeysDataSource(t *testing.T) {
	resName := "data.coolify_private_keys.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Without filters
			{
				Config: `data "coolify_private_keys" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "private_keys.#", "6"),
					// Check the first private key in the list
					resource.TestCheckResourceAttr(resName, "private_keys.0.id", "0"),
					resource.TestCheckResourceAttr(resName, "private_keys.0.name", "localhost's key"),
					resource.TestCheckResourceAttrSet(resName, "private_keys.0.private_key"),
					resource.TestCheckResourceAttrSet(resName, "private_keys.0.uuid"),
				),
			},
			// Single filter by name
			{
				Config: `
				data "coolify_private_keys" "test" {
					filter {
						name = "name"
						values = ["tf-acc-test-1", "tf-acc-test-2"]
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "private_keys.#", "2"),
					resource.TestCheckResourceAttr(resName, "private_keys.0.name", "tf-acc-test-1"),
					resource.TestCheckResourceAttr(resName, "private_keys.1.name", "tf-acc-test-2"),
				),
			},
			// Multiple filters
			{
				Config: `
				data "coolify_private_keys" "test" {
					filter {
						name = "name"
						values = ["tf-acc-test-1"]
					}
					filter {
						name = "description"
						values = ["Manually created for Datasource Acceptance Test"]
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "private_keys.#", "1"),
					resource.TestCheckResourceAttr(resName, "private_keys.0.name", "tf-acc-test-1"),
					resource.TestCheckResourceAttr(resName, "private_keys.0.description", "Manually created for Datasource Acceptance Test"),
				),
			},
			// Filter by non-string fields
			{
				Config: `
				data "coolify_private_keys" "test" {
					filter {
						name = "team_id"
						values = ["0"]
					}
					filter {
						name = "is_git_related"
						values = ["false"]
					}
					filter {
						name = "name"
						values = ["tf-acc-test-2"]
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "private_keys.#", "1"),
					resource.TestCheckResourceAttr(resName, "private_keys.0.name", "tf-acc-test-2"),
					resource.TestCheckResourceAttr(resName, "private_keys.0.team_id", "0"),
					resource.TestCheckResourceAttr(resName, "private_keys.0.is_git_related", "false"),
				),
			},
		},
	})
}
