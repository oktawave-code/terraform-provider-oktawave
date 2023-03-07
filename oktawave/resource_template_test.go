package oktawave

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/oktawave-code/odk"
)

func TestAccOktawaveTemplate_Basic(t *testing.T) {
	var template odk.Template
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
}
`

	templateUpdated := `
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
	name = "test-template1"
	description = "test-template1"
	version = "0.1"
	system_category_id = 1277
	default_type_id = 1047
	minimum_type_id = 1047
}
`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTemplateDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: templateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTemplateExists("oktawave_template.test-template", &template),
					// testAccCheckTemplateAttributes(&template),
					resource.TestCheckResourceAttr("oktawave_template.test-template", "name", "test-template"),
					resource.TestCheckResourceAttr("oktawave_template.test-template", "description", "test-template"),
					resource.TestCheckResourceAttrSet("oktawave_template.test-template", "instance_id"),
					resource.TestCheckResourceAttrSet("oktawave_template.test-template", "name"),
					resource.TestCheckResourceAttrSet("oktawave_template.test-template", "description"),
					resource.TestCheckResourceAttrSet("oktawave_template.test-template", "version"),
					resource.TestCheckResourceAttrSet("oktawave_template.test-template", "system_category_id"),
					resource.TestCheckResourceAttrSet("oktawave_template.test-template", "default_type_id"),
					resource.TestCheckResourceAttrSet("oktawave_template.test-template", "minimum_type_id"),
					resource.TestCheckResourceAttrSet("oktawave_template.test-template", "creation_date"),
					resource.TestCheckResourceAttrSet("oktawave_template.test-template", "last_change_date"),
					resource.TestCheckResourceAttrSet("oktawave_template.test-template", "creation_user_id"),
					resource.TestCheckResourceAttrSet("oktawave_template.test-template", "ethernet_controllers_number"),
					resource.TestCheckResourceAttrSet("oktawave_template.test-template", "ethernet_controllers_type_id"),
					resource.TestCheckResourceAttrSet("oktawave_template.test-template", "publication_status_id"),
					resource.TestCheckResourceAttr("oktawave_template.test-template", "disks.%", "1"),
					resource.TestCheckResourceAttr("oktawave_template.test-template", "software.%", "0"),
					resource.TestCheckResourceAttrSet("oktawave_template.test-template", "template_type_id"),
				),
			},
			{
				Config: templateUpdated,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTemplateExists("oktawave_template.test-template", &template),
					// testAccCheckTemplateAttributes(&template),
					resource.TestCheckResourceAttr("oktawave_template.test-template", "name", "test-template1"),
					resource.TestCheckResourceAttr("oktawave_template.test-template", "description", "test-template1"),
				),
			},
		},
	})
}

func testAccCheckTemplateExists(resourceName string, template *odk.Template) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Template ID is not set")
		}

		client := testAccProvider.Meta().(*ClientConfig).odkClient
		auth := testAccProvider.Meta().(*ClientConfig).odkAuth

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Failed to parse resource id %s.", rs.Primary.ID)
		}

		object, _, err := client.OCITemplatesApi.TemplatesGet_1(*auth, int32(id), nil)
		if err != nil {
			return err
		}

		if strconv.Itoa(int(object.Id)) != rs.Primary.ID {
			return fmt.Errorf("Template not found")
		}

		*template = object
		return nil
	}
}

// func testAccChecktemplateAttributes(template *odk.Template) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 	  if *widget.active != true {
// 		return fmt.Errorf("widget is not active")
// 	  }

// 	  return nil
// 	}
// }
