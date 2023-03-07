package oktawave

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestOktawave_DataSource_OksNode(t *testing.T) {
	resourcesConfig := `
resource "oktawave_oks_cluster" "test-cluster" {
	name = "tftest"
	version = "1.21"
}

resource "oktawave_oks_node" "test-node" {
	depends_on = [oktawave_oks_cluster.test-cluster]
	cluster_id = oktawave_oks_cluster.test-cluster.id
	subregion_id = 1
	type_id = 34
}`

	dataSourceConfig := `
data "oktawave_oks_node" "test-oks-node" {
	cluster_id = oktawave_oks_cluster.test-cluster.id
	id = oktawave_oks_node.test-node.id
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// CheckDestroy: testAccCheckOksNodeDatasourceDestroy, todo
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.oktawave_oks_node.test-oks-node", "subregion_id", "1"),
					resource.TestCheckResourceAttrPair("data.oktawave_oks_node.test-oks-node", "id", "oktawave_oks_node.test-node", "id"),
				),
			},
		},
	})
}

// func testAccCheckOksNodeDatasourceDestroy(s *terraform.State) error {
// 	client := testAccProvider.Meta().(*ClientConfig).oksClient
// 	auth := testAccProvider.Meta().(*ClientConfig).oksAuth

// 	clusters, _, err := client.ClustersApi.ClustersGet(*auth)
// 	if err != nil {
// 		return fmt.Errorf("Get oks clusters request failed. Caused by: %s.", err)
// 	}

// 	errors := []interface{}{}
// 	for _, rs := range s.RootModule().Resources {
// 		if rs.Type != "oktawave_oks_node" {
// 			continue
// 		}

// 		id := rs.Primary.ID

// 		for _, cluster := range clusters {
// 			if cluster.Name == id {
// 				errors = append(errors, fmt.Sprintf("Oks cluster with name %s not destroyed correctly.", cluster.Name))
// 			}
// 		}
// 	}

// 	err = testAccCheckGroupDatasourceDestroy(s)
// 	if err != nil {
// 		errors = append(errors, err)
// 	}

// 	if len(errors) > 0 {
// 		return fmt.Errorf("Test failed to cleanup. %s", errors)
// 	}

// 	return nil
// }
