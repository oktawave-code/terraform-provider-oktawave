package oktawave

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestOktawave_DataSource_Template(t *testing.T) {
	resourcesConfig := `
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

resource "oktawave_template" "template" {
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

	dataSourceConfig := `
data "oktawave_template" "template" {
	id = oktawave_template.template.id
}
`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTemplateDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.oktawave_template.template", "name", "test-template"),
				),
			},
		},
	})
}

func testAccCheckTemplateDatasourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ClientConfig).odkClient
	auth := testAccProvider.Meta().(*ClientConfig).odkAuth

	params := map[string]interface{}{
		"pageSize": int32(math.MaxInt16),
	}
	templates, _, err := client.OCITemplatesApi.TemplatesGet(*auth, params)
	if err != nil {
		return fmt.Errorf("Get templates request failed. Caused by: %s.", err)
	}

	errors := []interface{}{}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oktawave_template" {
			continue
		}

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Failed to parse resource id %s.", rs.Primary.ID)
		}

		for _, template := range templates.Items {
			if int64(template.Id) == id {
				errors = append(errors, fmt.Sprintf("Template with id %d not destroyed correctly.", template.Id))
				break
			}
		}
	}

	err = testAccCheckInstanceDatasourceDestroy(s)
	if err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("Test failed to cleanup. %s", errors)
	}

	return nil
}
