package oktawave

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/oktawave-code/odk"
)

func TestAccOktawaveSshKey_Basic(t *testing.T) {
	var ssh_key odk.SshKey
	instanceConfig := `
resource "oktawave_ssh_key" "test-ssh_key" {
	name = "test-key"
	value = "ssh-rsa TEST testx"
}
`
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSshKeyDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: instanceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSshKeyExists("oktawave_ssh_key.test-ssh_key", &ssh_key),
					// testAccCheckSsh_keyAttributes(&ssh_key),
					resource.TestCheckResourceAttr("oktawave_ssh_key.test-ssh_key", "name", "test-key"),
					resource.TestCheckResourceAttr("oktawave_ssh_key.test-ssh_key", "value", "ssh-rsa TEST testx"),
					resource.TestCheckResourceAttrSet("oktawave_ssh_key.test-ssh_key", "name"),
					resource.TestCheckResourceAttrSet("oktawave_ssh_key.test-ssh_key", "value"),
					resource.TestCheckResourceAttrSet("oktawave_ssh_key.test-ssh_key", "owner_user_id"),
					resource.TestCheckResourceAttrSet("oktawave_ssh_key.test-ssh_key", "creation_date"),
				),
			},
		},
	})
}

func testAccCheckSshKeyExists(resourceName string, ssh_key *odk.SshKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Ssh_key ID is not set")
		}

		client := testAccProvider.Meta().(*ClientConfig).odkClient
		auth := testAccProvider.Meta().(*ClientConfig).odkAuth

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Failed to parse resource id %s.", rs.Primary.ID)
		}

		object, _, err := client.AccountApi.AccountGetSshKey(*auth, int32(id), nil)
		if err != nil {
			return err
		}

		if strconv.Itoa(int(object.Id)) != rs.Primary.ID {
			return fmt.Errorf("Ssh_key not found")
		}

		*ssh_key = object
		return nil
	}
}

// func testAccCheckssh_keyAttributes(ssh_key *odk.Ssh_key) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 	  if *widget.active != true {
// 		return fmt.Errorf("widget is not active")
// 	  }

// 	  return nil
// 	}
// }
