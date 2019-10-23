# Terraform Oktawave Provider

Terraform provider for Oktawave Cloud.

# How to install and use:
1. Compile project using 
```bash
go build -o terraform-provider-oktawave.exe
```

 
 
2. Create config file witf .tf file extension.

3. Locate files terraform-provider-oktawave.exe, terraform.exe and you cfg.tf file in the same directory.

4. Run terraform via command line using

```bash
terraform init & terraform apply
```

5. To delete infrastructure use

```bash
terraform destroy
```

6. To update version of provider to the new one drop new terraform-provider-oktawave.exe file to directory with terraform.exe and use

```bash
terraform init
```

# List of supported resources:
	1. oktawave_oci
    
    2. oktawave_ovs
    
    3. oktawave_group
    
    4. oktawave_ip_address
    
    5. oktawave_load_balancer
    
    6. oktawave_opn
    
    7. oktawave_ssheky

# You can generate access_token using curl:
	curl -k -X POST -d "grant_type=password&username=youremail&password=yourpassword&scope=oktawave.api" -u "client_id:client_secret" 'https://id.oktawave.com/core/connect/token'

# Where can i find examples?
	You can find examples at projects directory /tf_config_examples

# Where can i find documentation for provider?
	You can find documentation for particular resource of provider in example config file at projects directory /tf_config_examples


