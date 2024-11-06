package provider_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-coolify/internal/provider"
)

func TestAccTeamDataSource(t *testing.T) {
	resName := "data.coolify_team.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// With ID
			{
				Config: `data "coolify_team" "test" {
					id = 1
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "id", "1"),
					resource.TestCheckResourceAttr(resName, "name", "Test"),
					resource.TestCheckResourceAttr(resName, "personal_team", "false"),
					resource.TestCheckResourceAttr(resName, "members.#", "1"),
				),
			},
			// Without ID, current API Key authenticated team
			{
				Config: `data "coolify_team" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "id", "0"),
					resource.TestCheckResourceAttr(resName, "name", "Root Team"),
					resource.TestCheckResourceAttr(resName, "members.#", "1"),
					resource.TestCheckResourceAttr(resName, "personal_team", "true"),
				),
			},
		},
	})
}

func TestTeamDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	ds := provider.NewTeamDataSource()
	resp := &datasource.SchemaResponse{}
	ds.Schema(ctx, datasource.SchemaRequest{}, resp)

	// Test id is optional
	idAttr := resp.Schema.Attributes["id"].(schema.Int64Attribute)
	if idAttr.Required {
		t.Error("id field should not be marked as required in schema")
	}
	if !idAttr.Optional {
		t.Error("id field should be marked as optional in schema")
	}

	// Test sensitive fields
	sensitiveFields := []string{"discord_webhook_url", "smtp_password", "telegram_token", "resend_api_key"}
	for _, field := range sensitiveFields {
		if attr, ok := resp.Schema.Attributes[field].(schema.StringAttribute); ok {
			if !attr.Sensitive {
				t.Errorf("%s field should be marked as sensitive in schema", field)
			}
		} else {
			t.Errorf("%s field should be a string attribute in schema", field)
		}
	}
}
