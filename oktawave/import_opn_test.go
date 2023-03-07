package oktawave

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktawaveOpn_importBasic(t *testing.T) {
	resourceName := "oktawave_opn.test-opn"
	opnConfig := `
resource "oktawave_opn" "test-opn" {
	name = "test-opn"
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOpnDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: opnConfig,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
