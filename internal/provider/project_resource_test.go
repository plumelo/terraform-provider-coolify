package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccProjectResource(t *testing.T) {
	resName := "coolify_project.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{ // Create and Read testing
				Config: `
				resource "coolify_project" "test" {
					name        = "TerraformAccTest"
					description = "Terraform acceptance testing"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", "TerraformAccTest"),
					resource.TestCheckResourceAttr(resName, "description", "Terraform acceptance testing"),
					// Verify dynamic values
					resource.TestCheckResourceAttrSet(resName, "uuid"),
					resource.TestCheckResourceAttrSet(resName, "id"),
				),
			},
			{ // ImportState testing
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					return s.RootModule().Resources[resName].Primary.Attributes["uuid"], nil
				},
			},
			{ // Update and Read testing
				Config: `
				resource "coolify_project" "test" {
					name        = "TerraformAccTestUpdated"
					description = "Terraform acceptance testing"
				}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resName, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(resName, tfjsonpath.New("name"), knownvalue.StringExact("TerraformAccTestUpdated")),
						plancheck.ExpectKnownValue(resName, tfjsonpath.New("description"), knownvalue.StringExact("Terraform acceptance testing")),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resName, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "uuid"),
					resource.TestCheckResourceAttr(resName, "name", "TerraformAccTestUpdated"),
					resource.TestCheckResourceAttr(resName, "description", "Terraform acceptance testing"),
				),
			},
		},
	})
}
