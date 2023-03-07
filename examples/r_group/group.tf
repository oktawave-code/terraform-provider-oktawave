resource "oktawave_group" "group" {
  name = "test_group"
}

output "test" {
  value = resource.oktawave_group.group
}
