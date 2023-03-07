package oktawave

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/oktawave-code/odk"
)

func TestAccOktawaveInstance_Basic(t *testing.T) {
	var instance odk.Instance
	instanceConfig := `
resource "oktawave_ip" "test-ip1" {
	subregion_id = 1
}

resource "oktawave_instance" "test-instance1" {
	depends_on = [oktawave_ip.test-ip1]
	name = "test-instance1"
	subregion_id = 1
	system_disk_class_id = 48
	template_id = 1021
	type_id = 1047
	public_ips = [oktawave_ip.test-ip1.id]
}
`
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: instanceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckInstanceExists("oktawave_instance.test-instance1", &instance),
					// testAccCheckInstanceAttributes(&instance),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "name", "test-instance1"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "subregion_id", "1"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "template_id", "1021"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "type_id", "1047"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "system_disk_class_id", "48"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "init_script", ""),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "opn_mac.%", "0"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "opn_ids.#", "0"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "ssh_keys_ids.#", "0"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "public_ips.#", "1"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "private_ip_address", ""),
					resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "authorization_method_id"),
					resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "system_disk_size"),
					resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "system_disk_id"),
					resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "creation_date"),
					resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "creation_user_id"),
					resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "is_locked"),
					resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "locking_date"),
					resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "ip_address"),
					resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "mac_address"),
					resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "dns_address"),
					resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "status_id"),
					resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "system_category_id"),
					resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "autoscaling_type_id"),
					resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "vmware_tools_status_id"),
					resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "monit_status_id"),
					resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "template_type_id"),
					resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "payment_type_id"),
					resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "scsi_controller_type_id"),
					resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "total_disks_capacity"),
					resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "cpu_number"),
					resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "ram_mb"),
					resource.TestCheckNoResourceAttr("oktawave_instance.test-instance1", "opn_ids"),
					resource.TestCheckNoResourceAttr("oktawave_instance.test-instance1", "ssh_keys_ids"),
					resource.TestCheckNoResourceAttr("oktawave_instance.test-instance1", "health_check_id"),
					resource.TestCheckNoResourceAttr("oktawave_instance.test-instance1", "support_type_id"),
				),
			},
		},
	})
}

func TestAccOktawaveInstance_InitScriptAndPublicIpAndOpn(t *testing.T) {
	var instance odk.Instance
	instanceConfig := `
resource "oktawave_ip" "test-ip1" {
	subregion_id = 1
}

resource "oktawave_opn" "test-opn1" {
	name = "test-opn"
}

resource "oktawave_instance" "test-instance1" {
	depends_on = [oktawave_ip.test-ip1]
	name = "test-instance1"
	subregion_id = 1
	system_disk_class_id = 48
	template_id = 1021
	type_id = 1047
	init_script = "ZWNobyAiSGVsbG8gd29ybGQi"
	public_ips = [oktawave_ip.test-ip1.id]
	opn_ids = [oktawave_opn.test-opn1.id]
}
`
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: instanceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckInstanceExists("oktawave_instance.test-instance1", &instance),
					// testAccCheckInstanceAttributes(&instance),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "init_script", "ZWNobyAiSGVsbG8gd29ybGQi"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "public_ips.#", "1"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "opn_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair("oktawave_instance.test-instance1", "opn_ids.*", "oktawave_opn.test-opn1", "id"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "opn_mac.%", "1"),
					resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "ip_address"),
					resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "mac_address"),
				),
			},
		},
	})
}

func TestAccOktawaveInstance_SshKeys(t *testing.T) {
	var instance odk.Instance
	instanceConfig := `
resource "oktawave_ssh_key" "test-ssh_key1" {
	name = "test-key1"
	value = "ssh-rsa TEST testx"
}

resource "oktawave_ssh_key" "test-ssh_key2" {
	name = "test-key2"
	value = "ssh-rsa TEST testx"
}

resource "oktawave_ip" "test-ip1" {
	subregion_id = 1
}

resource "oktawave_instance" "test-instance1" {
	depends_on = [oktawave_ip.test-ip1]
	name = "test-instance1"
	subregion_id = 1
	system_disk_class_id = 48
	template_id = 1021
	type_id = 1047
	authorization_method_id = 1398
	ssh_keys_ids = [
		oktawave_ssh_key.test-ssh_key1.id,
		oktawave_ssh_key.test-ssh_key2.id
	]
	public_ips = [oktawave_ip.test-ip1.id]
}
`
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: instanceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckInstanceExists("oktawave_instance.test-instance1", &instance),
					// testAccCheckInstanceAttributes(&instance),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "ssh_keys_ids.#", "2"),
					resource.TestCheckTypeSetElemAttrPair("oktawave_instance.test-instance1", "ssh_keys_ids.*", "oktawave_ssh_key.test-ssh_key2", "id"),
					resource.TestCheckTypeSetElemAttrPair("oktawave_instance.test-instance1", "ssh_keys_ids.*", "oktawave_ssh_key.test-ssh_key1", "id"),
				),
			},
		},
	})
}

func TestAccOktawaveInstance_Update(t *testing.T) {
	var instance odk.Instance
	instanceConfig := `
resource "oktawave_ip" "test-ip1" {
	subregion_id = 1
}

resource "oktawave_instance" "test-instance1" {
	name = "test-instance1"
	subregion_id = 1
	system_disk_class_id = 48
	template_id = 1021
	type_id = 1047
	public_ips = [oktawave_ip.test-ip1.id]
}
`

	instanceUpdatedConfig := `
resource "oktawave_ip" "test-ip1" {
	subregion_id = 1
}

resource "oktawave_instance" "test-instance1" {
	name = "test-instance2"
	subregion_id = 1
	system_disk_class_id = 48
	template_id = 1021
	type_id = 289
	public_ips = [oktawave_ip.test-ip1.id]
}
`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: instanceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckInstanceExists("oktawave_instance.test-instance1", &instance),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "name", "test-instance1"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "type_id", "1047"),
				),
			},
			{
				Config: instanceUpdatedConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckInstanceExists("oktawave_instance.test-instance1", &instance),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "name", "test-instance2"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "type_id", "289"),
				),
			},
		},
	})
}

func TestAccOktawaveInstance_PublicInterfaces(t *testing.T) {
	var instance odk.Instance
	instanceConfig := `
resource "oktawave_ip" "test-ip1" {
	subregion_id = 1
}


resource "oktawave_instance" "test-instance1" {
	depends_on = [oktawave_ip.test-ip1]
	name = "test-instance1"
	subregion_id = 1
	system_disk_class_id = 48
	template_id = 1021
	type_id = 1047
	public_ips = [oktawave_ip.test-ip1.id]
}
`

	instanceAddInterfaceConfig := `
resource "oktawave_ip" "test-ip1" {
	subregion_id = 1
}

resource "oktawave_ip" "test-ip2" {
	subregion_id = 1
}

resource "oktawave_instance" "test-instance1" {
	depends_on = [oktawave_ip.test-ip1, oktawave_ip.test-ip2]
	name = "test-instance1"
	subregion_id = 1
	system_disk_class_id = 48
	template_id = 1021
	type_id = 1047
	public_ips = [oktawave_ip.test-ip1.id, oktawave_ip.test-ip2.id]
}
`

	instanceRemoveInterfaceConfig := `
resource "oktawave_ip" "test-ip1" {
	subregion_id = 1
}

resource "oktawave_ip" "test-ip2" {
	subregion_id = 1
}

resource "oktawave_instance" "test-instance1" {
	depends_on = [oktawave_ip.test-ip1, oktawave_ip.test-ip2]
	name = "test-instance1"
	subregion_id = 1
	system_disk_class_id = 48
	template_id = 1021
	type_id = 1047
	public_ips = [oktawave_ip.test-ip1.id]
}
`

	instanceSwitchToOpnConfig := `
resource "oktawave_ip" "test-ip1" {
	subregion_id = 1
}

resource "oktawave_ip" "test-ip2" {
	subregion_id = 1
}

resource "oktawave_opn" "test-opn" {
	name = "test-opn"
}

resource "oktawave_instance" "test-instance1" {
	depends_on = [oktawave_ip.test-ip1, oktawave_ip.test-ip2, oktawave_opn.test-opn]
	name = "test-instance1"
	subregion_id = 1
	system_disk_class_id = 48
	template_id = 1021
	type_id = 1047
	public_ips = []
	opn_ids = [
		oktawave_opn.test-opn.id
	]
}
`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: instanceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckInstanceExists("oktawave_instance.test-instance1", &instance),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "public_ips.#", "1"),
				),
			},
			{
				Config: instanceAddInterfaceConfig, // add public interface
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckInstanceExists("oktawave_instance.test-instance1", &instance),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "public_ips.#", "2"),
					resource.TestCheckTypeSetElemAttrPair("oktawave_instance.test-instance1", "public_ips.*", "oktawave_ip.test-ip1", "id"),
					resource.TestCheckTypeSetElemAttrPair("oktawave_instance.test-instance1", "public_ips.*", "oktawave_ip.test-ip2", "id"),
					// TODO check interface, mac address etc.
				),
			},
			{
				Config: instanceRemoveInterfaceConfig, // remove public interface
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckInstanceExists("oktawave_instance.test-instance1", &instance),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "public_ips.#", "1"),
					// resource.TestCheckResourceAttr("oktawave_ip.test-ip2", "instance_id", ""), FIXME this does not refresh
					resource.TestCheckResourceAttrPair("oktawave_instance.test-instance1", "public_ips.0", "oktawave_ip.test-ip1", "id"),
				),
			},
			{
				Config: instanceSwitchToOpnConfig, // switch to opn
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckInstanceExists("oktawave_instance.test-instance1", &instance),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "public_ips.#", "0"),
					// resource.TestCheckResourceAttr("oktawave_ip.test-ip1", "instance_id", ""), FIXME, instance_id doesn't refresh
					resource.TestCheckResourceAttr("oktawave_ip.test-ip2", "instance_id", ""),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "opn_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair("oktawave_instance.test-instance1", "opn_ids.*", "oktawave_opn.test-opn", "id"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "opn_mac.%", "1"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "ip_address", ""),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "mac_address", ""),
				),
			},
		},
	})
}

func TestAccOktawaveInstance_WithoutPublicInterfaceWithOpn(t *testing.T) {
	var instance odk.Instance
	instanceConfig := `
resource "oktawave_opn" "test-opn1" {
	name = "test-opn1"
}

resource "oktawave_opn" "test-opn2" {
	name = "test-opn2"
}

resource "oktawave_instance" "test-instance1" {
	depends_on = [oktawave_opn.test-opn1, oktawave_opn.test-opn2]
	name = "test-instance1"
	subregion_id = 1
	system_disk_class_id = 48
	template_id = 1021
	type_id = 1047
	without_public_ip = true
	opn_ids = [
		oktawave_opn.test-opn1.id
	]
}
`

	instanceConfigAddOpn := `
resource "oktawave_opn" "test-opn1" {
	name = "test-opn1"
}

resource "oktawave_opn" "test-opn2" {
	name = "test-opn2"
}

resource "oktawave_instance" "test-instance1" {
	depends_on = [oktawave_opn.test-opn1]
	name = "test-instance1"
	subregion_id = 1
	system_disk_class_id = 48
	template_id = 1021
	type_id = 1047
	without_public_ip = true
	opn_ids = [
		oktawave_opn.test-opn1.id,
		oktawave_opn.test-opn2.id
	]
}`

	instanceConfigRemoveOpn := `
resource "oktawave_opn" "test-opn1" {
	name = "test-opn1"
}

resource "oktawave_opn" "test-opn2" {
	name = "test-opn2"
}

resource "oktawave_instance" "test-instance1" {
	depends_on = [oktawave_opn.test-opn1]
	name = "test-instance1"
	subregion_id = 1
	system_disk_class_id = 48
	template_id = 1021
	type_id = 1047
	without_public_ip = true
	opn_ids = [
		oktawave_opn.test-opn1.id
	]
}`

	instanceSwitchToPublicIp := `
resource "oktawave_opn" "test-opn1" {
	name = "test-opn1"
}

resource "oktawave_opn" "test-opn2" {
	name = "test-opn2"
}

resource "oktawave_ip" "test-ip" {
	subregion_id = 1
}

resource "oktawave_instance" "test-instance1" {
	depends_on = [oktawave_opn.test-opn1, oktawave_ip.test-ip]
	name = "test-instance1"
	subregion_id = 1
	system_disk_class_id = 48
	template_id = 1021
	type_id = 1047
	without_public_ip = true
	opn_ids = []
	public_ips = [
		oktawave_ip.test-ip.id
	]
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: instanceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckInstanceExists("oktawave_instance.test-instance1", &instance),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "public_ips.#", "0"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "opn_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair("oktawave_instance.test-instance1", "opn_ids.*", "oktawave_opn.test-opn1", "id"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "opn_mac.%", "1"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "ip_address", ""),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "mac_address", ""),
					// resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "private_ip_address"),
				),
			},
			{
				Config: instanceConfigAddOpn,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckInstanceExists("oktawave_instance.test-instance1", &instance),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "public_ips.#", "0"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "opn_ids.#", "2"),
					resource.TestCheckTypeSetElemAttrPair("oktawave_instance.test-instance1", "opn_ids.*", "oktawave_opn.test-opn1", "id"),
					resource.TestCheckTypeSetElemAttrPair("oktawave_instance.test-instance1", "opn_ids.*", "oktawave_opn.test-opn2", "id"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "opn_mac.%", "2"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "ip_address", ""),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "mac_address", ""),
					// resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "private_ip_address"),
				),
			},
			{
				Config: instanceConfigRemoveOpn,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckInstanceExists("oktawave_instance.test-instance1", &instance),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "public_ips.#", "0"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "opn_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair("oktawave_instance.test-instance1", "opn_ids.*", "oktawave_opn.test-opn1", "id"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "opn_mac.%", "1"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "ip_address", ""),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "mac_address", ""),
					// resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "private_ip_address"),
				),
			},
			{
				Config: instanceSwitchToPublicIp,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckInstanceExists("oktawave_instance.test-instance1", &instance),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "public_ips.#", "1"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "opn_ids.#", "0"),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "opn_mac.%", "0"),
					resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "ip_address"),
					resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "mac_address"),
					// resource.TestCheckResourceAttrSet("oktawave_instance.test-instance1", "private_ip_address"),
				),
			},
		},
	})
}

func TestAccOktawaveInstance_Disks(t *testing.T) {
	var instance odk.Instance
	instanceConfig := `
resource "oktawave_disk" "test-disk1" {
	name = "disk1"
	tier_id = 48
	subregion_id = 1
	capacity = 5
}

resource "oktawave_disk" "test-disk2" {
	name = "disk2"
	tier_id = 48
	subregion_id = 1
	capacity = 5
}

resource "oktawave_ip" "test-ip1" {
	subregion_id = 1
}

resource "oktawave_instance" "test-instance1" {
	depends_on = [oktawave_disk.test-disk1, oktawave_disk.test-disk2, oktawave_ip.test-ip1]
	name = "test-instance1"
	subregion_id = 1
	system_disk_class_id = 48
	template_id = 1021
	type_id = 1047
	disks_ids = [
		oktawave_disk.test-disk1.id
	]
	public_ips = [oktawave_ip.test-ip1.id]
}
`

	instanceAttachDiskConfig := `
resource "oktawave_disk" "test-disk1" {
	name = "disk1"
	tier_id = 48
	subregion_id = 1
	capacity = 5
}

resource "oktawave_disk" "test-disk2" {
	name = "disk2"
	tier_id = 48
	subregion_id = 1
	capacity = 5
}

resource "oktawave_ip" "test-ip1" {
	subregion_id = 1
}

resource "oktawave_instance" "test-instance1" {
	depends_on = [oktawave_disk.test-disk1, oktawave_disk.test-disk2, oktawave_ip.test-ip1]
	name = "test-instance1"
	subregion_id = 1
	system_disk_class_id = 48
	template_id = 1021
	type_id = 1047
	disks_ids = [
		oktawave_disk.test-disk1.id,
		oktawave_disk.test-disk2.id
	]
	public_ips = [oktawave_ip.test-ip1.id]
}
`

	instanceDetachDiskConfig := `
resource "oktawave_disk" "test-disk1" {
	name = "disk1"
	tier_id = 48
	subregion_id = 1
	capacity = 5
}

resource "oktawave_disk" "test-disk2" {
	name = "disk2"
	tier_id = 48
	subregion_id = 1
	capacity = 5
}

resource "oktawave_ip" "test-ip1" {
	subregion_id = 1
}

resource "oktawave_instance" "test-instance1" {
	depends_on = [oktawave_disk.test-disk1, oktawave_disk.test-disk2, oktawave_ip.test-ip1]
	name = "test-instance1"
	subregion_id = 1
	system_disk_class_id = 48
	template_id = 1021
	type_id = 1047
	disks_ids = []
	public_ips = [oktawave_ip.test-ip1.id]
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: instanceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckInstanceExists("oktawave_instance.test-instance1", &instance),
					// testAccCheckInstanceAttributes(&instance),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "disks_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair("oktawave_instance.test-instance1", "disks_ids.*", "oktawave_disk.test-disk1", "id"),
				),
			},
			{
				Config: instanceAttachDiskConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckInstanceExists("oktawave_instance.test-instance1", &instance),
					// testAccCheckInstanceAttributes(&instance),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "disks_ids.#", "2"),
					resource.TestCheckTypeSetElemAttrPair("oktawave_instance.test-instance1", "disks_ids.*", "oktawave_disk.test-disk1", "id"),
					resource.TestCheckTypeSetElemAttrPair("oktawave_instance.test-instance1", "disks_ids.*", "oktawave_disk.test-disk2", "id"),
				),
			},
			{
				Config: instanceDetachDiskConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckInstanceExists("oktawave_instance.test-instance1", &instance),
					// testAccCheckInstanceAttributes(&instance),
					resource.TestCheckResourceAttr("oktawave_instance.test-instance1", "disks_ids.#", "0"),
				),
			},
		},
	})
}

// Test instance behaviour after template creation
func TestAccOktawaveInstance_ConvertToTemplate(t *testing.T) {
	var instance odk.Instance
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
`

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

	instanceUpdateConfig := `
resource "oktawave_instance" "test-instance" {
	name = "test-instance"
	subregion_id = 1
	system_disk_class_id = 48
	template_id = 1021
	type_id = 1049
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

	instanceRemovedConfig := func(instanceId int32) string {
		return fmt.Sprintf(`
resource "oktawave_template" "template" {
	instance_id = %d
	name = "test-template"
	description = "test-template"
	version = "0.1"
	system_category_id = 1277
	default_type_id = 1047
	minimum_type_id = 1047
}`, instanceId)
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDatasourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckInstanceExists("oktawave_instance.test-instance", &instance),
				),
			},
			{
				Config: templateConfig,
			},
			{
				Config: instanceUpdateConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckInstanceDoesntExist("oktawave_instance.test-instance"),
					resource.TestCheckResourceAttrPair("oktawave_instance.test-instance", "converted_to_template_id", "oktawave_template.template", "id"),
				),
			},
			{
				// check if instance definition can be removed safely
				Config: instanceRemovedConfig(instance.Id),
			},
		},
	})
}

func testAccCheckInstanceExists(resourceName string, instance *odk.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Instance ID is not set")
		}

		client := testAccProvider.Meta().(*ClientConfig).odkClient
		auth := testAccProvider.Meta().(*ClientConfig).odkAuth

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Failed to parse resource id %s.", rs.Primary.ID)
		}

		object, _, err := client.OCIApi.InstancesGet_2(*auth, int32(id), nil)
		if err != nil {
			return err
		}

		if strconv.Itoa(int(object.Id)) != rs.Primary.ID {
			return fmt.Errorf("Instance not found")
		}

		*instance = object
		return nil
	}
}

func testAccCheckInstanceDoesntExist(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Instance ID is not set")
		}

		client := testAccProvider.Meta().(*ClientConfig).odkClient
		auth := testAccProvider.Meta().(*ClientConfig).odkAuth

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Failed to parse resource id %s.", rs.Primary.ID)
		}

		_, resp, err := client.OCIApi.InstancesGet_2(*auth, int32(id), nil)
		if err != nil {
			if resp.StatusCode == 404 {
				return nil
			}

			return fmt.Errorf("Failed to call for instance. Caused by: %s.", err)
		}

		return fmt.Errorf("Instance exists.")
	}
}
