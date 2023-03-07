data "oktawave_load_balancer" "lb" {
  group_name = "autotest-group"
}

output "test" {
  value = data.oktawave_load_balancer.lb
}
