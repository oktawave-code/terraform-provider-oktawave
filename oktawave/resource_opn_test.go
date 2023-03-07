package oktawave

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/oktawave-code/odk"
)

func TestAccOktawaveOpn_Basic(t *testing.T) {
	var opn odk.Opn
	opnConfig := `
resource "oktawave_opn" "test-opn" {
	name = "test-opn"
}
`

	opnUpdateConfig := `
resource "oktawave_opn" "test-opn" {
	name = "test-opn2"
}
`
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOpnDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: opnConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckOpnExists("oktawave_opn.test-opn", &opn),
					// testAccCheckOpnAttributes(&opn),
					resource.TestCheckResourceAttr("oktawave_opn.test-opn", "name", "test-opn"),
					resource.TestCheckResourceAttrSet("oktawave_opn.test-opn", "name"),
					resource.TestCheckResourceAttrSet("oktawave_opn.test-opn", "creation_user_id"),
					resource.TestCheckResourceAttrSet("oktawave_opn.test-opn", "creation_date"),
					resource.TestCheckResourceAttrSet("oktawave_opn.test-opn", "last_change_date"),
					resource.TestCheckResourceAttr("oktawave_opn.test-opn", "instance_ids.#", "0"),
				),
			},
			{
				Config: opnUpdateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckOpnExists("oktawave_opn.test-opn", &opn),
					// testAccCheckOpnAttributes(&opn),
					resource.TestCheckResourceAttr("oktawave_opn.test-opn", "name", "test-opn2"),
				),
			},
		},
	})
}

func testAccCheckOpnExists(resourceName string, opn *odk.Opn) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Opn ID is not set")
		}

		client := testAccProvider.Meta().(*ClientConfig).odkClient
		auth := testAccProvider.Meta().(*ClientConfig).odkAuth

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Failed to parse resource id %s.", rs.Primary.ID)
		}

		object, _, err := client.NetworkingApi.OpnsGet_1(*auth, int32(id), nil)
		if err != nil {
			return err
		}

		if strconv.Itoa(int(object.Id)) != rs.Primary.ID {
			return fmt.Errorf("Opn not found")
		}

		*opn = object
		return nil
	}
}

// func testAccCheckopnAttributes(opn *odk.Opn) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 	  if *widget.active != true {
// 		return fmt.Errorf("widget is not active")
// 	  }

// 	  return nil
// 	}
// }
