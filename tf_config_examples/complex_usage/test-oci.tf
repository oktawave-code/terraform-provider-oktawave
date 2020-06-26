# Description in variables.tf
module "test-oci" {
  providers = { oktawave = oktawave.ALIAS_NAME }
  source = "./oktawave-oci"
  instance_name = "test-oci"
  authorization_method_id = 1398
  ssh_keys_ids = [2025]
  disk_class = 48
  init_disk_size = 20
  subregion_id = 7
  type_id = 1423
  template_id = 1018
  opn_ids = [3670]
  ovs_disk_name = "test"
  ovs_space_capacity = 5
  ovs_tier_id = 48
  init_script = "ZmlsZSB7ICcvcm9vdC90ZXN0b3d5X3BsaWsudHh0JzoKICBjb250ZW50ID0+ICd0ZXN0IHRlc3QgdGVzdCcsCn0K"
}
