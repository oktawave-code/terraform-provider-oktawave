package main

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/oktawave-code/odk"
	"net/http"
	"os"
	"strconv"
	"testing"
)

func TestAccOktawave_Load_Balancer_Basic(t *testing.T) {
	var lb odk.LoadBalancer
	//Environment variable to manage sticket status
	//TODO: make different creationg and reading scenarios using environment variables
	token := os.Getenv("TOKEN")
	mockStatus := os.Getenv("MOCK_STATUS")
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOktawaveLoad_BalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOktawaveLoad_BalancerConfig_basic(token, mockStatus),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOktawaveLoad_BalancerExists("oktawave_load_balancer.my_lb", &lb),
					testAccCheckOktawaveLoad_BalancerAttributes_basic(&lb),
					resource.TestCheckResourceAttr("oktawave_load_balancer.my_lb", "ssl_enabled", "true"),
					resource.TestCheckResourceAttr("oktawave_load_balancer.my_lb", "service_type_id", "43"),
					resource.TestCheckResourceAttr("oktawave_load_balancer.my_lb", "port_number", "1250"),
					resource.TestCheckResourceAttr("oktawave_load_balancer.my_lb", "session_persistence_type_id", "46"),
					resource.TestCheckResourceAttr("oktawave_load_balancer.my_lb", "load_balancer_algorithm_id", "612"),
					resource.TestCheckResourceAttr("oktawave_load_balancer.my_lb", "ip_version_id", "115"),
					resource.TestCheckResourceAttr("oktawave_load_balancer.my_lb", "health_check_enabled", "true"),
					resource.TestCheckResourceAttr("oktawave_load_balancer.my_lb", "common_persistence_for_http_and_https_enabled", "true"),
				),
				//ExpectError: regexp.MustCompile("errors during apply: Resource OCI. CREATE. Unable to create instance . Ticket status: 137"),
			},
		},
	})
	httpmock.DeactivateAndReset()
}

func TestAccOktawave_Load_Balancer_Updated(t *testing.T) {
	var lb odk.LoadBalancer
	token := os.Getenv("TOKEN")
	mockStatus := os.Getenv("MOCK_STATUS")
	//Environment variable to manage sticket status
	//TODO: make different creationg and reading scenarios using environment variables
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOktawaveLoad_BalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOktawaveLoad_BalancerConfig_basic(token, mockStatus),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOktawaveLoad_BalancerExists("oktawave_load_balancer.my_lb", &lb),
					testAccCheckOktawaveLoad_BalancerAttributes_basic(&lb),
					resource.TestCheckResourceAttr("oktawave_load_balancer.my_lb", "ssl_enabled", "true"),
					resource.TestCheckResourceAttr("oktawave_load_balancer.my_lb", "service_type_id", "43"),
					resource.TestCheckResourceAttr("oktawave_load_balancer.my_lb", "port_number", "1250"),
					resource.TestCheckResourceAttr("oktawave_load_balancer.my_lb", "session_persistence_type_id", "46"),
					resource.TestCheckResourceAttr("oktawave_load_balancer.my_lb", "load_balancer_algorithm_id", "612"),
					resource.TestCheckResourceAttr("oktawave_load_balancer.my_lb", "ip_version_id", "115"),
					resource.TestCheckResourceAttr("oktawave_load_balancer.my_lb", "health_check_enabled", "true"),
					resource.TestCheckResourceAttr("oktawave_load_balancer.my_lb", "common_persistence_for_http_and_https_enabled", "true"),
				),
				//ExpectError: regexp.MustCompile("errors during apply: Resource OCI. CREATE. Unable to create instance . Ticket status: 137"),
			},
			{
				Config: testAccCheckOktawaveLoad_BalancerConfig_Updated(token, mockStatus),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOktawaveLoad_BalancerExists("oktawave_load_balancer.my_lb", &lb),
					testAccCheckOktawaveLoad_BalancerAttributes_Updated(&lb),
					resource.TestCheckResourceAttr("oktawave_load_balancer.my_lb", "ssl_enabled", "true"),
					resource.TestCheckResourceAttr("oktawave_load_balancer.my_lb", "service_type_id", "43"),
					resource.TestCheckResourceAttr("oktawave_load_balancer.my_lb", "port_number", "1255"),
					resource.TestCheckResourceAttr("oktawave_load_balancer.my_lb", "session_persistence_type_id", "46"),
					resource.TestCheckResourceAttr("oktawave_load_balancer.my_lb", "load_balancer_algorithm_id", "288"),
					resource.TestCheckResourceAttr("oktawave_load_balancer.my_lb", "ip_version_id", "115"),
					resource.TestCheckResourceAttr("oktawave_load_balancer.my_lb", "health_check_enabled", "true"),
					resource.TestCheckResourceAttr("oktawave_load_balancer.my_lb", "common_persistence_for_http_and_https_enabled", "true"),
				),
			},
		},
	})
}

func testAccCheckOktawaveLoad_BalancerExists(name string, lb *odk.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		client := testAccProvider.Meta().(*ClientConfig).oktaClient()
		auth := testAccProvider.Meta().(*ClientConfig).ctx
		id, _ := strconv.Atoi(rs.Primary.ID)
		groupId := id - 1
		foundLB, _, err := client.OCIGroupsApi.LoadBalancersGetLoadBalancer(*auth, (int32)(groupId), nil)
		if err != nil {
			return fmt.Errorf("Instance not found %s ", err)
		}
		*lb = foundLB
		return nil
	}
}

func testAccCheckOktawaveLoad_BalancerAttributes_basic(lb *odk.LoadBalancer) resource.TestCheckFunc {
	return func(t *terraform.State) error {
		if lb.SslEnabled != true {
			return fmt.Errorf("Bad load balancers ssl enabled option. Expected: true. Got: %s", strconv.FormatBool(lb.SslEnabled))
		}
		if lb.CommonPersistenceForHttpAndHttpsEnabled != true {
			return fmt.Errorf("Bad load balancers common persistence option. Expected: true. Got: %s", strconv.FormatBool(lb.CommonPersistenceForHttpAndHttpsEnabled))
		}
		if lb.HealthCheckEnabled != true {
			return fmt.Errorf("Bad load balancers health check option. Expected: true. Got: %s", strconv.FormatBool(lb.HealthCheckEnabled))
		}
		if lb.IpVersion.Id != 115 {
			return fmt.Errorf("Bad load balancers ip version id. Expected: 115. Got: %s", strconv.Itoa(int(lb.IpVersion.Id)))
		}
		if lb.SessionPersistenceType.Id != 46 {
			return fmt.Errorf("Bad load balancers session persistence type id. Expected: 46. Got: %s", strconv.Itoa(int(lb.SessionPersistenceType.Id)))
		}
		if lb.PortNumber != 1250 {
			return fmt.Errorf("Bad load balancers port. Expected: 1250. Got: %s", strconv.Itoa(int(lb.PortNumber)))
		}
		if lb.Algorithm.Id != 612 {
			return fmt.Errorf("Bad load balancers algorithm id. Expected: 612. Got: %s", strconv.Itoa(int(lb.Algorithm.Id)))
		}
		if lb.ServiceType.Id != 43 {
			return fmt.Errorf("Bad load balancers service type id. Expected: 43. Got: %s", strconv.Itoa(int(lb.ServiceType.Id)))
		}
		return nil
	}
}

func testAccCheckOktawaveLoad_BalancerAttributes_Updated(lb *odk.LoadBalancer) resource.TestCheckFunc {
	return func(t *terraform.State) error {
		if lb.SslEnabled != true {
			return fmt.Errorf("Bad load balancers ssl enabled option. Expected: true. Got: %s", strconv.FormatBool(lb.SslEnabled))
		}
		if lb.CommonPersistenceForHttpAndHttpsEnabled != true {
			return fmt.Errorf("Bad load balancers common persistence option. Expected: true. Got: %s", strconv.FormatBool(lb.CommonPersistenceForHttpAndHttpsEnabled))
		}
		if lb.HealthCheckEnabled != true {
			return fmt.Errorf("Bad load balancers health check option. Expected: true. Got: %s", strconv.FormatBool(lb.HealthCheckEnabled))
		}
		if lb.IpVersion.Id != 115 {
			return fmt.Errorf("Bad load balancers ip version id. Expected: 115. Got: %s", strconv.Itoa(int(lb.IpVersion.Id)))
		}
		if lb.SessionPersistenceType.Id != 46 {
			return fmt.Errorf("Bad load balancers session persistence type id. Expected: 46. Got: %s", strconv.Itoa(int(lb.SessionPersistenceType.Id)))
		}
		if lb.PortNumber != 1255 {
			return fmt.Errorf("Bad load balancers port. Expected: 1250. Got: %s", strconv.Itoa(int(lb.PortNumber)))
		}
		if lb.Algorithm.Id != 288 {
			return fmt.Errorf("Bad load balancers algorithm id. Expected: 288. Got: %s", strconv.Itoa(int(lb.Algorithm.Id)))
		}
		if lb.ServiceType.Id != 43 {
			return fmt.Errorf("Bad load balancers service type id. Expected: 43. Got: %s", strconv.Itoa(int(lb.ServiceType.Id)))
		}
		return nil
	}
}

func testAccCheckOktawaveLoad_BalancerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ClientConfig).oktaClient()
	auth := testAccProvider.Meta().(*ClientConfig).ctx
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oktawave_load_balancer" {
			id, _ := strconv.Atoi(rs.Primary.ID)
			_, resp, err := client.OCIGroupsApi.LoadBalancersGetLoadBalancer(*auth, int32(id-1), nil)
			if err != nil {
				if resp != nil && resp.StatusCode == 404 {
					return nil
				}
				return fmt.Errorf("Error waitiing LB to be destroyed. Error: %v", err.Error())
			}
			return nil
		}

	}
	return fmt.Errorf("Check destroy Group  test function error: resource type was not found")
}

func testAccCheckOktawaveLoad_BalancerConfig_basic(token string, mockStatus string) string {
	port_number := 1250
	ssl_enabled := true
	cmnPersistForHttps := true
	ipVersion := 115
	sessionPersistenceType := 46
	algorithmId := 612
	serviceType := 43
	if mockStatus == "1" {
		httpmock.Activate()
		mockGroupPost("my_group", true, 1403)
		mockGetGroup("my_group", true, 1403)
		mockGetGroupAssignments_empty()
		mockDeleteGroup()
		mockPostLB(int32(serviceType), int32(port_number), int32(sessionPersistenceType), int32(algorithmId), int32(ipVersion), ssl_enabled, cmnPersistForHttps)
		mockGetLB(int32(serviceType), int32(port_number), int32(sessionPersistenceType), int32(algorithmId), int32(ipVersion), ssl_enabled, cmnPersistForHttps)
		mockDeleteLB()
	}
	return fmt.Sprintf(`
provider "oktawave" {
 access_token="%s"

 api_url = "https://api.oktawave.com/beta/"
}


resource "oktawave_group" "my_group"{
	group_name="my_group"
	//group_instance_ids=[78914]
}




resource "oktawave_load_balancer" "my_lb"{
	group_id = oktawave_group.my_group.id
	port_number=1250
	ssl_enabled=true
	service_type_id = 43
	session_persistence_type_id=46
	load_balancer_algorithm_id=612
	ip_version_id=115
	health_check_enabled=true
	common_persistence_for_http_and_https_enabled=true
	//load_balancer_ip_id=2823
	depends_on=[oktawave_group.my_group]
}`, token)
}

func testAccCheckOktawaveLoad_BalancerConfig_Updated(token string, mockStatus string) string {
	port_number := 1255
	ssl_enabled := true
	cmnPersistForHttps := true
	ipVersion := 115
	sessionPersistenceType := 46
	algorithmId := 288
	serviceType := 43
	if mockStatus == "1" {
		mockPutLB(int32(serviceType), int32(port_number), int32(sessionPersistenceType), int32(algorithmId), int32(ipVersion), ssl_enabled, cmnPersistForHttps)
	}

	return fmt.Sprintf(`
provider "oktawave" {
 access_token="%s"

 api_url = "https://api.oktawave.com/beta/"
}


resource "oktawave_group" "my_group"{
	group_name="my_group"
	//group_instance_ids=[78914]
}




resource "oktawave_load_balancer" "my_lb"{
	group_id = oktawave_group.my_group.id
	port_number=1255
	load_balancer_algorithm_id=288
	//load_balancer_ip_id=2823
	depends_on=[oktawave_group.my_group]
}`, token)
}

func mockPostLB(serviceType int32, port int32, sessionPersistenceType int32, algorithmId int32, ipVersion int32, sslEnabled bool, isCommonPersistence bool) {
	httpmock.RegisterResponder(http.MethodPost, "https://api.oktawave.com/beta//groups/1658/load_balancer",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(http.StatusOK, odk.LoadBalancer{
				GroupId:     1658,
				GroupName:   "my_group",
				IpAddress:   "",
				IpV6Address: "",
				ServiceType: &odk.DictionaryItem{
					Id: serviceType,
				},
				PortNumber:       port,
				TargetPortNumber: 0,
				SessionPersistenceType: &odk.DictionaryItem{
					Id: sessionPersistenceType,
				},
				Algorithm: &odk.DictionaryItem{
					Id: algorithmId,
				},
				IpVersion: &odk.DictionaryItem{
					Id: ipVersion,
				},
				HealthCheckEnabled:                      true,
				SslEnabled:                              sslEnabled,
				CommonPersistenceForHttpAndHttpsEnabled: isCommonPersistence,
				Servers:                                 make([]odk.LoadBalancerServer, 0, 0),
			})
		})
}

func mockGetLB(serviceType int32, port int32, sessionPersistenceType int32, algorithmId int32, ipVersion int32, sslEnabled bool, isCommonPersistence bool) {
	httpmock.RegisterResponder(http.MethodGet, "https://api.oktawave.com/beta//groups/1658/load_balancer",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(http.StatusOK, odk.LoadBalancer{
				GroupId:     1658,
				GroupName:   "my_group",
				IpAddress:   "",
				IpV6Address: "",
				ServiceType: &odk.DictionaryItem{
					Id: serviceType,
				},
				PortNumber:       port,
				TargetPortNumber: 0,
				SessionPersistenceType: &odk.DictionaryItem{
					Id: sessionPersistenceType,
				},
				Algorithm: &odk.DictionaryItem{
					Id: algorithmId,
				},
				IpVersion: &odk.DictionaryItem{
					Id: ipVersion,
				},
				HealthCheckEnabled:                      true,
				SslEnabled:                              sslEnabled,
				CommonPersistenceForHttpAndHttpsEnabled: isCommonPersistence,
				Servers:                                 make([]odk.LoadBalancerServer, 0, 0),
			})
		})
}

func mockDeleteLB() {
	httpmock.RegisterResponder(http.MethodDelete, "https://api.oktawave.com/beta//groups/1658/load_balancer",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(http.StatusOK, odk.Object{})
		})
}

func mockPutLB(serviceType int32, port int32, sessionPersistenceType int32, algorithmId int32, ipVersion int32, sslEnabled bool, isCommonPersistence bool) {
	httpmock.RegisterResponder(http.MethodPut, "https://api.oktawave.com/beta//groups/1658/load_balancer",
		func(req *http.Request) (*http.Response, error) {
			mockGetLB(serviceType, port, sessionPersistenceType, algorithmId, ipVersion, sslEnabled, isCommonPersistence)
			return httpmock.NewJsonResponse(http.StatusOK, odk.LoadBalancer{
				GroupId:     1658,
				GroupName:   "my_group",
				IpAddress:   "",
				IpV6Address: "",
				ServiceType: &odk.DictionaryItem{
					Id: serviceType,
				},
				PortNumber:       port,
				TargetPortNumber: 0,
				SessionPersistenceType: &odk.DictionaryItem{
					Id: sessionPersistenceType,
				},
				Algorithm: &odk.DictionaryItem{
					Id: algorithmId,
				},
				IpVersion: &odk.DictionaryItem{
					Id: ipVersion,
				},
				HealthCheckEnabled:                      true,
				SslEnabled:                              sslEnabled,
				CommonPersistenceForHttpAndHttpsEnabled: isCommonPersistence,
				Servers:                                 make([]odk.LoadBalancerServer, 0, 0),
			})
		})
}
