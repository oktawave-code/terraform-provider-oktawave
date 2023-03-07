package oktawave

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestOktawave_DataSource_Opns(t *testing.T) {
	resourcesConfig := `
resource "oktawave_opn" "test-opn1" {
	name = "test-opn1"
}

resource "oktawave_opn" "test-opn2" {
	name = "test-opn2"
}

resource "oktawave_opn" "test-opn3" {
	name = "test-opn3"
}
`

	dataSourceConfig := `
data "oktawave_opns" "test-opns" {
	filter {
		key = "name"
		values = ["test-opn1", "test-opn3"]
	}
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
					resource.TestCheckResourceAttr("data.oktawave_opns.test-opns", "opns.#", "2"),
					resource.TestCheckResourceAttr("data.oktawave_opns.test-opns", "opns.0.name", "test-opn1"),
					resource.TestCheckResourceAttr("data.oktawave_opns.test-opns", "opns.1.name", "test-opn3"),
					resource.TestCheckResourceAttr("data.oktawave_opns.test-opns", "opns.1.private_ips.#", "0"),
					resource.TestCheckResourceAttrPair("data.oktawave_opns.test-opns", "opns.0.id", "oktawave_opn.test-opn1", "id"),
					resource.TestCheckResourceAttrPair("data.oktawave_opns.test-opns", "opns.1.id", "oktawave_opn.test-opn3", "id"),
				),
			},
		},
	})
}
