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
	"strings"
	"testing"
	"time"
)

func TestAccOktawave_Opn_Basic(t *testing.T) {
	var opn odk.Opn
	token := os.Getenv("TOKEN")
	mockStatus := os.Getenv("MOCK_STATUS")
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOktawaveOpnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOktawaveOpnConfig_basic(token, mockStatus),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOktawaveOpnExists("oktawave_opn.my_opn", &opn),
					testAccCheckOktawaveOpnAttributes_basic(&opn),
					resource.TestCheckResourceAttr("oktawave_opn.my_opn", "opn_name", "test_opn"),
				),
				//ExpectError: regexp.MustCompile("errors during apply: Resource OCI. CREATE. Unable to create instance . Ticket status: 137"),
			},
		},
	})

	httpmock.DeactivateAndReset()
}

func TestAccOktawave_Opn_Update(t *testing.T) {
	var opn odk.Opn
	token := os.Getenv("TOKEN")
	mockStatus := os.Getenv("MOCK_STATUS")
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOktawaveOpnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOktawaveOpnConfig_basic(token, mockStatus),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOktawaveOpnExists("oktawave_opn.my_opn", &opn),
					testAccCheckOktawaveOpnAttributes_basic(&opn),
					resource.TestCheckResourceAttr("oktawave_opn.my_opn", "opn_name", "test_opn"),
				),
				//ExpectError: regexp.MustCompile("errors during apply: Resource OCI. CREATE. Unable to create instance . Ticket status: 137"),
			},
			{
				Config: testAccCheckOktawaveOpnConfig_updated(token, mockStatus),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOktawaveOpnExists("oktawave_opn.my_opn", &opn),
					testAccCheckOktawaveOpnAttributes_updated(&opn),
					resource.TestCheckResourceAttr("oktawave_opn.my_opn", "opn_name", "test_opn_upd"),
				),
				//ExpectError: regexp.MustCompile("errors during apply: Resource OCI. CREATE. Unable to create instance . Ticket status: 137"),
			},
		},
	})

	httpmock.DeactivateAndReset()
}

func testAccCheckOktawaveOpnExists(name string, opn *odk.Opn) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No OPN ID is set")
		}

		client := testAccProvider.Meta().(*ClientConfig).oktaClient()
		auth := testAccProvider.Meta().(*ClientConfig).ctx
		id, _ := strconv.Atoi(rs.Primary.ID)
		opnCollection, _, err := client.NetworkingApi.OpnsGet(*auth, nil)
		if err != nil {
			return fmt.Errorf("Error occured when tryed to get list of OPNs: %s", err)
		}
		foundOpn, err := findOpnById(int32(id), opnCollection)
		if err != nil {
			return err
		}
		log.Printf("OPN was found by id %s and name %v", rs.Primary.ID, foundOpn.Name)
		*opn = foundOpn
		return nil
	}
}

func testAccCheckOktawaveOpnAttributes_basic(opn *odk.Opn) resource.TestCheckFunc {
	return func(t *terraform.State) error {
		if opn.Name != "test_opn" {
			return fmt.Errorf("Bad OPN name. Expected: test_opn. Got: %v", opn.Name)
		}
		return nil
	}
}

func testAccCheckOktawaveOpnAttributes_updated(opn *odk.Opn) resource.TestCheckFunc {
	return func(t *terraform.State) error {
		if opn.Name != "test_opn_upd" {
			return fmt.Errorf("Bad OPN name. Expected: test_opn_upd. Got: %v", opn.Name)
		}
		return nil
	}
}

func testAccCheckOktawaveOpnDestroy(s *terraform.State) error {
	if os.Getenv("MOCK_STATUS") == "1" {
		httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, "Instance not found"))
	}
	client := testAccProvider.Meta().(*ClientConfig).oktaClient()
	auth := testAccProvider.Meta().(*ClientConfig).ctx
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oktawave_opn" {
			id, _ := strconv.Atoi(rs.Primary.ID)
			opnCollection, _, err := client.NetworkingApi.OpnsGet(*auth, nil)
			if err != nil {
				return fmt.Errorf("Error waitiing OPN to be destroyed. Error: %v", err.Error())
			}
			_, err = findOpnById(int32(id), opnCollection)
			log.Printf("[DEBUG] OPN id %s", strconv.Itoa(id))
			if err != nil && strings.Contains(err.Error(), "is not present on list of OPNs") {
				return nil
			}
			return fmt.Errorf("Error destroying OPN. OPN was not destroyed")
		}
	}
	return fmt.Errorf("Error waiting OPN to be destroyed: Resource is not present on list of resources")
}

func testAccCheckOktawaveOpnConfig_basic(token string, mockStatus string) string {
	opn_name := "test_opn"
	if mockStatus == "1" {
		httpmock.Activate()
		mockPostOpn(136, opn_name)
		mockGetTicket(136, opn_name, `https://api.oktawave.com/beta//tickets/1`)
		mockGetOpn(opn_name, 1658)
		mockDeleteOpn(136, opn_name)
	}

	return fmt.Sprintf(`
provider "oktawave" {
  access_token="%s"
  
  api_url = "https://api.oktawave.com/beta/"
}

resource "oktawave_opn" "my_opn"{
	opn_name="%s"
}`, token, opn_name)
}

func testAccCheckOktawaveOpnConfig_updated(token string, mockStatus string) string {
	opn_name := "test_opn_upd"
	if mockStatus == "1" {
		mockPutOpn(opn_name)
	}
	return fmt.Sprintf(`
provider "oktawave" {
  access_token="%s"
  
  api_url = "https://api.oktawave.com/beta/"
}

resource "oktawave_opn" "my_opn"{
	opn_name="%s"
}`, token, opn_name)
}

func mockPostOpn(int32ticketStatusId int32, opnName string) {
	httpmock.RegisterResponder(http.MethodPost, "https://api.oktawave.com/beta//opns",
		func(req *http.Request) (*http.Response, error) {
			ticket := odk.Ticket{
				Id:            1,
				CreationDate:  time.Now(),
				CreationUser:  &(odk.UserResource{}),
				EndDate:       time.Now(),
				Status:        &(odk.DictionaryItem{int32ticketStatusId, "", &(odk.Resource{})}),
				OperationType: &(odk.DictionaryItem{}),
				ObjectId:      (int32)(1658),
				ObjectType:    &(odk.DictionaryItem{}),
				ObjectName:    opnName,
				Progress:      (int32)(100),
			}
			return httpmock.NewJsonResponse(200, ticket)
		})
}

func mockGetOpn(opn_name string, id int32) {
	httpmock.RegisterResponder(http.MethodGet, "https://api.oktawave.com/beta//opns",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, odk.ApiCollectionOpn{
				Items: []odk.Opn{
					{
						Id:             id,
						Name:           opn_name,
						CreationDate:   time.Now(),
						LastChangeDate: time.Now(),
						CreationUser:   &odk.UserResource{},
						PrivateIps:     make([]odk.PrivateIp, 0, 0),
					},
				},
			})
		})
}

func mockPutOpn(opn_name string) {
	httpmock.RegisterResponder(http.MethodPut, "https://api.oktawave.com/beta//opns/1658",
		func(req *http.Request) (*http.Response, error) {
			mockGetOpn(opn_name, 1658)
			return httpmock.NewJsonResponse(200, odk.Object{})
		})
}

func mockDeleteOpn(int32ticketStatusId int32, opnName string) {
	httpmock.RegisterResponder(http.MethodDelete, "https://api.oktawave.com/beta//opns/1658",
		func(req *http.Request) (*http.Response, error) {
			mockGetOpn(opnName, int32ticketStatusId)
			return httpmock.NewJsonResponse(200, odk.Ticket{
				Id:            1,
				CreationDate:  time.Now(),
				CreationUser:  &(odk.UserResource{}),
				EndDate:       time.Now(),
				Status:        &(odk.DictionaryItem{int32ticketStatusId, "", &(odk.Resource{})}),
				OperationType: &(odk.DictionaryItem{}),
				ObjectId:      (int32)(1658),
				ObjectType:    &(odk.DictionaryItem{}),
				ObjectName:    opnName,
				Progress:      (int32)(100),
			})
		})
}
