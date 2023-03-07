data "oktawave_group" "group" {
  name = "autotest-group"
}

output "test" {
  value = data.oktawave_group.group
}
