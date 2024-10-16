package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// func TestAccPrivateKeyResource(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			// Test Create and Read
// 			{
// 				Config: providerConfig + `
// 					resource "coolify_private_key" "testacc" {
// 						name        = "TerraformAccTest"
// 						description = "Terraform acceptance testing"
// 						private_key = ` + mockPrivateKey + `
// 					}
// 				`,
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttrSet("coolify_private_key.testacc", "uuid"),
// 					resource.TestCheckResourceAttr("coolify_private_key.testacc", "name", "TerraformAccTest"),
// 					resource.TestCheckResourceAttr("coolify_private_key.testacc", "description", "Terraform acceptance testinga"),
// 					resource.TestCheckResourceAttr("coolify_private_key.testacc", "private_key", mockPrivateKey),
// 				// resource.TestCheckResourceAttr("ubicloud_firewall.testacc", "project_id", GetTestAccProjectId()),
// 				// resource.TestCheckResourceAttr("ubicloud_firewall.testacc", "name", "tf-testacc"),
// 				// resource.TestCheckResourceAttr("ubicloud_firewall.testacc", "description", "Terraform acceptance testing"),
// 				// resource.TestCheckResourceAttr("ubicloud_firewall.testacc", "firewall_rules.#", "0"),
// 				),
// 			},
// 			// Test ImportState
// 			// {
// 			// 	ResourceName:        "ubicloud_firewall.testacc",
// 			// 	ImportState:         true,
// 			// 	ImportStateIdPrefix: fmt.Sprintf("%s,", GetTestAccProjectId()),
// 			// 	ImportStateVerify:   true,
// 			// },
// 		},
// 	})
// }

func TestAccPrivateKeyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				resource "coolify_private_key" "test" {
					name        = "TerraformAccTest"
					description = "Terraform acceptance testing"
					private_key = "FakePrivateKey\n"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("coolify_private_key.test", "uuid"),
					resource.TestCheckResourceAttr("coolify_private_key.test", "name", "TerraformAccTest"),
					resource.TestCheckResourceAttr("coolify_private_key.test", "description", "Terraform acceptance testing"),
					resource.TestCheckResourceAttr("coolify_private_key.test", "private_key", "FakePrivateKey\n"),
				),
			},
			{
				ResourceName:      "coolify_private_key.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					return s.RootModule().Resources["coolify_private_key.test"].Primary.Attributes["uuid"], nil
				},
			},
		},
	})
}

func TestAccPrivateKeyResourceUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				Config: providerConfig + `
				resource "coolify_private_key" "test" {
					name        = "TerraformAccTest"
					description = "Terraform acceptance testing"
					private_key = "FakePrivateKey\n"
				}
				`,
			},
			// Update
			{
				Config: providerConfig + `
				resource "coolify_private_key" "test" {
					name        = "TerraformAccTest"
					description = "Terraform acceptance testing"
					private_key = "NewFakePrivateKey\n"
				}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("coolify_private_key.test", plancheck.ResourceActionUpdate),
					},
				},
			},
		},
	})
}

// func TestAccPrivateKeyResourceDeleted(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: `
// 				resource "coolify_private_key" "test" {
// 					space_id = "0p38pssr0fi3"
// 					environment_id = "master"
// 					app_definition_id = "1WkQ2J9LERPtbMTdUfSHka"
// 				}
// 				`,
// 			},
// 			{
// 				Config: `
// 				resource "coolify_private_key" "test" {
// 					name        = "TerraformAccTest"
// 					description = "Terraform acceptance testing"
// 					private_key = "FakePrivateKey\n"
// 				}

// 				import {
// 					id = "0p38pssr0fi3/master/1WkQ2J9LERPtbMTdUfSHka"
// 					to = coolify_private_key.test_dup
// 				}

// 				resource "coolify_private_key" "test_dup" {
// 					space_id = "0p38pssr0fi3"
// 					environment_id = "master"
// 					app_definition_id = "1WkQ2J9LERPtbMTdUfSHka"
// 				}
// 				`,
// 			},
// 			{
// 				Config: `
// 				resource "coolify_private_key" "test" {
// 					space_id = "0p38pssr0fi3"
// 					environment_id = "master"
// 					app_definition_id = "1WkQ2J9LERPtbMTdUfSHka"
// 				}
// 				`,
// 				ExpectNonEmptyPlan: true,
// 				ConfigPlanChecks: resource.ConfigPlanChecks{
// 					PostApplyPostRefresh: []plancheck.PlanCheck{
// 						plancheck.ExpectResourceAction("coolify_private_key.test", plancheck.ResourceActionCreate),
// 					},
// 				},
// 			},
// 		},
// 	})
// }

// //nolint:paralleltest
// func TestAccAppInstallationResourceDeleted(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: `
// 				resource "coolify_private_key" "test" {
// 					space_id = "0p38pssr0fi3"
// 					environment_id = "master"
// 					app_definition_id = "1WkQ2J9LERPtbMTdUfSHka"
// 				}
// 				`,
// 			},
// 			{
// 				Config: `
// 				resource "coolify_private_key" "test" {
// 					space_id = "0p38pssr0fi3"
// 					environment_id = "master"
// 					app_definition_id = "1WkQ2J9LERPtbMTdUfSHka"
// 				}

// 				import {
// 					id = "0p38pssr0fi3/master/1WkQ2J9LERPtbMTdUfSHka"
// 					to = coolify_private_key.test_dup
// 				}

// 				resource "coolify_private_key" "test_dup" {
// 					space_id = "0p38pssr0fi3"
// 					environment_id = "master"
// 					app_definition_id = "1WkQ2J9LERPtbMTdUfSHka"
// 				}
// 				`,
// 			},
// 			{
// 				Config: `
// 				resource "coolify_private_key" "test" {
// 					space_id = "0p38pssr0fi3"
// 					environment_id = "master"
// 					app_definition_id = "1WkQ2J9LERPtbMTdUfSHka"
// 				}
// 				`,
// 				ExpectNonEmptyPlan: true,
// 				ConfigPlanChecks: resource.ConfigPlanChecks{
// 					PostApplyPostRefresh: []plancheck.PlanCheck{
// 						plancheck.ExpectResourceAction("coolify_private_key.test", plancheck.ResourceActionCreate),
// 					},
// 				},
// 			},
// 		},
// 	})
// }

// func TestAccAppInstallationResourceImport(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: `
// 				resource "coolify_private_key" "test" {
// 					space_id = "0p38pssr0fi3"
// 					environment_id = "master"
// 					app_definition_id = "1WkQ2J9LERPtbMTdUfSHka"
// 				}
// 				`,
// 			},
// 			{
// 				ResourceName:  "coolify_private_key.test",
// 				ImportState:   true,
// 				ImportStateId: "a",
// 				ExpectError:   regexp.MustCompile(`Resource Import Passthrough Multipart ID Mismatch`),
// 			},
// 			{
// 				ResourceName:  "coolify_private_key.test",
// 				ImportState:   true,
// 				ImportStateId: "a/b",
// 				ExpectError:   regexp.MustCompile(`Resource Import Passthrough Multipart ID Mismatch`),
// 			},
// 			{
// 				ResourceName:  "coolify_private_key.test",
// 				ImportState:   true,
// 				ImportStateId: "a/b/c/d",
// 				ExpectError:   regexp.MustCompile(`Resource Import Passthrough Multipart ID Mismatch`),
// 			},
// 			{
// 				ResourceName:  "coolify_private_key.test",
// 				ImportState:   true,
// 				ImportStateId: "0p38pssr0fi3/master/1WkQ2J9LERPtbMTdUfSHka",
// 			},
// 		},
// 	})
// }

// //nolint:paralleltest
// func TestAccAppInstallationResourceImportNotFound(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: `
// 				resource "coolify_private_key" "test" {
// 					space_id = "0p38pssr0fi3"
// 					environment_id = "test"
// 					app_definition_id = "nonexistent"
// 				}
// 				`,
// 				PlanOnly:           true,
// 				ExpectNonEmptyPlan: true,
// 			},
// 			{
// 				ResourceName:  "coolify_private_key.test",
// 				ImportState:   true,
// 				ImportStateId: "0p38pssr0fi3/test/nonexistent",
// 				ExpectError:   regexp.MustCompile(`Cannot import non-existent remote object`),
// 			},
// 		},
// 	})
// }

// //nolint:paralleltest
// func TestAccAppInstallationResourceCreateNotFound(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: `
// 				resource "coolify_private_key" "test" {
// 					space_id = "0p38pssr0fi3"
// 					environment_id = "master"
// 					app_definition_id = "12345"
// 				}
// 				`,
// 				ExpectError: regexp.MustCompile(`Failed to create app installation`),
// 			},
// 		},
// 	})
// }

// //nolint:paralleltest
// func TestAccAppInstallationResourceDeleted(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: `
// 				resource "coolify_private_key" "test" {
// 					space_id = "0p38pssr0fi3"
// 					environment_id = "master"
// 					app_definition_id = "1WkQ2J9LERPtbMTdUfSHka"
// 				}
// 				`,
// 			},
// 			{
// 				Config: `
// 				resource "coolify_private_key" "test" {
// 					space_id = "0p38pssr0fi3"
// 					environment_id = "master"
// 					app_definition_id = "1WkQ2J9LERPtbMTdUfSHka"
// 				}

// 				import {
// 					id = "0p38pssr0fi3/master/1WkQ2J9LERPtbMTdUfSHka"
// 					to = coolify_private_key.test_dup
// 				}

// 				resource "coolify_private_key" "test_dup" {
// 					space_id = "0p38pssr0fi3"
// 					environment_id = "master"
// 					app_definition_id = "1WkQ2J9LERPtbMTdUfSHka"
// 				}
// 				`,
// 			},
// 			{
// 				Config: `
// 				resource "coolify_private_key" "test" {
// 					space_id = "0p38pssr0fi3"
// 					environment_id = "master"
// 					app_definition_id = "1WkQ2J9LERPtbMTdUfSHka"
// 				}
// 				`,
// 				ExpectNonEmptyPlan: true,
// 				ConfigPlanChecks: resource.ConfigPlanChecks{
// 					PostApplyPostRefresh: []plancheck.PlanCheck{
// 						plancheck.ExpectResourceAction("coolify_private_key.test", plancheck.ResourceActionCreate),
// 					},
// 				},
// 			},
// 		},
// 	})
// }
