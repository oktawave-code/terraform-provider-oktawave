data "oktawave_template" "template" {
  name = "autotest-template"
}

output "test" {
  value = data.oktawave_template.template
}
