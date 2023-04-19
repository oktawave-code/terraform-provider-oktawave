package oktawave

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestOktawave_DataSource_Instances(t *testing.T) {
	resourcesConfig := `
resource "oktawave_ip" "test-ip1" {
	subregion_id = 1
}

resource "oktawave_instance" "test-instance1" {
	public_ips = [oktawave_ip.test-ip1.id]
	depends_on = [oktawave_ip.test-ip1]
	name = "test-instance1"
	subregion_id = 1
	system_disk_class_id = 48
	template_id = 1021
	type_id = 1047
}


resource "oktawave_ip" "test-ip2" {
	subregion_id = 1
}

resource "oktawave_instance" "test-instance2" {
	public_ips = [oktawave_ip.test-ip2.id]
	depends_on = [oktawave_ip.test-ip2]
	name = "test-instance2"
	subregion_id = 1
	system_disk_class_id = 48
	template_id = 1021
	type_id = 1047
}


resource "oktawave_ip" "test-ip3" {
	subregion_id = 1
}

resource "oktawave_instance" "test-instance3" {
	public_ips = [oktawave_ip.test-ip3.id]
	depends_on = [oktawave_ip.test-ip3]
	name = "test-instance3"
	subregion_id = 1
	system_disk_class_id = 48
	template_id = 1021
	type_id = 1047
}`

	dataSourceConfig := `
data "oktawave_instances" "instances" {
	filter {
		key = "name"
		values = ["test-instance1", "test-instance3"]
	}
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.oktawave_instances.instances", "items.#", "2"),
					resource.TestCheckResourceAttr("data.oktawave_instances.instances", "items.0.subregion_id", "1"),
					resource.TestCheckResourceAttr("data.oktawave_instances.instances", "items.1.subregion_id", "1"),
					resource.TestCheckResourceAttrPair("data.oktawave_instances.instances", "items.0.id", "oktawave_instance.test-instance1", "id"),
					resource.TestCheckResourceAttrPair("data.oktawave_instances.instances", "items.1.id", "oktawave_instance.test-instance3", "id"),
				),
			},
		},
	})
}
