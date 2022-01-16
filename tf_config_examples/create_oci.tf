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

//you can omit optional parametres and parametres with default value
resource "oktawave_oci" "my_oci" {

		# Optional: true
        # Type: int
        # Default value: 1399
        # Available values: 1398 (ssh keys),1399 (login&password).
		authorization_method_id=1399 
             
        # Required: true
        # Type: int
        # Available values: 48(Tier 1), 49(Tier2), 50(Tier3), 895(Tier4), 896(Tier5)
        # Comment: this is class of instances disk which is created as part of instance setup
        disk_class =896
        
        # Optional: true
        # Type: int
        # Available values: id of ip address that you want to set as default
        # Comment: Replace default ip address that would be created as part of instance setup
        # init_ip_address_id=id
        
        # Optional: true
        # Type: int
        # Default value: 5
      	# Available values: any available capacity for disk(in GB) which is created as part of instance setup
	init_disk_size = 5
        
        # Required: true
        # Type: string
        # Available values: any string of available length
 	instance_name ="my_instance"
 		
        # Required: true
        # Type: int
        # Available values: 1(PL001), 4(PL002), 5(PL003), 6(PL004), 7(PL005)
        subregion_id =4
 		
        # Required: true
        # Type: int
        # Available values: 
                # "ID: 69  Zimbra"
                # "ID: 82  MS SQL 2012 Express (@ Windows 2012)"
                # "ID: 97  MS SQL 2014 Standard (@ Windows 2012 R2)"
                # "ID: 98  Windows Server 2012 R2"
                # "ID: 99  MS SQL 2014 Web (@ Windows 2012 R2)"
                # "ID: 101  MS SQL 2014 Express (@ Windows 2012 R2)"
                # "ID: 643  MS SQL 2016 Express"
                # "ID: 644  MS SQL 2016 Standard"
                # "ID: 645  MS SQL 2016 Web"
                # "ID: 1018  Ubuntu Server 18.04 LTS autobuild"
                # "ID: 1021  Ubuntu Server 20.04 LTS autobuild"
                # "ID: 1027  Windows Server 2019 autobuild"
                # "ID: 1028  Windows Server Core 2019 autobuild"
                # "ID: 1032  openSUSE Leap 15.1 autobuild"
                # "ID: 1037  FreeBSD 12.1 autobuild"
                # "ID: 1039  FreeBSD 11.3 autobuild"
                # "ID: 1042  Debian 9 autobuild"
                # "ID: 1043  Debian 10 autobuild"
                # "ID: 1045  Ubuntu Server 16.04 LTS autobuild"
                # "ID: 1048  CentOS 8 autobuild"
                # "ID: 1051  CentOS 7 autobuild"
                # "ID: 1063  Fedora Server 32 autobuild"
                # "ID: 1091  pfSense 2.4.5 autobuild"
                # "ID: 1095  Windows Server 2016 autobuild"

        template_id = 1021
        
        # Required: true
        # Type: int
        # Available values: 
		
			# v1.standard: 1047(v1.standard-1.05), 289(v1.standard-1.09), 34(v1.standard-2.2), 35(v1.standard-4.4), 36(v1.standard-8.8), 1048(v1.standard-16.16), 1766(v1.standard-24.24)
		
			# v1.highmemory-1/2.x: 1263(v1.highmemory-1.4), 1428(v1.highmemory-2.4), 1264(v1.highmemory-2.8)
        
			# v1.highmemory-4.x: 1423(v1.highmemory-4.8), 1265(v1.highmemory-4.16), 1420(v1.highmemory-4.32), 1767(v1.highmemory.4.48), 1421(v1.highmemory-4.64), 
		
			# v1.highmemory-8.x: 1424(v1.highmemory-8.16), 1266(v1.highmemory-8.32), 1768(v1.highmemory.8.48), 1422(v1.highmemory-8.64), 1574(v1.highmemory-8.96), 
		
			# v1.highmemory-16.x: 1757(v1.highmemory-16.24), 1049(v1.highmemory-16.32), 1769(v1.highmemory.16.48), 1050(v1.highmemory-16.64), 1267(v1.highmemory-16.96),
        
			# v1.highmemory-24.x:1765(v1.highmemory-24.32), 1770(v1.highmemory.24.48), 1759(v1.highmemory-24.64), 
		
			# v1.highcpu: 1268(v1.highcpu-2.09), 1269(v1.highcpu-4.2), 1270(v1.highcpu-8.4), 1271(v1.highcpu-16.8)
            
 		type_id = 1268
        
        
        # Optional: true
        # Type: int
        # Defauilt value: 1
        # Available values: number of instance you want to create.
        # Comment: For now function for his attribute doesn't work, so omit it
	instances_count = 1
		
        #Optional: true
        #Type: bool
        #Default value: false
        #Available values: true, false
        isfreemium=false 
}

		# ADDITIONAL INFO/INFORMACJE DODATKOWE
		
        # OPTIONAL SETS/ZBIORY OPCJONALNE
		
        # Optional: true
       	# Type: int set
        # Available values: set of int ids of OPNs to which you want to connect your instance
        # opn_ids=[1, 2, 3]
        
        
        # Optional: true
        # Type: int set
        # Available values: set of ids(e.g. [1, 2, 3]) of ssh keys
        # ssh_keys_ids=[1, 2, 3] //optional. Required if you want to use authorization_method_id 1398

	# Optional: true
	# Type: int set
	# Available values: set of ids(e.g. [1, 2, 3]) of ips that you want to attach to this instance
	# Comment: Here you define non-default ip ids. Default ip id should be defined in init_ip_address_id field
	# ip_address_ids=[1, 2, 3]

        
        
		#COMPUTED VALUES/WARTOŚCI WYLICZALNE AUTOMATYCZNIE
	# creation_date
        # Computed: true
        # Type: string
        # Available values: computed value, not allowed to set
            
        # islocked
        # Computed: true
        # Type: bool
        # Available values: computed value, not allowed to set
        
        # creation_userid
        # Computed: true
        # Type: bool
        # Available values: computed value, not allowed to set
        
        # init_ip_address
        # Computed: true
        # Type: string
        # Available values: computed value, not allowed to set
            
        # dns_address
        # Computed: true
        # Type: string
        # Available values: computed value, not allowed to set

	# init_disk_id=id
	# Computed: true
        # Type: int
        # Available values: computed value, not allowed to set
	
	# opn_mac = { opn_name = "mac_address" }
	# Computed: true
	# Type: String
        # Available values: computed value, not allowed to set
