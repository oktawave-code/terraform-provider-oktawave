data "oktawave_disk" "disk" {
  name = "autotest-disk"
}

output "test" {
  value = data.oktawave_disk.disk
}
