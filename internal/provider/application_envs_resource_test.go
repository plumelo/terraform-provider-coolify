package provider_test

import (
	"context"
	"fmt"
	"testing"

	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"terraform-provider-coolify/internal/provider"
)

func TestAccApplicationEnvsResource(t *testing.T) {
	resName := "coolify_application_envs.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{ // Create and Read testing
				Config: `
					resource "coolify_application_envs" "test" {
						uuid = "mc8gw00wscww4gskgk0gwgw0"
						env {
							key        = "key1"
							value      = "value1"
						}
						env {
							key        = "key1"
							value      = "value1-preview"
							is_preview = true
						}
						env {
							key        = "key2"
							value      = "value2"
						}
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "env.#"),
					resource.TestCheckResourceAttr(resName, "env.0.key", "key1"),
					resource.TestCheckResourceAttr(resName, "env.0.value", "value1"),
				),
			},
			{ // ImportState testing
				ResourceName:                         resName,
				ImportState:                          true,
				ImportStateVerifyIdentifierAttribute: "uuid",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					return s.RootModule().Resources[resName].Primary.Attributes["uuid"], nil
				},
				ImportStateCheck: func(is []*terraform.InstanceState) error {
					if len(is) != 1 {
						return fmt.Errorf("expected 1 state, got %d", len(is))
					}
					// Expect 1 additional pre-existing env
					if is[0].Attributes["env.#"] != "4" {
						return fmt.Errorf("expected 4 envs, got %s", is[0].Attributes["env.#"])
					}
					if is[0].Attributes["uuid"] != "mc8gw00wscww4gskgk0gwgw0" {
						return fmt.Errorf("expected uuid to be mc8gw00wscww4gskgk0gwgw0, got %s", is[0].Attributes["uuid"])
					}
					return nil
				},
			},
			{ // Update and Read testing
				Config: `
					resource "coolify_application_envs" "test" {
						uuid = "mc8gw00wscww4gskgk0gwgw0"
						env {
							key        = "key1-1"
							value      = "value1-1"
						}
						env {
							key        = "key2"
							value      = "value2-2"
						}
					}`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resName, plancheck.ResourceActionUpdate),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resName, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "env.#"),
					resource.TestCheckResourceAttr(resName, "env.0.key", "key1-1"),
					resource.TestCheckResourceAttr(resName, "env.0.value", "value1-1"),
					resource.TestCheckResourceAttr(resName, "env.1.key", "key2"),
					resource.TestCheckResourceAttr(resName, "env.1.value", "value2-2"),
				),
			},
		},
	})
}

func TestApplicationEnvsResourceSchema(t *testing.T) {
	ctx := context.Background()
	rs := provider.NewApplicationEnvsResource()
	resp := &tfresource.SchemaResponse{}
	rs.Schema(ctx, tfresource.SchemaRequest{}, resp)
}
