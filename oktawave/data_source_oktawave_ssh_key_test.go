package oktawave

import (
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestOktawave_DataSource_SshKey(t *testing.T) {
	resourcesConfig := `
resource "oktawave_ssh_key" "key1" {
	name = "test-key"
	value = "ssh-rsa TEST testx"
}
	`

	dataSourceConfig := `
data "oktawave_ssh_key" "key" {
	id = oktawave_ssh_key.key1.id
}
`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSshKeyDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + dataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.oktawave_ssh_key.key", "name", "test-key"),
					resource.TestCheckResourceAttrPair("data.oktawave_ssh_key.key", "id", "oktawave_ssh_key.key1", "id"),
				),
			},
		},
	})
}

func testAccCheckSshKeyDatasourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ClientConfig).odkClient
	auth := testAccProvider.Meta().(*ClientConfig).odkAuth

	params := map[string]interface{}{
		"pageSize": int32(math.MaxInt16),
	}
	keys, _, err := client.AccountApi.AccountGetSshKeys(*auth, params)
	if err != nil {
		return fmt.Errorf("Get ssh keys request failed. Caused by: %s.", err)
	}

	errors := []interface{}{}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "oktawave_ssh_key" {
			continue
		}

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Failed to parse resource id %s.", rs.Primary.ID)
		}

		for _, key := range keys.Items {
			if int64(key.Id) == id {
				errors = append(errors, fmt.Sprintf("Ssh key with id %d not destroyed correctly.", key.Id))
				break
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("Test failed to cleanup. %s", errors)
	}

	return nil
}
