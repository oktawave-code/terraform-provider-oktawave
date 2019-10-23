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
	comment="example ip_updated
}







	# ADDITIONAL INFO/INFORMACJE DODATKOWE

    # COMPUTED VALUES/WARTOŚCI WYLICZALNE AUTOMATYCZNIE
	
    #restore_rev_dns
	# Computed: true
    # Type: bool
    # Available values: computed attribute. Not available to set manually
	
	#restore_rev_dns_v6
	# Computed: true
    # Type: bool
    # Available values: computed attribute. Not available to set manually

    # rev_dns
	# Computed: true
	# Type: string
	# Available values: computed attribute. Not available to set manually

	# address
	# Computed: true
	# Type: string
	# Available values: computed attribute. Not available to set manually
    
	# address_v6
	# Computed: true
	# Type: string
	# Available values: computed attribute. Not available to set manually    
    
    #restore_rev_dns
	# Computed: true
    # Type: bool
    # Available values: computed attribute. Not available to set manually

	# gateway
	# Computed: true
	# Type: string
	# Available values: computed attribute. Not available to set manually

	# netmask
	# Computed: true
	# Type: string
	# Available values: computed attribute. Not available to set manually

	# mac_address
	# Computed: true
	# Type: string
	# Available values: computed attribute. Not available to set manually

	# interface_id
	# Computed: true
	# Type: int
	# Available values: computed attribute. Not available to set manually

	# dns_prefix
	# Computed: true
	# Type: string
	# Available values: computed attribute. Not available to set manually

	# dhcp_branch
	# Computed: true
	# Type: string
	# Available values: computed attribute. Not available to set manually

	# type_id
	# Computed: true
	# Type: int
	# Available values: computed attribute. Not available to set manually

	# creation_user_id
	# Computed: true
	# Type: int
	# Available values: computed attribute. Not available to set manually
