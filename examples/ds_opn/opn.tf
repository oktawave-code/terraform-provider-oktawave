data "oktawave_opn" "opn" {
  name = "autotest-opn"
}

output "test" {
  value = data.oktawave_opn.opn
}
