package provider_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccPrivateKeyResource(t *testing.T) {
	resName := "coolify_private_key.test"
	randomName := getRandomResourceName("pk")

	_, privateKey1, err := acctest.RandSSHKeyPair(t.Name())
	if err != nil {
		t.Fatal(err)
	}
	_, privateKey2, err := acctest.RandSSHKeyPair(t.Name())
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{ // Create and Read testing
				Config: testAccPrivateKeyResourceConfig(randomName, privateKey1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", randomName),
					resource.TestCheckResourceAttr(resName, "private_key", privateKey1),
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
				Config: testAccPrivateKeyResourceConfig(randomName, privateKey2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resName, plancheck.ResourceActionUpdate),
						plancheck.ExpectUnknownValue(resName, tfjsonpath.New("fingerprint")),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resName, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "uuid"),
					resource.TestCheckResourceAttr(resName, "name", randomName),
					resource.TestCheckResourceAttr(resName, "private_key", privateKey2),
				),
			},
		},
	})
}

func testAccPrivateKeyResourceConfig(name, privateKey string) string {
	return fmt.Sprintf(`
		resource "coolify_private_key" "test" {
			name        = "%s"
			description = "Terraform acceptance testing"
			private_key = "%s"
		}
	`,
		name,
		strings.ReplaceAll(privateKey, "\n", "\\n"),
	)
}
