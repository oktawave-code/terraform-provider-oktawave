package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/oktawave-code/odk"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestAccOktawave_IpAddress_Basic(t *testing.T) {
	var ip odk.Ip
	mockStatus := os.Getenv("MOCK_STATUS")
	token := os.Getenv("TOKEN")
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOktawaveIpAddrDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOktawaveIpAddrConfig_basic(token, mockStatus),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOktawaveIpAddrExists("oktawave_ip_address.my_ip1", &ip),
					testAccCheckOktawaveIpAddrAttributes_basic(&ip),
					resource.TestCheckResourceAttr("oktawave_ip_address.my_ip1", "comment", ""),
					resource.TestCheckResourceAttr("oktawave_ip_address.my_ip1", "subregion_id", "4"),
					resource.TestCheckResourceAttr("oktawave_ip_address.my_ip1", "type_id", "1106"),
				),
				//ExpectError: regexp.MustCompile("errors during apply: Resource OCI. CREATE. Unable to create instance . Ticket status: 137"),
			},
		},
	})
	httpmock.DeactivateAndReset()
}

func TestAccOktawave_IpAddress_Update(t *testing.T) {
	var ip odk.Ip
	mockStatus := os.Getenv("MOCK_STATUS")
	token := os.Getenv("TOKEN")
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOktawaveIpAddrDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOktawaveIpAddrConfig_basic(token, mockStatus),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOktawaveIpAddrExists("oktawave_ip_address.my_ip1", &ip),
					testAccCheckOktawaveIpAddrAttributes_basic(&ip),
					resource.TestCheckResourceAttr("oktawave_ip_address.my_ip1", "comment", ""),
					resource.TestCheckResourceAttr("oktawave_ip_address.my_ip1", "subregion_id", "4"),
					resource.TestCheckResourceAttr("oktawave_ip_address.my_ip1", "type_id", "1106"),
				),
				//ExpectError: regexp.MustCompile("errors during apply: Resource OCI. CREATE. Unable to create instance . Ticket status: 137"),
			},
			{
				Config: testAccCheckOktawaveIpAddrConfig_Updated(token, mockStatus),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOktawaveIpAddrExists("oktawave_ip_address.my_ip1", &ip),
					testAccCheckOktawaveIpAddrAttributes_Updated(&ip),
					resource.TestCheckResourceAttr("oktawave_ip_address.my_ip1", "comment", "example ip_updated"),
					resource.TestCheckResourceAttr("oktawave_ip_address.my_ip1", "subregion_id", "4"),
					resource.TestCheckResourceAttr("oktawave_ip_address.my_ip1", "type_id", "1106"),
				),
			},
		},
	})

	httpmock.DeactivateAndReset()
}

func testAccCheckOktawaveIpAddrExists(name string, ip *odk.Ip) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ip ID is set")
		}

		client := testAccProvider.Meta().(*ClientConfig).oktaClient()
		auth := testAccProvider.Meta().(*ClientConfig).ctx
		id, _ := strconv.Atoi(rs.Primary.ID)
		foundIp, _, err := client.OCIInterfacesApi.InstancesGetInstanceIp(*auth, int32(id), nil)
		if err != nil {
			return fmt.Errorf("ip was not found by id: %s", rs.Primary.ID)
		}
		log.Printf("IP was found by id %s ", rs.Primary.ID)
		*ip = foundIp
		return nil
	}
}

func testAccCheckOktawaveIpAddrAttributes_basic(ip *odk.Ip) resource.TestCheckFunc {
	return func(t *terraform.State) error {
		if ip.Subregion.Id != 4 {
			return fmt.Errorf("Bad IP subregion id. Expected: 4. Got: %v", strconv.Itoa(int(ip.Subregion.Id)))
		}
		if ip.Comment != "" {
			return fmt.Errorf("Bad IP comment. Expected: . Got: %v", ip.Comment)
		}
		if ip.Type_.Id != 1106 {
			return fmt.Errorf("Bad ip type id. Expected: 1106. Got: %s", strconv.Itoa(int(ip.Type_.Id)))
		}
		return nil
	}
}

func testAccCheckOktawaveIpAddrAttributes_Updated(ip *odk.Ip) resource.TestCheckFunc {
	return func(t *terraform.State) error {
		if ip.Comment != "example ip_updated" {
			return fmt.Errorf("Bad IP comment. Expected: example ip_updated. Got: %v", ip.Comment)
		}
		if ip.Type_.Id != 1106 {
			return fmt.Errorf("Bad ip type id. Expected: 1106. Got: %s", strconv.Itoa(int(ip.Type_.Id)))
		}
		//if ip.Instance==nil{
		//	return fmt.Errorf("Bad ip attachment. Instance parameter expected not to be null")
		//}
		return nil
	}
}

func testAccCheckOktawaveIpAddrDestroy(s *terraform.State) error {
	if os.Getenv("MOCK_STATUS") == "1" {
		httpmock.RegisterNoResponder(httpmock.NewStringResponder(http.StatusNotFound, ""))
	}
	client := testAccProvider.Meta().(*ClientConfig).oktaClient()
	auth := testAccProvider.Meta().(*ClientConfig).ctx
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oktawave_ip_address" {
			id, _ := strconv.Atoi(rs.Primary.ID)
			_, resp, err := client.OCIInterfacesApi.InstancesGetInstanceIp(*auth, int32(id), nil)
			if err != nil && resp.StatusCode != 404 {
				return fmt.Errorf("Error waitiing IP to be destroyed. Error: %v", err.Error())
			}
			break
		}
	}
	return nil
}

func testAccCheckOktawaveIpAddrConfig_basic(token string, mockStatus string) string {
	subregion_id := 4
	//comment := "example ip"
	if mockStatus == "1" {
		httpmock.Activate()
		mockPostIP()
		mockPostAttachIpTicket()
		mockGetAttachIpTicket()
		mockGetIp(int32(subregion_id), "")
		mockDeleteIp()
	}
	return fmt.Sprintf(`
provider "oktawave" {
  access_token="%s"
  
  api_url = "https://api.oktawave.com/beta/"
}

resource "oktawave_ip_address" "my_ip1"{
	subregion_id=%s
}`, token, strconv.Itoa(subregion_id))
}

func testAccCheckOktawaveIpAddrConfig_Updated(token string, mockStatus string) string {

	subregion_id := 4
	comment := "example ip_updated"

	if mockStatus == "1" {
		mockPutIP()
	}
	return fmt.Sprintf(`
provider "oktawave" {
  access_token="%s"
  
  api_url = "https://api.oktawave.com/beta/"
}

resource "oktawave_ip_address" "my_ip1"{
	subregion_id=%s
	comment="%s"
}`, token, strconv.Itoa(subregion_id), comment)
}

func mockPostIP() {
	httpmock.RegisterResponder(http.MethodPost, "https://api.oktawave.com/beta//instances/ip_addresses",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(http.StatusOK, odk.Ip{
				Id:           1000,
				Address:      "",
				AddressV6:    "",
				Gateway:      "",
				NetMask:      "",
				MacAddress:   "",
				InterfaceId:  0,
				DnsPrefix:    "",
				Subregion:    &odk.BaseResource{},
				Type_:        &odk.DictionaryItem{},
				OwnerAccount: &odk.BaseResource{},
				Comment:      "",
				RevDns:       "",
				RevDnsV6:     "",
				CreationUser: &odk.UserResource{},
			})
		})
}

func mockPostAttachIpTicket() {
	httpmock.RegisterResponder(http.MethodPost, "https://api.oktawave.com/beta/instances/1658/attach_ip_ticket",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(http.StatusOK, odk.Ticket{
				Id:            1001,
				CreationDate:  time.Now(),
				CreationUser:  &odk.UserResource{},
				EndDate:       time.Now(),
				Status:        &odk.DictionaryItem{Id: TICKET_STATUS__SUCCESS},
				OperationType: &odk.DictionaryItem{},
				ObjectId:      1000,
				ObjectName:    "",
				Progress:      100,
			})
		})
}

func mockGetAttachIpTicket() {
	httpmock.RegisterResponder(http.MethodGet, "https://api.oktawave.com/beta/ticket/1001",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(http.StatusOK, odk.Ticket{
				Id:            1001,
				CreationDate:  time.Now(),
				CreationUser:  &odk.UserResource{},
				EndDate:       time.Now(),
				Status:        &odk.DictionaryItem{Id: TICKET_STATUS__SUCCESS},
				OperationType: &odk.DictionaryItem{},
				ObjectId:      1000,
				ObjectName:    "",
				Progress:      100,
			})
		})
}

func mockGetIp(subregion_id int32, comment string) {
	httpmock.RegisterResponder(http.MethodGet, "https://api.oktawave.com/beta//instances/ip_addresses/1000",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(http.StatusOK, odk.Ip{
				Id:           1000,
				Address:      "",
				AddressV6:    "",
				Gateway:      "",
				NetMask:      "",
				Instance:     &odk.BaseResource{},
				MacAddress:   "",
				InterfaceId:  0,
				DnsPrefix:    "",
				Subregion:    &odk.BaseResource{Id: subregion_id},
				Type_:        &odk.DictionaryItem{Id: 1106},
				OwnerAccount: &odk.BaseResource{},
				Comment:      comment,
				RevDns:       "",
				RevDnsV6:     "",
				CreationUser: &odk.UserResource{},
			})
		})

}

func mockDeleteIp() {
	httpmock.RegisterResponder(http.MethodDelete, "https://api.oktawave.com/beta//instances/ip_addresses/1000",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(http.StatusOK, odk.Object{})
		})
}

func mockPutIP() {
	httpmock.RegisterResponder(http.MethodPut, "https://api.oktawave.com/beta//instances/ip_addresses/1000",
		func(req *http.Request) (*http.Response, error) {

			mockGetIp(4, "example ip_updated")

			return httpmock.NewJsonResponse(http.StatusOK, odk.Object{})
		})
}
