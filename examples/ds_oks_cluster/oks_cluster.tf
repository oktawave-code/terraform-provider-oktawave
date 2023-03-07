data "oktawave_oks_cluster" "oks_cluster" {
  name = "autotest-07eajln4n"
}

output "test" {
  value = data.oktawave_oks_cluster.oks_cluster
}
