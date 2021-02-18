resource "oktawave_oci" "OCI" {
  instance_name = var.instance_name
  authorization_method_id = var.authorization_method_id
  ssh_keys_ids = var.ssh_keys_ids
  disk_class = var.disk_class
  init_disk_size = var.init_disk_size
  ip_address_ids = var.ip_address_ids
  subregion_id = var.subregion_id
  type_id = var.type_id
  template_id = var.template_id
  instances_count = var.instances_count
  isfreemium = var.isfreemium
  opn_ids = var.opn_ids
  init_script = var.init_script
  # path to puppet manifest file
  # init_script = filebase64("${path.module}/${var.init_script_file}")
  without_public_ip = var.without_public_ip
}

resource "oktawave_ovs" "OVS" {
  disk_name = var.ovs_disk_name
  space_capacity = var.ovs_space_capacity
  tier_id = var.ovs_tier_id
  subregion_id = var.subregion_id
  connections_with_instanceids = [data.oktawave_oci.oci.id]
  depends_on = [oktawave_oci.OCI]
}
