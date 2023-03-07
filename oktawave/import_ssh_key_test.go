package oktawave

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOktawaveSshKey_importBasic(t *testing.T) {
	resourceName := "oktawave_ssh_key.test-ssh_key"
	ssh_keyConfig := `
resource "oktawave_ssh_key" "test-ssh_key" {
	name = "test-key"
	value = "ssh-rsa TEST testx"
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSshKeyDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: ssh_keyConfig,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
