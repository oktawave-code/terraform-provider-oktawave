data "oktawave_ssh_key" "key" {
  name = "key1"
}

output "test_key" {
  value = data.oktawave_ssh_key.key
}
