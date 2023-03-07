package oktawave

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestOktawave_DataSource_SshKeys(t *testing.T) {
	resourcesConfig := `
resource "oktawave_ssh_key" "key1" {
	name = "test-key1"
	value = "ssh-rsa TEST testx"
}

resource "oktawave_ssh_key" "key2" {
	name = "test-key2"
	value = "ssh-rsa TEST testx"
}


resource "oktawave_ssh_key" "key3" {
	name = "test-key3"
	value = "ssh-rsa TEST testx"
}`

	dataSourceConfig := `
data "oktawave_ssh_keys" "keys" {
	filter {
		key = "name"
		values = ["test-key1", "test-key3"]
	}
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
					resource.TestCheckResourceAttr("data.oktawave_ssh_keys.keys", "ssh_keys.#", "2"),
					resource.TestCheckResourceAttr("data.oktawave_ssh_keys.keys", "ssh_keys.0.name", "test-key1"),
					resource.TestCheckResourceAttr("data.oktawave_ssh_keys.keys", "ssh_keys.1.name", "test-key3"),
					resource.TestCheckResourceAttrPair("data.oktawave_ssh_keys.keys", "ssh_keys.0.id", "oktawave_ssh_key.key1", "id"),
					resource.TestCheckResourceAttrPair("data.oktawave_ssh_keys.keys", "ssh_keys.1.id", "oktawave_ssh_key.key3", "id"),
				),
			},
		},
	})
}
