terraform {
    required_providers {
        oktawave = {
            version = "~> 1.0.0"
            source  = "oktawave.com/iaac/oktawave"
        }
    }
}

provider "oktawave" {
  access_token="your-access-token"
  api_url = "https://api.oktawave.com/beta/"
}

resource "oktawave_opn" "my_opn"{
	# Required: true
    # Type: string
    # Available values: string of available length that represent name of private network
	opn_name="test_opn"
    
}

    # ADDITIONAL INFO/ INFORMACJE DODATKOWE
    
    
    # COMPUTED VALUES/ WARTOŚCI WYLICZALNE AUTOMATYCZNIE
    
	# creation_user_id
	# Computed: true
	# Type: int
	# Available values: computed attribute is not available to set manually
	
	# creation_date
	# Computed: true
	# Type: string
	# Available values: computed attribute is not available to set manually
	
	# last_change_date
	# Computed: true
	# Type: string
	# Available values: computed attribute is not available to set manually
