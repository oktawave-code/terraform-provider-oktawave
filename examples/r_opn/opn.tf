resource "oktawave_opn" "opn" {
  name = "test_opn"
}

output "test" {
  value = resource.oktawave_opn.opn
}
