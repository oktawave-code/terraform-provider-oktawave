resource "oktawave_sshKey" "key" {
  ssh_key_name = "test"
  ssh_key_value = "ssh-rsa TEST test"
}

data "oktawave_sshKey" "key" {
  ssh_key_name = "key1"
}

output "test_ssh" {
  value = data.oktawave_sshKey.key
}
