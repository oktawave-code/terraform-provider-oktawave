package oktawave

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestOktawave_DataSource_Opn(t *testing.T) {
	resourcesConfig := `
resource "oktawave_opn" "test-opn" {
	name = "test-opn"
}
`

	dataSourceConfig := `
data "oktawave_opn" "test-opn" {
	id = oktawave_opn.test-opn.id
}
`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOpnDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.oktawave_opn.test-opn", "name", "test-opn"),
					resource.TestCheckResourceAttrPair("data.oktawave_opn.test-opn", "id", "oktawave_opn.test-opn", "id"),
				),
			},
		},
	})
}

func testAccCheckOpnDatasourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ClientConfig).odkClient
	auth := testAccProvider.Meta().(*ClientConfig).odkAuth

	params := map[string]interface{}{
		"pageSize": int32(math.MaxInt16),
	}
	opns, _, err := client.NetworkingApi.OpnsGet(*auth, params)
	if err != nil {
		return fmt.Errorf("Get opns request failed. Caused by: %s.", err)
	}

	errors := []interface{}{}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oktawave_opn" {
			continue
		}

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Failed to parse resource id %s.", rs.Primary.ID)
		}

		for _, opn := range opns.Items {
			if int64(opn.Id) == id {
				errors = append(errors, fmt.Sprintf("Opn with id %d not destroyed correctly.", opn.Id))
				break
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("Test failed to cleanup. %s", errors)
	}

	return nil
}
