package oktawave

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestOktawave_DataSource_Groups(t *testing.T) {
	resourcesConfig := `
resource "oktawave_group" "group1" {
	name = "test-group1"
}

resource "oktawave_group" "group2" {
	name = "test-group2"
}

resource "oktawave_group" "group3" {
	name = "test-group3"
}`

	dataSourceConfig := `
data "oktawave_groups" "test-groups" {
	filter {
		key = "name"
		values = ["test-group1", "test-group3"]
	}

	filter {
		key = "is_load_balancer"
		values = ["false"]
	}
}`

	// TODO test empty datasource
	// TODO test errors

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
					resource.TestCheckResourceAttr("data.oktawave_groups.test-groups", "groups.#", "2"),
					resource.TestCheckResourceAttr("data.oktawave_groups.test-groups", "groups.0.name", "test-group1"),
					resource.TestCheckResourceAttr("data.oktawave_groups.test-groups", "groups.1.name", "test-group3"),
					resource.TestCheckResourceAttrPair("data.oktawave_groups.test-groups", "groups.0.id", "oktawave_group.group1", "id"),
					resource.TestCheckResourceAttrPair("data.oktawave_groups.test-groups", "groups.1.id", "oktawave_group.group3", "id"),
				),
			},
		},
	})
}
