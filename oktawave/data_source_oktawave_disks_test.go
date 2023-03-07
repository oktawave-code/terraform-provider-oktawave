package oktawave

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestOktawave_DataSource_Disks(t *testing.T) {
	resourcesConfig := `
resource "oktawave_disk" "disk1" {
	name = "disk1"
	tier_id = 48
	subregion_id = 1
}

resource "oktawave_disk" "disk2" {
	name = "disk2"
	tier_id = 48
	subregion_id = 1
}


resource "oktawave_disk" "disk3" {
	name = "disk3"
	tier_id = 48
	subregion_id = 5
}
`

	dataSourceConfig := `
data "oktawave_disks" "disks" {
	filter {
		key = "name"
		values = ["disk1", "disk3"]
	}
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
					resource.TestCheckResourceAttr("data.oktawave_disks.disks", "disks.#", "2"),
					resource.TestCheckResourceAttr("data.oktawave_disks.disks", "disks.0.name", "disk1"),
					resource.TestCheckResourceAttr("data.oktawave_disks.disks", "disks.1.name", "disk3"),
					resource.TestCheckResourceAttrPair("data.oktawave_disks.disks", "disks.0.id", "oktawave_disk.disk1", "id"),
					resource.TestCheckResourceAttrPair("data.oktawave_disks.disks", "disks.1.id", "oktawave_disk.disk3", "id"),
				),
			},
		},
	})
}
