package oktawave

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestOktawave_DataSource_LoadBalancer(t *testing.T) {
	resourcesConfig := `
resource "oktawave_group" "test-lb-group" {
	name = "test-lb-group"
}

resource "oktawave_load_balancer" "test-lb1" {
	depends_on = [oktawave_group.test-lb-group]
	group_id = oktawave_group.test-lb-group.id
}`

	dataSourceConfig := `
data "oktawave_load_balancer" "test-lb" {
	group_id = oktawave_group.test-lb-group.id
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoadBalancerDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.oktawave_load_balancer.test-lb", "group_name", "test-lb-group"),
					resource.TestCheckResourceAttrPair("data.oktawave_load_balancer.test-lb", "group_id", "oktawave_load_balancer.test-lb1", "group_id"),
				),
			},
		},
	})
}

func testAccCheckLoadBalancerDatasourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ClientConfig).odkClient
	auth := testAccProvider.Meta().(*ClientConfig).odkAuth

	params := map[string]interface{}{
		"pageSize": int32(math.MaxInt16),
	}
	lbs, _, err := client.OCIGroupsApi.GroupsGetLoadBalancers(*auth, params)
	if err != nil {
		return fmt.Errorf("Get load balancers request failed. Caused by: %s.", err)
	}

	errors := []interface{}{}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oktawave_load_balancer" {
			continue
		}

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Failed to parse resource id %s.", rs.Primary.ID)
		}

		for _, lb := range lbs.Items {
			if int64(lb.GroupId) == id {
				errors = append(errors, fmt.Sprintf("Load Balancer for group with id %d not destroyed correctly.", lb.GroupId))
			}
		}
	}

	err = testAccCheckGroupDatasourceDestroy(s)
	if err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("Test failed to cleanup. %s", errors)
	}

	return nil
}
