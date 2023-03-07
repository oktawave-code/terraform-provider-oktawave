package oktawave

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	oks "github.com/oktawave-code/oks-sdk"
)

func TestAccOktawaveOksCluster_Basic(t *testing.T) {
	var oksCluster oks.K44SClusterDetailsDto
	clusterConfig := `
resource "oktawave_oks_cluster" "test-cluster" {
	name = "tftest"
	version = "1.21"
}
`
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOksClusterDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: clusterConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckOksClusterExists("oktawave_oks_cluster.test-cluster", &oksCluster),
					resource.TestCheckResourceAttr("oktawave_oks_cluster.test-cluster", "name", "tftest"),
					resource.TestCheckResourceAttr("oktawave_oks_cluster.test-cluster", "version", "1.21"),
					resource.TestCheckResourceAttrSet("oktawave_oks_cluster.test-cluster", "creation_date"),
				),
			},
		},
	})
}

func testAccCheckOksClusterExists(resourceName string, cluster *oks.K44SClusterDetailsDto) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Oks cluster ID is not set")
		}

		client := testAccProvider.Meta().(*ClientConfig).oksClient
		auth := testAccProvider.Meta().(*ClientConfig).oksAuth

		id := rs.Primary.ID

		object, _, err := client.ClustersApi.ClustersNameGet(*auth, id)
		if err != nil {
			return err
		}

		if object.Name != rs.Primary.ID {
			return fmt.Errorf("Cluster not found")
		}

		*cluster = object
		return nil
	}
}
