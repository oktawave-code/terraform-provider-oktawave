package oktawave

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktawaveIp_importBasic(t *testing.T) {
	resourceName := "oktawave_ip.test-ip"
	ipConfig := `
resource "oktawave_ip" "test-ip" {
	subregion_id = 1
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIpDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: ipConfig,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
