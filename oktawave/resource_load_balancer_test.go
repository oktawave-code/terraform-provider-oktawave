package oktawave

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/oktawave-code/odk"
)

func TestAccOktawaveLoadBalancer_Basic(t *testing.T) {
	var load_balancer odk.LoadBalancer
	lbConfig := `
resource "oktawave_group" "test-lb-group" {
	name = "test-lb-group"
}

resource "oktawave_load_balancer" "test-load_balancer" {
	depends_on = [oktawave_group.test-lb-group]
	group_id = oktawave_group.test-lb-group.id
}`

	lbUpdateConfig := `
resource "oktawave_group" "test-lb-group" {
	name = "test-lb-group"
}

resource "oktawave_load_balancer" "test-load_balancer" {
	depends_on = [oktawave_group.test-lb-group]
	group_id = oktawave_group.test-lb-group.id
	algorithm_id = 281
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoadBalancerDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: lbConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLoadBalancerExists("oktawave_load_balancer.test-load_balancer", &load_balancer),
					// testAccCheckLoad_balancerAttributes(&load_balancer),
					resource.TestCheckResourceAttrSet("oktawave_load_balancer.test-load_balancer", "group_id"),
					resource.TestCheckResourceAttrSet("oktawave_load_balancer.test-load_balancer", "service_type_id"),
					resource.TestCheckResourceAttrSet("oktawave_load_balancer.test-load_balancer", "port_number"),
					resource.TestCheckResourceAttrSet("oktawave_load_balancer.test-load_balancer", "target_port_number"),
					resource.TestCheckResourceAttrSet("oktawave_load_balancer.test-load_balancer", "ssl_target_port_number"),
					resource.TestCheckResourceAttrSet("oktawave_load_balancer.test-load_balancer", "session_persistence_type_id"),
					resource.TestCheckResourceAttrSet("oktawave_load_balancer.test-load_balancer", "algorithm_id"),
					resource.TestCheckResourceAttrSet("oktawave_load_balancer.test-load_balancer", "ip_version_id"),
					resource.TestCheckResourceAttrSet("oktawave_load_balancer.test-load_balancer", "health_check_enabled"),
					resource.TestCheckResourceAttrSet("oktawave_load_balancer.test-load_balancer", "ssl_enabled"),
					resource.TestCheckResourceAttrSet("oktawave_load_balancer.test-load_balancer", "common_persistence_enabled"),
					resource.TestCheckResourceAttrSet("oktawave_load_balancer.test-load_balancer", "proxy_protocol_version_id"),
					resource.TestCheckResourceAttrSet("oktawave_load_balancer.test-load_balancer", "group_name"),
					resource.TestCheckResourceAttrSet("oktawave_load_balancer.test-load_balancer", "address"),
					resource.TestCheckResourceAttr("oktawave_load_balancer.test-load_balancer", "address_v6", ""),
				),
			},
			{
				Config: lbUpdateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLoadBalancerExists("oktawave_load_balancer.test-load_balancer", &load_balancer),
					// testAccCheckLoad_balancerAttributes(&load_balancer),
					resource.TestCheckResourceAttr("oktawave_load_balancer.test-load_balancer", "algorithm_id", "281"),
				),
			},
		},
	})
}

func testAccCheckLoadBalancerExists(resourceName string, load_balancer *odk.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Load_balancer ID is not set")
		}

		client := testAccProvider.Meta().(*ClientConfig).odkClient
		auth := testAccProvider.Meta().(*ClientConfig).odkAuth

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Failed to parse resource id %s.", rs.Primary.ID)
		}

		object, _, err := client.OCIGroupsApi.LoadBalancersGetLoadBalancer(*auth, int32(id), nil)
		if err != nil {
			return err
		}

		if strconv.Itoa(int(object.GroupId)) != rs.Primary.ID {
			return fmt.Errorf("Load_balancer not found")
		}

		*load_balancer = object
		return nil
	}
}

// func testAccCheckload_balancerAttributes(load_balancer *odk.Load_balancer) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 	  if *widget.active != true {
// 		return fmt.Errorf("widget is not active")
// 	  }

// 	  return nil
// 	}
// }
