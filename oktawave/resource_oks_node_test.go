package oktawave

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	oks "github.com/oktawave-code/oks-sdk"
)

func TestAccOktawaveOksNode_Basic(t *testing.T) {
	var oksNode oks.K44sInstance
	clusterConfig := `
resource "oktawave_oks_cluster" "test-cluster" {
	name = "tftest"
	version = "1.21"
}

resource "oktawave_oks_node" "test-node" {
	depends_on = [oktawave_oks_cluster.test-cluster]
	cluster_id = oktawave_oks_cluster.test-cluster.id
	subregion_id = 1
	type_id = 34
}
`
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOksClusterDatasourceDestroy, // todo - node destroy
		Steps: []resource.TestStep{
			{
				Config: clusterConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckOksNodeExists("oktawave_oks_cluster.test-cluster", "oktawave_oks_node.test-node", &oksNode),
					resource.TestCheckResourceAttrPair("oktawave_oks_node.test-node", "cluster_id", "oktawave_oks_cluster.test-cluster", "id"),
					resource.TestCheckResourceAttr("oktawave_oks_node.test-node", "subregion_id", "1"),
					resource.TestCheckResourceAttr("oktawave_oks_node.test-node", "type_id", "34"),
					resource.TestCheckResourceAttrSet("oktawave_oks_node.test-node", "name"),
					resource.TestCheckResourceAttrSet("oktawave_oks_node.test-node", "creation_date"),
					resource.TestCheckResourceAttrSet("oktawave_oks_node.test-node", "status_id"),
					resource.TestCheckResourceAttrSet("oktawave_oks_node.test-node", "ip_address"),
					resource.TestCheckResourceAttrSet("oktawave_oks_node.test-node", "total_disks_capacity"),
					resource.TestCheckResourceAttrSet("oktawave_oks_node.test-node", "cpu_number"),
					resource.TestCheckResourceAttrSet("oktawave_oks_node.test-node", "ram_mb"),
				),
			},
		},
	})
}

func testAccCheckOksNodeExists(clusterResourceName string, resourceName string, node *oks.K44sInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rsCluster, ok := s.RootModule().Resources[clusterResourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", clusterResourceName)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rsCluster.Primary.ID == "" {
			return fmt.Errorf("Oks cluster ID is not set")
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Oks node ID is not set")
		}

		client := testAccProvider.Meta().(*ClientConfig).oksClient
		auth := testAccProvider.Meta().(*ClientConfig).oksAuth

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Failed to parse resource id %s.", rs.Primary.ID)
		}

		nodes, _, err := client.ClustersApi.ClustersInstancesNameGet(*auth, rsCluster.Primary.ID)
		if err != nil {
			return err
		}

		filterFn := func(node oks.K44sInstance) bool {
			return int64(node.Id) == id
		}

		results := filter(nodes, filterFn)

		if len(results) != 1 {
			return fmt.Errorf("Node not found")
		}

		*node = results[0]
		return nil
	}
}
