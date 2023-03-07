package oktawave

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestOktawave_DataSource_Group(t *testing.T) {
	resourcesConfig := `
resource "oktawave_group" "group1" {
	name = "test-group"
}`

	dataSourceConfig := `
data "oktawave_group" "test-group" {
	id = oktawave_group.group1.id
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGroupDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.oktawave_group.test-group", "name", "test-group"),
					resource.TestCheckResourceAttrPair("data.oktawave_group.test-group", "id", "oktawave_group.group1", "id"),
				),
			},
		},
	})
}

func testAccCheckGroupDatasourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ClientConfig).odkClient
	auth := testAccProvider.Meta().(*ClientConfig).odkAuth

	params := map[string]interface{}{
		"pageSize": int32(math.MaxInt16),
	}
	groups, _, err := client.OCIGroupsApi.GroupsGetGroups(*auth, params)
	if err != nil {
		return fmt.Errorf("Get groups request failed. Caused by: %s.", err)
	}

	errors := []interface{}{}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oktawave_group" {
			continue
		}

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Failed to parse resource id %s.", rs.Primary.ID)
		}

		for _, group := range groups.Items {
			if int64(group.Id) == id {
				errors = append(errors, fmt.Sprintf("Group with id %d not destroyed correctly.", group.Id))
				break
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("Test failed to cleanup. %s", errors)
	}

	return nil
}
