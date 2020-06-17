data "oktawave_oci" "oci" {
  instance_name = var.instance_name
  subregion_id = var.subregion_id
  template_id = var.template_id
  type_id = var.type_id
}
