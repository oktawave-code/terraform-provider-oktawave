resource "oktawave_load_balancer" "lb" {
  group_name = "test_opn"
}

output "test" {
  value = resource.oktawave_load_balancer.lb
}
