package oktawave

import (
	"fmt"
	"math"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestOktawave_DataSource_Ip(t *testing.T) {
	resourcesConfig := `
resource "oktawave_ip" "test-ip1" {
	subregion_id = 1
}`

	dataSourceConfig := `
data "oktawave_ip" "test-datasource-ip1" {
	address = oktawave_ip.test-ip1.address
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
					resource.TestCheckResourceAttr("data.oktawave_ip.test-datasource-ip1", "subregion_id", "1"),
					resource.TestCheckResourceAttrPair("data.oktawave_ip.test-datasource-ip1", "id", "oktawave_ip.test-ip1", "id"),
				),
			},
		},
	})
}

func testAccCheckIpDatasourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ClientConfig).odkClient
	auth := testAccProvider.Meta().(*ClientConfig).odkAuth

	params := map[string]interface{}{
		"pageSize": int32(math.MaxInt16),
	}
	ips, _, err := client.FloatingIPsApi.FloatingIpsGetIps(*auth, params)
	if err != nil {
		return fmt.Errorf("Get ips request failed. Caused by: %s.", err)
	}

	errors := []interface{}{}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oktawave_ip" {
			continue
		}

		// id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		// if err != nil {
		// 	return fmt.Errorf("Failed to parse resource id %s.", rs.Primary.ID)
		// }

		address := rs.Primary.ID
		for _, ip := range ips.Items {
			if ip.Address == address {
				errors = append(errors, fmt.Sprintf("Ip with id %d not destroyed correctly.", ip.Id))
				break
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("Test failed to cleanup. %s", errors)
	}

	return nil
}
