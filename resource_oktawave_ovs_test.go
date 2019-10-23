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

func TestAccOktawave_Ovs_Basic(t *testing.T) {
	var volume odk.Disk
	token := os.Getenv("TOKEN")
	mockStatus := os.Getenv("MOCK_STATUS")
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOktawaveOvsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOktawaveOVSConfig_basic(token, mockStatus),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOktawaveOVSExists("oktawave_ovs.my_ovs", &volume),
					testAccCheckOktawaveOVSAttributes_basic(&volume),
					resource.TestCheckResourceAttr("oktawave_ovs.my_ovs", "disk_name", "my_disk2"),
					resource.TestCheckResourceAttr("oktawave_ovs.my_ovs", "space_capacity", "5"),
					resource.TestCheckResourceAttr("oktawave_ovs.my_ovs", "tier_id", "896"),
					resource.TestCheckResourceAttr("oktawave_ovs.my_ovs", "is_shared", "false"),
					resource.TestCheckResourceAttr("oktawave_ovs.my_ovs", "subregion_id", "4"),
					resource.TestCheckResourceAttr("oktawave_ovs.my_ovs", "is_locked", "false"),
					resource.TestCheckResourceAttr("oktawave_ovs.my_ovs", "isfreemium", "false"),
				),
				//ExpectError: regexp.MustCompile("errors during apply: Resource OCI. CREATE. Unable to create instance . Ticket status: 137"),
			},
		},
	})
}

func TestAccOktawave_Ovs_TierUpdate(t *testing.T) {
	var volume odk.Disk
	token := os.Getenv("TOKEN")
	mockStatus := os.Getenv("MOCK_STATUS")
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOktawaveOvsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOktawaveOVSConfig_basic(token, mockStatus),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOktawaveOVSExists("oktawave_ovs.my_ovs", &volume),
					testAccCheckOktawaveOVSAttributes_basic(&volume),
					resource.TestCheckResourceAttr("oktawave_ovs.my_ovs", "tier_id", "896"),
				),
			},
			{
				Config: testAccCheckOktawaveOVSConfig_TierUpdated(token, mockStatus),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOktawaveOVSExists("oktawave_ovs.my_ovs", &volume),
					testAccCheckOktawaveOVSAttributes_TierUpdate(&volume),
					resource.TestCheckResourceAttr("oktawave_ovs.my_ovs", "tier_id", "895"),
				),
			},
		},
	})
}

func testAccCheckOktawaveOVSExists(name string, volume *odk.Disk) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No volume ID is set")
		}

		client := testAccProvider.Meta().(*ClientConfig).oktaClient()
		auth := testAccProvider.Meta().(*ClientConfig).ctx
		id, _ := strconv.Atoi(rs.Primary.ID)
		foundVolume, _, err := client.OVSApi.DisksGet(*auth, int32(id), nil)
		if err != nil {
			return fmt.Errorf("Volume was not found by id: %s", rs.Primary.ID)
		}
		log.Printf("Volume was found by id %s and name %v", rs.Primary.ID, foundVolume.Name)
		*volume = foundVolume
		return nil
	}
}

func testAccCheckOktawaveOVSAttributes_basic(volume *odk.Disk) resource.TestCheckFunc {
	return func(t *terraform.State) error {
		if volume.Name != "my_disk2" {
			return fmt.Errorf("Bad OVS name. Expected: my_disk2. Got: %v", volume.Name)
		}
		if volume.Subregion.Id != 4 {
			volSubregionId := strconv.Itoa((int)(volume.Subregion.Id))
			return fmt.Errorf("Bad OVS subregion id. Expected: 4. Got: %v", volSubregionId)
		}
		if volume.IsFreemium != false {
			voltIsFreemium := strconv.FormatBool(volume.IsFreemium)
			return fmt.Errorf("Bad OVS freemium option. Expected: false. Got: %s", voltIsFreemium)
		}
		if volume.SpaceCapacity != 5 {
			volDiskSize := strconv.Itoa(int(volume.SpaceCapacity))
			return fmt.Errorf("Bad OVS disk capacity. Expected: 5. Got: %s", volDiskSize)
		}
		if volume.Tier.Id != 896 {
			return fmt.Errorf("Bad OVS tier id. Expected: 896. Got: %s", strconv.Itoa(int(volume.Tier.Id)))
		}
		if volume.IsShared != false {
			return fmt.Errorf("Bad OVS shared option. Expected: false. Got: true")
		}
		if volume.IsLocked != false {
			return fmt.Errorf("Bad OVS locked option. Expected: false. Got: true")
		}
		return nil
	}
}

func testAccCheckOktawaveOVSAttributes_TierUpdate(volume *odk.Disk) resource.TestCheckFunc {
	return func(t *terraform.State) error {
		if volume.Tier.Id != 895 {
			return fmt.Errorf("Bad OVS tier id. Expected: 896. Got: %s", strconv.Itoa(int(volume.Tier.Id)))
		}
		return nil
	}
}

//
////TODO API bug: throw 403 error when volume is not found
func testAccCheckOktawaveOvsDestroy(s *terraform.State) error {
	if os.Getenv("MOCK_STATUS") == "1" {
		httpmock.RegisterNoResponder(httpmock.NewStringResponder(403, "403"))
	}
	client := testAccProvider.Meta().(*ClientConfig).oktaClient()
	auth := testAccProvider.Meta().(*ClientConfig).ctx
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oktawave_ovs" {
			id, _ := strconv.Atoi(rs.Primary.ID)
			_, _, err := client.OVSApi.DisksGet(*auth, int32(id), nil)
			if err != nil && !strings.Contains(err.Error(), "403") {
				return fmt.Errorf("Error waitiing OVS to be destroyed. Error: %v", err.Error())
			}
			return nil
		}

	}
	return fmt.Errorf("Check destroy OVS  test function error: resource type was not found")
}

func testAccCheckOktawaveOVSConfig_basic(token string, mockStatus string) string {
	disk_name := "my_disk2"
	tier_id := 896
	subregion_id := 4
	if mockStatus == "1" {
		httpmock.Activate()
		mockPostOVS(136, disk_name)
		mockGetOVS(disk_name, int32(tier_id), int32(subregion_id))
		mockGetTicket(136, "my_disk2", `https://api.oktawave.com/beta/tickets/1`)
		mockDeleteOvs(136, disk_name)
		mockPutOVS(136, disk_name, int32(tier_id), int32(subregion_id))
	}
	return fmt.Sprintf(`
provider "oktawave" {
		access_token="%s"
		api_url = "https://api.oktawave.com/beta"
}

resource "oktawave_ovs" "my_ovs"{
	disk_name="%s"
	tier_id = %s
	subregion_id=%s
}`, token, disk_name, strconv.Itoa(tier_id), strconv.Itoa(subregion_id))
}

func testAccCheckOktawaveOVSConfig_TierUpdated(token string, mockStatus string) string {
	ovs_name := "my_disk2"
	tier_id := 895
	subregion_id := 4

	if mockStatus == "1" {
		mockPutOVS(136, ovs_name, int32(tier_id), int32(subregion_id))
	}
	return fmt.Sprintf(`
provider "oktawave" {
		access_token="%s"
		api_url = "https://api.oktawave.com/beta"
}

resource "oktawave_ovs" "my_ovs"{
	disk_name="%s"
	tier_id = %s
	subregion_id=%s
}`, token, ovs_name, strconv.Itoa(tier_id), strconv.Itoa(subregion_id))
}

func mockPostOVS(int32ticketStatusId int32, ovsName string) {
	httpmock.RegisterResponder(http.MethodPost, "https://api.oktawave.com/beta/disks",
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
				ObjectName:    ovsName,
				Progress:      (int32)(100),
			}
			return httpmock.NewJsonResponse(200, ticket)
		})
}

func mockGetOVS(ovs_name string, tier_id int32, subregion_id int32) {
	httpmock.RegisterResponder(http.MethodGet, "https://api.oktawave.com/beta/disks/1658",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, odk.Disk{
				Id:             1658,
				Name:           ovs_name,
				SpaceCapacity:  5,
				Tier:           &odk.DictionaryItem{Id: tier_id},
				CreationDate:   time.Now(),
				CreationUser:   &odk.UserResource{},
				IsShared:       false,
				SharedDiskType: &odk.DictionaryItem{},
				Subregion: &odk.BaseResource{
					Id: subregion_id,
				},
				IsLocked:    false,
				Connections: make([]odk.DiskConnection, 0, 0),
				IsFreemium:  false,
			})
		})
}

func mockPutOVS(int32ticketStatusId int32, ovs_name string, tier_id int32, subregion_id int32) {
	httpmock.RegisterResponder(http.MethodPut, "https://api.oktawave.com/beta/disks/1658",
		func(req *http.Request) (*http.Response, error) {
			mockGetOVS(ovs_name, tier_id, subregion_id)
			ticket := odk.Ticket{
				Id:            1,
				CreationDate:  time.Now(),
				CreationUser:  &(odk.UserResource{}),
				EndDate:       time.Now(),
				Status:        &(odk.DictionaryItem{int32ticketStatusId, "", &(odk.Resource{})}),
				OperationType: &(odk.DictionaryItem{}),
				ObjectId:      (int32)(1658),
				ObjectType:    &(odk.DictionaryItem{}),
				ObjectName:    ovs_name,
				Progress:      (int32)(100),
			}
			return httpmock.NewJsonResponse(200, ticket)
		})
}

func mockDeleteOvs(int32ticketStatusId int32, ovsName string) {
	httpmock.RegisterResponder(http.MethodDelete, "https://api.oktawave.com/beta/disks/1658",
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
				ObjectName:    ovsName,
				Progress:      (int32)(100),
			}
			return httpmock.NewJsonResponse(200, ticket)
		})
}
