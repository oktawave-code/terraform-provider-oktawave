# Required empty provider while using aliases for multiple accounts
provider "oktawave" {
    access_token = ""
}

# Account 1
provider "oktawave" {
  alias = "ALIAS_NAME"
  access_token = "TOKEN"
  api_url = "https://api.oktawave.com/beta"
}

# Account 2
provider "oktawave" {
  alias = "ALIAS_NAME2"
  access_token = "TOKEN"
  api_url = "https://api.oktawave.com/beta"
}
