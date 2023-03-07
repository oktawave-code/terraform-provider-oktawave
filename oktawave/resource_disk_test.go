package oktawave

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/oktawave-code/odk"
)

func TestAccOktawaveDisk_Basic(t *testing.T) {
	var disk odk.Disk
	diskConfig := `
resource "oktawave_disk" "test-disk" {
	name = "disk1"
	tier_id = 48
	subregion_id = 1
	capacity = 5
}
`

	updateConfig := `
resource "oktawave_disk" "test-disk" {
	name = "disk1-updated"
	tier_id = 49
	subregion_id = 4
	capacity = 6
}
`
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDiskDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: diskConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDiskExists("oktawave_disk.test-disk", &disk),
					// testAccCheckDiskAttributes(&disk),
					resource.TestCheckResourceAttr("oktawave_disk.test-disk", "name", "disk1"),
					resource.TestCheckResourceAttr("oktawave_disk.test-disk", "tier_id", "48"),
					resource.TestCheckResourceAttr("oktawave_disk.test-disk", "subregion_id", "1"),
					resource.TestCheckResourceAttr("oktawave_disk.test-disk", "capacity", "5"),
					resource.TestCheckResourceAttrSet("oktawave_disk.test-disk", "creation_user_id"),
					resource.TestCheckResourceAttrSet("oktawave_disk.test-disk", "creation_date"),
					resource.TestCheckResourceAttrSet("oktawave_disk.test-disk", "is_shared"),
					resource.TestCheckResourceAttrSet("oktawave_disk.test-disk", "is_locked"),
					resource.TestCheckResourceAttrSet("oktawave_disk.test-disk", "locking_date"),
					resource.TestCheckResourceAttrSet("oktawave_disk.test-disk", "is_freemium"),
					resource.TestCheckNoResourceAttr("oktawave_disk.test-disk", "shared_disk_type_id"),
					resource.TestCheckResourceAttr("oktawave_disk.test-disk", "instance_ids.#", "0"),
				),
			},
			{
				Config: updateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("oktawave_disk.test-disk", "name", "disk1-updated"),
					resource.TestCheckResourceAttr("oktawave_disk.test-disk", "tier_id", "49"),
					resource.TestCheckResourceAttr("oktawave_disk.test-disk", "subregion_id", "4"),
					resource.TestCheckResourceAttr("oktawave_disk.test-disk", "capacity", "6"),
				),
			},
		},
	})
}

func testAccCheckDiskExists(resourceName string, disk *odk.Disk) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Disk ID is not set")
		}

		client := testAccProvider.Meta().(*ClientConfig).odkClient
		auth := testAccProvider.Meta().(*ClientConfig).odkAuth

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Failed to parse resource id %s.", rs.Primary.ID)
		}

		object, _, err := client.OVSApi.DisksGet(*auth, int32(id), nil)
		if err != nil {
			return err
		}

		if strconv.Itoa(int(object.Id)) != rs.Primary.ID {
			return fmt.Errorf("Disk not found")
		}

		*disk = object
		return nil
	}
}

// func testAccCheckdiskAttributes(disk *odk.Disk) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 	  if *widget.active != true {
// 		return fmt.Errorf("widget is not active")
// 	  }

// 	  return nil
// 	}
// }
