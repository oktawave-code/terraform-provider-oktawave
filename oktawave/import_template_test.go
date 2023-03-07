package oktawave

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktawaveTemplate_importBasic(t *testing.T) {
	resourceName := "oktawave_template.test-template"
	templateConfig := `
resource "oktawave_ip" "test-ip" {
	subregion_id = 1
}

resource "oktawave_instance" "test-instance" {
	depends_on = [oktawave_ip.test-ip]
	name = "test-instance"
	subregion_id = 1
	system_disk_class_id = 48
	template_id = 1021
	type_id = 1047
	public_ips = [oktawave_ip.test-ip.id]
}

resource "oktawave_template" "test-template" {
	depends_on = [oktawave_instance.test-instance]
	instance_id = oktawave_instance.test-instance.id
	name = "test-template"
	description = "test-template"
	version = "0.1"
	system_category_id = 1277
	default_type_id = 1047
	minimum_type_id = 1047
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTemplateDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: templateConfig,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id"},
			},
		},
	})
}
