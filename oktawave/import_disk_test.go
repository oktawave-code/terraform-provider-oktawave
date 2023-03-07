package oktawave

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktawaveDisk_importBasic(t *testing.T) {
	resourceName := "oktawave_disk.test-disk"
	diskConfig := `
resource "oktawave_disk" "test-disk" {
	name = "disk1"
	tier_id = 48
	subregion_id = 1
	capacity = 5
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDiskDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: diskConfig,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
