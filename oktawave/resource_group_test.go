package oktawave

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/oktawave-code/odk"
)

func TestAccOktawaveGroup_Basic(t *testing.T) {
	var group odk.Group
	groupConfig := `
resource "oktawave_group" "test-group" {
	name = "test-group"
	affinity_rule_type_id = 1403
}
`

	groupUpdateConfig := `
resource "oktawave_group" "test-group" {
	name = "test-group1"
	affinity_rule_type_id = 1404
}
`
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGroupDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: groupConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckGroupExists("oktawave_group.test-group", &group),
					// testAccCheckGroupAttributes(&group),
					resource.TestCheckResourceAttr("oktawave_group.test-group", "name", "test-group"),
					resource.TestCheckResourceAttr("oktawave_group.test-group", "affinity_rule_type_id", "1403"),
					resource.TestCheckResourceAttr("oktawave_group.test-group", "group_instance_ip_ids.#", "0"),
					resource.TestCheckResourceAttrSet("oktawave_group.test-group", "is_load_balancer"),
					resource.TestCheckResourceAttrSet("oktawave_group.test-group", "instances_count"),
					resource.TestCheckResourceAttrSet("oktawave_group.test-group", "schedulers_count"),
					resource.TestCheckResourceAttrSet("oktawave_group.test-group", "autoscaling_type_id"),
					resource.TestCheckResourceAttrSet("oktawave_group.test-group", "last_change_date"),
					resource.TestCheckResourceAttrSet("oktawave_group.test-group", "creation_user_id"),
				),
			},
			{
				Config: groupUpdateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckGroupExists("oktawave_group.test-group", &group),
					// testAccCheckGroupAttributes(&group),
					resource.TestCheckResourceAttr("oktawave_group.test-group", "name", "test-group1"),
					resource.TestCheckResourceAttr("oktawave_group.test-group", "affinity_rule_type_id", "1404"),
				),
			},
		},
	})
}

func TestAccOktawaveGroup_Assignments(t *testing.T) {
	var group odk.Group
	groupConfig := `
resource "oktawave_ip" "test-ip1" {
	subregion_id = 1
}

resource "oktawave_instance" "test-instance1" {
	depends_on = [oktawave_ip.test-ip1]
	name = "test-instance1"
	subregion_id = 1
	system_disk_class_id = 48
	template_id = 1021
	type_id = 1047
	public_ips = [oktawave_ip.test-ip1.id]
}

resource "oktawave_group" "test-group" {
	depends_on = [oktawave_instance.test-instance1]
	name = "test-group"
	affinity_rule_type_id = 1403

	assignment {
		instance_id = oktawave_instance.test-instance1.id
		ip_id = oktawave_ip.test-ip1.id
	}
}
`
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGroupDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: groupConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckGroupExists("oktawave_group.test-group", &group),
					// testAccCheckGroupAttributes(&group),
					resource.TestCheckResourceAttr("oktawave_group.test-group", "assignment.#", "1"),
				),
			},
		},
	})
}

func testAccCheckGroupExists(resourceName string, group *odk.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Group ID is not set")
		}

		client := testAccProvider.Meta().(*ClientConfig).odkClient
		auth := testAccProvider.Meta().(*ClientConfig).odkAuth

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Failed to parse resource id %s.", rs.Primary.ID)
		}

		object, _, err := client.OCIGroupsApi.GroupsGetGroup(*auth, int32(id), nil)
		if err != nil {
			return err
		}

		if strconv.Itoa(int(object.Id)) != rs.Primary.ID {
			return fmt.Errorf("Group not found")
		}

		*group = object
		return nil
	}
}

// func testAccCheckgroupAttributes(group *odk.Group) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 	  if *widget.active != true {
// 		return fmt.Errorf("widget is not active")
// 	  }

// 	  return nil
// 	}
// }
