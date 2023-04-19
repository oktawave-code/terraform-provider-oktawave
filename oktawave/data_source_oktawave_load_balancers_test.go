package oktawave

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestOktawave_DataSource_LoadBalancers(t *testing.T) {
	resourcesConfig := `
resource "oktawave_group" "test-lb-group1" {
	name = "test-lb-group1"
}

resource "oktawave_load_balancer" "test-lbs-lb1" {
	depends_on = [oktawave_group.test-lb-group1]
	group_id = oktawave_group.test-lb-group1.id
}


`

	dataSourceConfig := `
data "oktawave_load_balancers" "test_datasource_lbs" {
	filter {
		key = "group_name"
		values = ["test-lb-group1", "test-lb-group3"]
	}
}
`

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
					resource.TestCheckResourceAttr("data.oktawave_load_balancers.test_datasource_lbs", "items.#", "1"),
					resource.TestCheckResourceAttr("data.oktawave_load_balancers.test_datasource_lbs", "items.0.group_name", "test-lb-group1"),
					// resource.TestCheckResourceAttr("data.oktawave_load_balancers.test_datasource_lbs", "items.1.name", "test-lb-group3"),
					resource.TestCheckResourceAttrPair("data.oktawave_load_balancers.test_datasource_lbs", "items.0.group_id", "oktawave_group.test-lb-group1", "id"),
					// resource.TestCheckResourceAttrPair("data.oktawave_load_balancers.test_datasource_lbs", "items.1.group_id", "oktawave_group.test-lb-group3", "id"),
				),
			},
		},
	})
}
