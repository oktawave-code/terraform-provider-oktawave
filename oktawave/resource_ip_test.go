package oktawave

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/oktawave-code/odk"
)

func TestAccOktawaveIp_Basic(t *testing.T) {
	var ip odk.Ip
	instanceConfig := `
resource "oktawave_ip" "test-ip" {
	subregion_id = 1
}
`
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIpDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: instanceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIpExists("oktawave_ip.test-ip", &ip),
					// testAccCheckIpAttributes(&ip),
					resource.TestCheckResourceAttr("oktawave_ip.test-ip", "subregion_id", "1"),
					resource.TestCheckResourceAttrSet("oktawave_ip.test-ip", "subregion_id"),
					resource.TestCheckResourceAttr("oktawave_ip.test-ip", "comment", ""),
					resource.TestCheckResourceAttrSet("oktawave_ip.test-ip", "type_id"),
					resource.TestCheckResourceAttrSet("oktawave_ip.test-ip", "address"),
					resource.TestCheckResourceAttrSet("oktawave_ip.test-ip", "gateway"),
					resource.TestCheckResourceAttrSet("oktawave_ip.test-ip", "netmask"),
					resource.TestCheckResourceAttr("oktawave_ip.test-ip", "instance_id", ""),
					resource.TestCheckResourceAttrSet("oktawave_ip.test-ip", "interface_id"),
					resource.TestCheckResourceAttrSet("oktawave_ip.test-ip", "dhcp_branch"),
					resource.TestCheckResourceAttrSet("oktawave_ip.test-ip", "mode_id"),
					resource.TestCheckResourceAttrSet("oktawave_ip.test-ip", "creation_user_id"),
					resource.TestCheckResourceAttr("oktawave_ip.test-ip", "rev_dns", ""),
					resource.TestCheckResourceAttr("oktawave_ip.test-ip", "rev_dns_v6", ""),
					resource.TestCheckResourceAttr("oktawave_ip.test-ip", "address_v6", ""),
					resource.TestCheckResourceAttr("oktawave_ip.test-ip", "mac_address", ""),
					resource.TestCheckResourceAttr("oktawave_ip.test-ip", "dns_prefix", ""),
				),
			},
		},
	})
}

func TestAccOktawaveIp_Update(t *testing.T) {
	var ip odk.Ip

	ipConfig := `
resource "oktawave_ip" "test-ip" {
	subregion_id = 1
	comment = "test-comment"
}
`

	ipConfigUpdated := `
resource "oktawave_ip" "test-ip" {
	subregion_id = 4
	comment = "test-comment2"
}
`
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIpDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: ipConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIpExists("oktawave_ip.test-ip", &ip),
					// testAccCheckIpAttributes(&ip),
					resource.TestCheckResourceAttr("oktawave_ip.test-ip", "subregion_id", "1"),
					resource.TestCheckResourceAttr("oktawave_ip.test-ip", "comment", "test-comment"),
				),
			},
			{
				Config: ipConfigUpdated,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIpExists("oktawave_ip.test-ip", &ip),
					// testAccCheckIpAttributes(&ip),
					resource.TestCheckResourceAttr("oktawave_ip.test-ip", "subregion_id", "4"),
					resource.TestCheckResourceAttr("oktawave_ip.test-ip", "comment", "test-comment2"),
				),
			},
		},
	})
}

func TestAccOktawaveIp_Attached(t *testing.T) {
	var ip odk.Ip
	instanceConfig := `
resource "oktawave_ip" "test-ip" {
	subregion_id = 1
}

resource "oktawave_instance" "test-instance" {
	depends_on = [oktawave_ip.test-ip]
	name = "test-instance1"
	subregion_id = 1
	system_disk_class_id = 48
	template_id = 1021
	type_id = 1047
	public_ips = [oktawave_ip.test-ip.id]
}
`
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIpDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: instanceConfig, // FIXME - this needs to be applied twice to force IP data synchronization
			},
			{
				Config: instanceConfig, // FIXME - this needs to be applied twice to force IP data synchronization
			},
			{
				Config: instanceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIpExists("oktawave_ip.test-ip", &ip),
					// testAccCheckIpAttributes(&ip),
					resource.TestCheckResourceAttrSet("oktawave_ip.test-ip", "rev_dns"),
					resource.TestCheckResourceAttrSet("oktawave_ip.test-ip", "rev_dns_v6"),
					resource.TestCheckResourceAttrSet("oktawave_ip.test-ip", "address_v6"),
					resource.TestCheckResourceAttrSet("oktawave_ip.test-ip", "mac_address"),
					resource.TestCheckResourceAttrSet("oktawave_ip.test-ip", "dns_prefix"),
					resource.TestCheckResourceAttrPair("oktawave_ip.test-ip", "instance_id", "oktawave_instance.test-instance", "id"),
				),
			},
		},
	})
}

func testAccCheckIpExists(resourceName string, ip *odk.Ip) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Ip ID is not set")
		}

		client := testAccProvider.Meta().(*ClientConfig).odkClient
		auth := testAccProvider.Meta().(*ClientConfig).odkAuth

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Failed to parse resource id %s.", rs.Primary.ID)
		}

		object, _, err := client.OCIInterfacesApi.InstancesGetInstanceIp(*auth, int32(id), nil)
		if err != nil {
			return err
		}

		if strconv.Itoa(int(object.Id)) != rs.Primary.ID {
			return fmt.Errorf("Ip not found")
		}

		*ip = object
		return nil
	}
}

// func testAccCheckipAttributes(ip *odk.Ip) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 	  if *widget.active != true {
// 		return fmt.Errorf("widget is not active")
// 	  }

// 	  return nil
// 	}
// }
