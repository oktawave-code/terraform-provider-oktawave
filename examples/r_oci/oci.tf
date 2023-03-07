resource "oktawave_instance" "instance" {
  name = "test_instance"
  subregion_id = 1014
  system_disk_class_id = 49
  template_id = 3
  type_id = 1047
}

output "test" {
  value = resource.oktawave_instance.instance
}
