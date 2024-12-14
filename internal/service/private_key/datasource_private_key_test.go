package private_key_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-coolify/internal/acctest"
)

func TestAccPrivatekeyDataSource(t *testing.T) {
	randomName := acctest.GetRandomResourceName("pk-ds")
	resName := "data.coolify_private_key." + randomName
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateKeyDataSourceConfig(randomName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "id", "0"),
					resource.TestCheckResourceAttr(resName, "name", "localhost's key"),
					resource.TestCheckResourceAttrSet(resName, "private_key"),
					resource.TestCheckResourceAttr(resName, "uuid", acctest.PrivateKeyUUID),
				),
			},
		},
	})
}

func testAccPrivateKeyDataSourceConfig(randomName string) string {
	return fmt.Sprintf(`
		data "coolify_private_key" "%[1]s" {
			uuid = "%[2]s"
		}
	`,
		randomName,
		acctest.PrivateKeyUUID,
	)
}
