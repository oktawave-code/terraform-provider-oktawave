# Terraform Provider

Oktawave terraform provider source code

# How to build

```shell
go build -o terraform-provider-oktawave
```

# How to test

With acceptance tests (uses real API to create resources):
```shell
OKTAWAVE_ACCESS_TOKEN={access_token} TF_ACC=1 go test -v ./...
```

With unit tests only (fast):
```shell
go test -v ./...
```

# How to use

```shell
OKTAWAVE_ACCESS_TOKEN={access_token} OKTAWAVE_DC={DC1|DC2} terraform {terraform commands here}
```

# Env vars supported by this plugin

- OKTAWAVE_ACCESS_TOKEN (access_token) - needed to authorize to Oktawave apis
- OKTAWAVE_DC (dc) - data center selector: "DC1" or "DC2"
- OKTAWAVE_ODK_API_URL (odk_api_url) - manual api url override
- OKTAWAVE_ODK_API_SKIP_TLS (odk_api_skip_tls) - manual disabling of certificate check
- OKTAWAVE_OKS_API_URL (oks_api_url) - manual api url override
- OKTAWAVE_OKS_API_SKIP_TLS (oks_api_skip_tls) - manual disabling of certificate check

# You can generate access_token using curl:
```shell
curl -k -X POST -d "grant_type=password&username=youremail&password=yourpassword&scope=oktawave.api" -u "client_id:client_secret" 'https://id.oktawave.com/core/connect/token'
```
