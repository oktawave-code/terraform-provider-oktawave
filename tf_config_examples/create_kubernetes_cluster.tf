terraform {
    required_providers {
        oktawave = {
            version = "~> 1.0.0"
            source  = "oktawave.com/iaac/oktawave"
        }
    }
}

provider "oktawave"{
	access_token="your-access-token"
	api_url="https://k44s-api.i.k44s.oktawave.com"
}

resource "oktawave_kubernetes_cluster" "my_cluster" {
	# Required: true
	# Type: string
	# Available values: any string with minimal length=6 and maximum length=10
    name="tfclusr"
	
	# Required: true
	# Type: string
	# Available values: 1.21, 1.18
    version="1.21"
}


	# ADDITIONAL INFO/INFORMACJE DODATKOWE
	
	
	# COMPUTED VALUES/WARTOŚCI WYLICZALNE AUTOMATYCZNIE
	
	
	# Computed: true
	# Type: bool
	# Available values: computed attribute. Not available to set manually
	# is_running = true/false
