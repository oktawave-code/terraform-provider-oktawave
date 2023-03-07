resource "oktawave_disk" "disk" {
  name = "test_disk"
  tier_id = 49
  subregion_id = 1014
}

output "test" {
  value = resource.oktawave_disk.disk
}
