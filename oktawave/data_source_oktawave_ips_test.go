package oktawave

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestOktawave_DataSource_Ips(t *testing.T) {
	resourcesConfig := `
resource "oktawave_ip" "test-ip1" {
	subregion_id = 1
	comment = "test-ip1"
}
resource "oktawave_ip" "test-ip2" {
	subregion_id = 4
	comment = "test-ip2"
}
resource "oktawave_ip" "test-ip3" {
	subregion_id = 5
	comment = "test-ip3"
}`

	dataSourceConfig := `
data "oktawave_ips" "test-datasource-ips" {
	filter {
		key = "subregion_id"
		values = [1,5]
	}

	filter {
		key = "comment"
		values = ["test-ip1", "test-ip3"]
	}
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIpDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.oktawave_ips.test-datasource-ips", "items.#", "2"),
					resource.TestCheckResourceAttr("data.oktawave_ips.test-datasource-ips", "items.0.subregion_id", "1"),
					resource.TestCheckResourceAttr("data.oktawave_ips.test-datasource-ips", "items.1.subregion_id", "5"),
					resource.TestCheckResourceAttrPair("data.oktawave_ips.test-datasource-ips", "items.0.id", "oktawave_ip.test-ip1", "id"),
					resource.TestCheckResourceAttrPair("data.oktawave_ips.test-datasource-ips", "items.1.id", "oktawave_ip.test-ip3", "id"),
				),
			},
		},
	})
}
