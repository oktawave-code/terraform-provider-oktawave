package oktawave

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktawaveInstance_importBasic(t *testing.T) {
	resourceName := "oktawave_instance.test-instance"
	instanceConfig := `
resource "oktawave_ip" "test-ip1" {
	subregion_id = 1
}

resource "oktawave_instance" "test-instance" {
	depends_on = [oktawave_ip.test-ip1]
	name = "test-instance1"
	subregion_id = 1
	system_disk_class_id = 48
	template_id = 1021
	type_id = 1047
	public_ips = [oktawave_ip.test-ip1.id]
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: instanceConfig,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"authorization_method_id"},
			},
		},
	})
}
