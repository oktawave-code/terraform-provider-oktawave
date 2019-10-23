provider "oktawave" {
 access_token="your-access-token"
 api_url = "https://api.oktawave.com/beta/"
}


resource "oktawave_group" "my_group"{
	group_name="my_group"
}




resource "oktawave_load_balancer" "my_lb"{
	# Required: true
    # Type: int
	# ForceNew: true
    # Available values: group id
	group_id = oktawave_group.my_group.id
	
	# Optional: true
    # Type: int
    # Available values: any port from 1 to 65536
	port_number=1255
	
	# Optional: true
    # Type: bool
    # Available values: true, false
	ssl_enabled=true
	
	
	# Optional: true
    # Type: int
	# Default value: 43
    # Available values: 43 (HTTP (80)), 44 (HTTPS (443)), 45(SMTP (25)), 155 (PORT)
	service_type_id=43
	
	# Optional: true
	# Type: int
	# Available values: target port from 0 to 65536
	# target_port_number
	
	# Optional: true
	# Type: int
	# Default value: 46
	# Available values: 46 (Persistance by source IP), 280 (Persistance by Cookie), 47 (None persistance)
	session_persistence_type_id=46
	
	# Optional: true
	# Type: int
	# Default value: 612
	# Available values: 281 (Least Connection), 282 (Least Response Time), 288 (Source IP Hash), 612 (Round Robin)
	load_balancer_algorithm_id=612
	
	# Optional: true
	# Type: int
	# Default value: 115
	# Available values: 115 (IPv4), 116 (IPv6)
	ip_version_id=115
	
	# Optional: true
	# Type: bool
	# Default value: true
	# Available values: true, false
	health_check_enabled=true
	
	# Optional: true
	# Type: bool
	# Default value: true
	# Available values: true, false	
	common_persistence_for_http_and_https_enabled=true
	
	depends_on=[oktawave_group.my_group]
}