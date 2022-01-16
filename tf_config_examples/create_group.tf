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



resource "oktawave_group" "my_group"{
	# Required: true
    # Type: string
    # Available values: any string of available length that will represent name of your group.
	group_name="my_group1"
	
	# Optional: true
    # Type: int
	# Default value: 1403
    # Available values: 1403(No separation), 1404(Minimize separation), 1405(Maximize separation).
	affinity_rule_type_id=1403
}







	#ADDITIONAL INFO/INFORMACJE DODATKOWE

    #ADDITIONAL MAPS/DODATKOWE MAPY

    # Optional: true
    # Type: map[string]int
    # Available values: int-int pair representing instance id as a key(it's format is string but in fact it should be int) and instance ip id as value.
	# Comment: if you set value 0 - pair instance id-it's default ip will be attached to group
	# group_instance_ip_ids={
	#	"78982":0
	# }


    #COMPUTED VALUES/WARTOŚCI WYLICZALNE AUTOMATYCZNIE

	# schedulers_count
	# Computed: true
    # Type: int
	# Available values: computed attribute. Not available to set manually

	# last_change_date
	# Computed: true
    # Type: string
	# Available values: computed attribute. Not available to set manually


	# creation_user_id
	# Computed: true
    # Type: int
	# Available values: computed attribute. Not available to set manually


	# instances_count
	# Computed: true
    # Type: int
	# Available values: computed attribute. Not available to set manually
