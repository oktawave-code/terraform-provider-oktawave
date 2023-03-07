package oktawave

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestOktawave_DataSource_Disk(t *testing.T) {
	resourcesConfig := `
resource "oktawave_disk" "disk1" {
	name = "disk1"
	tier_id = 48
	subregion_id = 1
}
`

	dataSourceConfig := `
data "oktawave_disk" "disk1" {
	id = oktawave_disk.disk1.id
}
`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDiskDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.oktawave_disk.disk1", "name", "disk1"),
					resource.TestCheckResourceAttrPair("data.oktawave_disk.disk1", "id", "oktawave_disk.disk1", "id"),
				),
			},
		},
	})
}

func testAccCheckDiskDatasourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ClientConfig).odkClient
	auth := testAccProvider.Meta().(*ClientConfig).odkAuth

	params := map[string]interface{}{
		"pageSize": int32(math.MaxInt16),
	}
	disks, _, err := client.OVSApi.DisksGetDisks(*auth, params)
	if err != nil {
		return fmt.Errorf("Get disks request failed. Caused by: %s.", err)
	}

	errors := []interface{}{}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oktawave_disk" {
			continue
		}

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Failed to parse resource id %s.", rs.Primary.ID)
		}

		for _, group := range disks.Items {
			if int64(group.Id) == id {
				errors = append(errors, fmt.Sprintf("Disk with id %d not destroyed correctly.", group.Id))
				break
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("Test failed to cleanup. %s", errors)
	}

	return nil
}
