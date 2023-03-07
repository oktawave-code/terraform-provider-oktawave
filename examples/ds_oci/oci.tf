data "oktawave_instance" "instance" {
  name = "autotest-instance"
}

output "test" {
  value = data.oktawave_instance.instance
}
