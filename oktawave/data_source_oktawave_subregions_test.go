package oktawave

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestOktawave_DataSource_Subregions(t *testing.T) {
	dataSourceConfig := `
data "oktawave_subregions" "subregions" {
	filter {
		key = "name"
		values = ["PL-001", "PL-002"]
	}
}
`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.oktawave_subregions.subregions", "items.#", "2"),
					resource.TestCheckResourceAttr("data.oktawave_subregions.subregions", "items.0.name", "PL-001"),
					resource.TestCheckResourceAttr("data.oktawave_subregions.subregions", "items.1.name", "PL-002"),
				),
			},
		},
	})
}
