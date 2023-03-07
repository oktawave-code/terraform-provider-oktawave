package oktawave

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktawaveLoadBalancer_importBasic(t *testing.T) {
	resourceName := "oktawave_load_balancer.test-load_balancer"
	load_balancerConfig := `
resource "oktawave_group" "test-lb-group" {
	name = "test-lb-group"
}

resource "oktawave_load_balancer" "test-load_balancer" {
	depends_on = [oktawave_group.test-lb-group]
	group_id = oktawave_group.test-lb-group.id
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoadBalancerDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: load_balancerConfig,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
