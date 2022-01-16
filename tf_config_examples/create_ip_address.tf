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



resource "oktawave_ip_address" "my_ip"{
	# Required: true
    # Type: int
	# ForceNew: true
    # Available values: 1(PL001), 4(PL002), 5(PL003), 6(PL004), 7(PL005)
	# Comment: ip subregion == instance subregion
	subregion_id=4
	
	# Optional: true
    # Type: string
    # Available values: string of available length
	comment="example ip_updated"
}



