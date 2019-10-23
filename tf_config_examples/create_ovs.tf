provider "oktawave" {
  access_token="your-access-token"
  api_url = "https://api.oktawave.com/beta"
}

resource "oktawave_ovs" "my_ovs"{


	# Required: true
    # Type: string
    # Available values: any string of available length that will represent name of your disk.
	disk_name="my_disk2"
    
    
    # Optional: true
    # Type: int
    # Default value: 5
    # Available values: any available capacity for disk(in GB).
	space_capacity=5
    
    
    # Required: true
  	# Type: int
    # Available values: 48(Tier 1), 49(Tier2), 50(Tier3), 895(Tier4), 896(Tier5)
    # Comment: this is class of disk
	tier_id = 895
    
    
    # Optional: true
    # Type: bool
    # Default value: false
    # ForceNew: true
    # Available values: true, false
    # Comment: If set true - shared disk type id attribute is required
	is_shared = false //false by default
    
    
    # Optional: true
    # Type: int
    # Available values: 1411( LINUX (OCFS/GFS)), 1412( Windows (MSCS))
    # Comment: this is disk file system. Required when is_shared=true
	#shared_disk_type_id = 1411/1412
    
    
    # Required: true
    # Type: int
    # Available values: 1(PL001), 4(PL002), 5(PL003), 6(PL004), 7(PL005)
	subregion_id=4
    
    
    # Optional: true
    # Type: int
    # Default value: false
    # Available values: true. false
	is_locked=false //false by default
    
    
    # Optional: true
    # Type: int set
    # Available values: set of ints representing instance ids to which you want to attach disk(e.g [1, 2, ...])  
	#connections_with_instanceid=[ids of instances that should be connected with volume, ...]
    
    
    # Optional: true
    # Type: bool
    # Default value: false
    # Available values: set of ints representing instance ids to which you want to attach disk(e.g [1, 2, ...])      
	isfreemium = false
}