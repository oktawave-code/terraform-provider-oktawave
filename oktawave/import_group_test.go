package oktawave

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktawaveGroup_importBasic(t *testing.T) {
	resourceName := "oktawave_group.test-group"
	groupConfig := `
resource "oktawave_group" "test-group" {
	name = "test-group"
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGroupDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: groupConfig,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
