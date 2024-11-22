package provider_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"terraform-provider-coolify/internal/provider"
)

func TestAccPostgresqlDatabaseResource(t *testing.T) {
	resName := "coolify_postgresql_database.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{ // Create and Read testing
				Config: `
				resource "coolify_postgresql_database" "test" {
					name        = "TerraformAccTest"
					description = "Terraform acceptance testing"

					server_uuid = "rg8ks8c"
					project_uuid = "uoswco88w8swo40k48o8kcwk"
					environment_name = "production"

					image = "postgres:16-alpine"
					postgres_db = "postgres"
					postgres_user = "postgres"
					postgres_password = "password"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "name", "TerraformAccTest"),
					resource.TestCheckResourceAttr(resName, "description", "Terraform acceptance testing"),
					resource.TestCheckResourceAttr(resName, "server_uuid", "rg8ks8c"),
					resource.TestCheckResourceAttr(resName, "project_uuid", "uoswco88w8swo40k48o8kcwk"),
					resource.TestCheckResourceAttr(resName, "environment_name", "production"),
					resource.TestCheckResourceAttr(resName, "instant_deploy", "false"),

					resource.TestCheckResourceAttrSet(resName, "uuid"),
					resource.TestCheckResourceAttrSet(resName, "internal_db_url"),
				),
			},
			{ // ImportState testing
				ResourceName:                         resName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "uuid",
				ExpectError: regexp.MustCompile(
					`("instant_deploy")`,
				),
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					r := s.RootModule().Resources[resName].Primary.Attributes
					return fmt.Sprintf("%s/%s/%s/%s",
						r["server_uuid"],
						r["project_uuid"],
						r["environment_name"],
						r["uuid"],
					), nil
				},
			},
			{ // Update and Read testing
				Config: `
				resource "coolify_postgresql_database" "test" {
					name        = "TerraformAccTestUpdated"
					description = "Terraform acceptance testing"

					server_uuid = "rg8ks8c"
					project_uuid = "uoswco88w8swo40k48o8kcwk"
					environment_name = "production"

					image = "postgres:16-alpine"
					postgres_db = "postgres"
					postgres_user = "postgres"
					postgres_password = "password"

					is_public = false
					// public_port = 1024 
					instant_deploy = false

					postgres_initdb_args = "()"
					postgres_host_auth_method = "()"
					postgres_conf	= " "

					limits_memory = "0"
					limits_memory_swap = "0"
					limits_memory_swappiness = "60"
					limits_memory_reservation = "0"
					limits_cpus = "0"
					// limits_cpuset = null
					limits_cpu_shares = "1024"
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
					resource.TestCheckResourceAttrSet(resName, "internal_db_url"),
					resource.TestCheckResourceAttr(resName, "name", "TerraformAccTestUpdated"),
					resource.TestCheckResourceAttr(resName, "description", "Terraform acceptance testing"),
					resource.TestCheckResourceAttr(resName, "server_uuid", "rg8ks8c"),
				),
			},
		},
	})
}

func TestDatabaseResourceSchema(t *testing.T) {
	ctx := context.Background()
	rs := provider.NewPostgresqlDatabaseResource()
	resp := &tfresource.SchemaResponse{}
	rs.Schema(ctx, tfresource.SchemaRequest{}, resp)
}
