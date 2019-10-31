provider "oktawave"{
	access_token="your-access-token"
}

resource "oktawave_kubernetes_node" "my_node" {


	# Required: true
	# Type: int
	# Available values:

			# v1.standard: 34(v1.standard-2.2), 35(v1.standard-4.4), 36(v1.standard-8.8), 1048(v1.standard-16.16), 1766(v1.standard-24.24)

			# v1.highmemory-1/2.x:, 1428(v1.highmemory-2.4), 1264(v1.highmemory-2.8)

			# v1.highmemory-4.x: 1423(v1.highmemory-4.8), 1265(v1.highmemory-4.16), 1420(v1.highmemory-4.32), 1767(v1.highmemory.4.48), 1421(v1.highmemory-4.64),

			# v1.highmemory-8.x: 1424(v1.highmemory-8.16), 1266(v1.highmemory-8.32), 1768(v1.highmemory.8.48), 1422(v1.highmemory-8.64), 1574(v1.highmemory-8.96),

			# v1.highmemory-16.x: 1757(v1.highmemory-16.24), 1049(v1.highmemory-16.32), 1769(v1.highmemory.16.48), 1050(v1.highmemory-16.64), 1267(v1.highmemory-16.96),

			# v1.highmemory-24.x:1765(v1.highmemory-24.32), 1770(v1.highmemory.24.48), 1759(v1.highmemory-24.64),

			# v1.highcpu:1269(v1.highcpu-4.2), 1270(v1.highcpu-8.4), 1271(v1.highcpu-16.8)
			
	type_id=35
	
	
	# Required: true
    # Type: int  
	# Available values: 1(PL001), 4(PL002), 5(PL003), 6(PL004), 7(PL005)
    subregion_id=4
	
	
	# Required: true
	# Type: int
	# Available values: full name of cluster you want connect node to
	cluster_name="tfcluster-01uuep4mg"
}