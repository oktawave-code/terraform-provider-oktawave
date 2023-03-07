data "oktawave_oks_node" "oks_node" {
  name = "k44sdev-autotest-07eajln4n-tcdl2lnn"
}

output "test" {
  value = data.oktawave_oks_node.oks_node
}
