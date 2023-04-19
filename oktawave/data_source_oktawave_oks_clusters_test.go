package oktawave

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestOktawave_DataSource_OksClusters(t *testing.T) {
	resourcesConfig := `
resource "oktawave_oks_cluster" "test-oks-cluster1" {
	name = "tftest1"
	version = "1.21"
}
`
	// FIXME: problem with cluster status causes tests to hang or fail too often
	//
	// resource "oktawave_oks_cluster" "test-oks-cluster2" {
	// 	name = "tftest2"
	// 	version = "1.21"
	// }

	// resource "oktawave_oks_cluster" "test-oks-cluster3" {
	// 	name = "tftest3"
	// 	version = "1.21"
	// }
	// `

	dataSourceConfig := `
data "oktawave_oks_clusters" "test-oks-clusters" {
	filter {
		key = "name"
		values = [oktawave_oks_cluster.test-oks-cluster1.id]
	}
}
`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOksClusterDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.oktawave_oks_clusters.test-oks-clusters", "items.#", "1"),
					resource.TestCheckResourceAttrPair("data.oktawave_oks_clusters.test-oks-clusters", "items.0.name", "oktawave_oks_cluster.test-oks-cluster1", "id"),
					// resource.TestCheckResourceAttrPair("data.oktawave_oks_clusters.test-oks-clusters", "items.1.name", "oktawave_oks_cluster.test-oks-cluster3", "name"),
				),
			},
		},
	})
}
