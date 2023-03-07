package oktawave

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestOktawave_DataSource_OksCluster(t *testing.T) {
	resourcesConfig := `
resource "oktawave_oks_cluster" "test-oks-cluster" {
	name = "tftest"
	version = "1.21"
}`

	dataSourceConfig := `
data "oktawave_oks_cluster" "test-oks-cluster" {
	name = oktawave_oks_cluster.test-oks-cluster.id
}`

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
					resource.TestCheckResourceAttr("data.oktawave_oks_cluster.test-oks-cluster", "version", "1.21"),
					resource.TestCheckResourceAttrPair("data.oktawave_oks_cluster.test-oks-cluster", "name", "oktawave_oks_cluster.test-oks-cluster", "id"),
				),
			},
		},
	})
}

func testAccCheckOksClusterDatasourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ClientConfig).oksClient
	auth := testAccProvider.Meta().(*ClientConfig).oksAuth

	clusters, _, err := client.ClustersApi.ClustersGet(*auth)
	if err != nil {
		return fmt.Errorf("Get oks clusters request failed. Caused by: %s.", err)
	}

	errors := []interface{}{}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oktawave_oks_cluster" {
			continue
		}

		id := rs.Primary.ID

		for _, cluster := range clusters {
			if cluster.Name == id {
				errors = append(errors, fmt.Sprintf("Oks cluster with name %s not destroyed correctly.", cluster.Name))
			}
		}
	}

	err = testAccCheckGroupDatasourceDestroy(s)
	if err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("Test failed to cleanup. %s", errors)
	}

	return nil
}
