package oktawave

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestOktawave_DataSource_InstanceTypes(t *testing.T) {
	dataSourceConfig := `
data "oktawave_instance_types" "types" {
	filter {
		key = "name"
		values = ["v1.highcpu-16.12", "v1.highcpu-16.8"]
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
					resource.TestCheckResourceAttr("data.oktawave_instance_types.types", "items.#", "2"),
					resource.TestCheckResourceAttr("data.oktawave_instance_types.types", "items.0.name", "v1.highcpu-16.12"),
				),
			},
		},
	})
}
