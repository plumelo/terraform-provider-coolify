package provider_test

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccPrivateKeyResource(t *testing.T) {
	resName := "coolify_private_key.test"

	privateKey := generatePrivateKey(t)
	privateKeyUpdated := generatePrivateKey(t)

	// Escape newlines for inclusion in HCL configuration
	hclPrivateKey := strings.ReplaceAll(privateKey, "\n", "\\n")
	hclPrivateKeyUpdated := strings.ReplaceAll(privateKeyUpdated, "\n", "\\n")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{ // Create and Read testing
				Config: `
				resource "coolify_private_key" "test" {
					name        = "TerraformAccTest"
					description = "Terraform acceptance testing"
					private_key = "` + hclPrivateKey + `"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", "TerraformAccTest"),
					resource.TestCheckResourceAttr(resName, "description", "Terraform acceptance testing"),
					resource.TestCheckResourceAttr(resName, "private_key", privateKey),
					// Verify dynamic values
					resource.TestCheckResourceAttrSet(resName, "uuid"),
					resource.TestCheckResourceAttrSet(resName, "id"),
					resource.TestCheckResourceAttrSet(resName, "created_at"),
					resource.TestCheckResourceAttrSet(resName, "updated_at"),
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
				resource "coolify_private_key" "test" {
					name        = "TerraformAccTestUpdated"
					description = "Terraform acceptance testing"
					private_key = "` + hclPrivateKeyUpdated + `"
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
					resource.TestCheckResourceAttr(resName, "private_key", privateKeyUpdated),
				),
			},
		},
	})
}

func generatePrivateKey(t *testing.T) string {
	t.Helper()
	// Generate an Ed25519 private key
	_, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	// Marshal the private key into PKCS8 format
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("failed to marshal private key: %v", err)
	}

	// Encode the private key into PEM format
	privateKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	return string(privateKeyPem)
}
