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
   api_url = "https://api.oktawave.com/beta"
}

resource "oktawave_oci" "my_oci" {
	authorization_method_id=1399 
	disk_class =896
	disk_size = 5 
	instance_name ="my_instance"
	subregion_id =4
	template_id =94
	type_id = 1268
	instances_count = 1 
	ip_address_id = 0 
	ssh_keys_ids=[] 
	isfreemium=false
}

resource "oktawave_ovs" "my_ovs"{
	disk_name="my_disk2"
	space_capacity=5 //5 by default
	tier_id = 895
	is_shared = false //false by default
	subregion_id=4
	is_locked=false 
	connections_with_instanceids=[oktawave_oci.my_oci.id] //instances subregion == ovs subregion. Other scenario == error
	isfreemium = false 
	
	depends_on=[oktawave_oci.my_oci]
}