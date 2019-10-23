provider "oktawave" {
  access_token="your-access-token"
  api_url = "https://api.oktawave.com/beta"
}

resource "oktawave_sshKey" "my_key"{
	# Required: true
    # Type: string
	# ForceNew: true
    # Available values: string of available length that represent name of ssh key
	ssh_key_name="my_sshKey"
	
	# Required: true
    # Type: string
	# ForceNew: true
    # Available values: ssh-rsa key value
	ssh_key_value="your-ssh-key"
}

# ADDITIONAL INFO

# COMPUTED VALUES

# owner_user_id
# Computed: true
# Type: int
# Available values: computed attribute is not allowed to set manually

# creation_date
# Computed: true
# Type: string
# Available values: computed attribute is not allowed to set manually
