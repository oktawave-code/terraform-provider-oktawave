package oktawave

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestOktawave_DataSource_Instance(t *testing.T) {
	resourcesConfig := `
resource "oktawave_ip" "test-ip1" {
	subregion_id = 1
}

resource "oktawave_instance" "test-instance" {
	public_ips = [oktawave_ip.test-ip1.id]
	depends_on = [oktawave_ip.test-ip1]
	name = "test-instance"
	subregion_id = 1
	system_disk_class_id = 48
	template_id = 1021
	type_id = 1047
}`

	dataSourceConfig := `
data "oktawave_instance" "test-instance" {
	id = oktawave_instance.test-instance.id
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.oktawave_instance.test-instance", "name", "test-instance"),
					resource.TestCheckResourceAttrPair("data.oktawave_instance.test-instance", "id", "oktawave_instance.test-instance", "id"),
				),
			},
		},
	})
}

func testAccCheckInstanceDatasourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ClientConfig).odkClient
	auth := testAccProvider.Meta().(*ClientConfig).odkAuth

	params := map[string]interface{}{
		"pageSize": int32(math.MaxInt16),
	}
	list, _, err := client.OCIApi.InstancesGet(*auth, params)
	if err != nil {
		return fmt.Errorf("Get instances request failed. Caused by: %s.", err)
	}

	errors := []interface{}{}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oktawave_instance" {
			continue
		}

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Failed to parse resource id %s.", rs.Primary.ID)
		}

		for _, item := range list.Items {
			if int64(item.Id) == id {
				errors = append(errors, fmt.Sprintf("Instance with id %d not destroyed correctly.", item.Id))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("Test failed to cleanup. %s", errors)
	}

	return nil
}
