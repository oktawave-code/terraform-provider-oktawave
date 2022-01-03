package main

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/oktawave-code/odk"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

//TODO: CANNOT MOCK UPDATE TEST: WHEN MOCKING 2 GETS FOR THE SAME ID - THE LAST ONE IMMEDIEATELY REPLACE FIRST ONE
//TODO: TRIED TO TRICK THE SYSTEM: REPLACE IDS IN MOCKS. BUT IN THAT CASE I CANNOT MAKE LOCAL STATE TEST AND STATE DRIFTS APPEARS
//TODO: SYSTEM STORE THE STATE OF INFRASTRUCTURE CREATED IN CREATE METHOD. AS A RESULT, IN UPDATE TEST I TRY TO CALLBACK TO INFRASTRUCTURE
//TODO: THAT IS NOT STORED BY TERRAFORM AND WHEN LOCAL STATE TEST BEGIN I GET ERRORS
func TestAccOktawave_Oci_Basic(t *testing.T) {
	var instance odk.Instance
	mockStatus := os.Getenv("MOCK_STATUS")
	token := os.Getenv("TOKEN")
	//TODO: make different creationg and reading scenarios using environment variables
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOktawaveOciDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOktawaveOCIConfig_basic(mockStatus, token),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOktawaveOCIExists("oktawave_oci.my_oci", &instance),
					testAccCheckOktawaveOCIAttributes(&instance),
					resource.TestCheckResourceAttr("oktawave_oci.my_oci", "authorization_method_id", "1399"),
					resource.TestCheckResourceAttr("oktawave_oci.my_oci", "disk_class", "896"),
					resource.TestCheckResourceAttr("oktawave_oci.my_oci", "init_disk_size", "5"),
					resource.TestCheckResourceAttr("oktawave_oci.my_oci", "instance_name", "my_instance2"),
					resource.TestCheckResourceAttr("oktawave_oci.my_oci", "instances_count", "1"),
					resource.TestCheckResourceAttr("oktawave_oci.my_oci", "subregion_id", "4"),
					resource.TestCheckResourceAttr("oktawave_oci.my_oci", "template_id", "94"),
					resource.TestCheckResourceAttr("oktawave_oci.my_oci", "type_id", "1268"),
					resource.TestCheckResourceAttr("oktawave_oci.my_oci", "isfreemium", "false"),
				),
				//ExpectError: regexp.MustCompile("errors during apply: Resource OCI. CREATE. Unable to create instance . Ticket status: 137"),
			},
		},
	})

}

func TestAccOktawave_Oci_Update(t *testing.T) {
	var instance odk.Instance
	token := os.Getenv("TOKEN")
	mockStatus := os.Getenv("MOCK_STATUS")
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOktawaveOciDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOktawaveOCIConfig_basic(mockStatus, token),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOktawaveOCIExists("oktawave_oci.my_oci", &instance),
					testAccCheckOktawaveOCIAttributes(&instance),
					resource.TestCheckResourceAttr("oktawave_oci.my_oci", "instance_name", "my_instance2"),
					resource.TestCheckResourceAttr("oktawave_oci.my_oci", "type_id", "1268"),
				),
			},
			{
				Config: testAccCheckOktawaveOCIConig_RenameAndRetyped(mockStatus, token),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOktawaveOCIExists("oktawave_oci.my_oci", &instance),
					testAccCheckOktawaveOCI_RenameAndRetyped(&instance),
					resource.TestCheckResourceAttr("oktawave_oci.my_oci", "instance_name", "my_instance1"),
					resource.TestCheckResourceAttr("oktawave_oci.my_oci", "type_id", "1270"),
				),
			},
		},
	})
	httpmock.DeactivateAndReset()
}

func testAccCheckOktawaveOCIExists(name string, instance *odk.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}
		res, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}
		id := int32(res)
		client := testAccProvider.Meta().(*ClientConfig).oktaClient()
		auth := testAccProvider.Meta().(*ClientConfig).ctx
		foundInstance, _, err := client.OCIApi.InstancesGet_2(*auth, id, nil)
		if err != nil {
			return fmt.Errorf("Instance not found %s ", err)
		}
		log.Printf("[INFO]Test OCI. Found instance name: %s", foundInstance.Name)
		*instance = foundInstance
		return nil
	}
}

//Check whether OCI attributes correspond to expected value
func testAccCheckOktawaveOCIAttributes(instance *odk.Instance) resource.TestCheckFunc {
	return func(t *terraform.State) error {
		if instance.Name != "my_instance2" {
			return fmt.Errorf("Bad Basic OCI name. Expected: my_instance2. Got: %v", instance.Name)
		}
		if instance.Subregion.Id != 4 {
			instSubregionId := strconv.Itoa((int)(instance.Subregion.Id))
			return fmt.Errorf("Bad Basic OCI subregion id. Expected: 4. Got: %v", instSubregionId)
		}
		if instance.Template.Id != 94 {
			instTemplateId := strconv.Itoa((int)(instance.Template.Id))
			return fmt.Errorf("Bad OCI template id. Expected: 94. Got: %s", instTemplateId)
		}
		if instance.Type_.Id != 1268 {
			instTypeId := strconv.Itoa((int)(instance.Type_.Id))
			return fmt.Errorf("Bad Basic OCI type id. Expected: 1268. Got: %s", instTypeId)
		}
		if instance.IsFreemium != false {
			instIsFreemium := strconv.FormatBool(instance.IsFreemium)
			return fmt.Errorf("Bad Basic OCI freemium option. Expected: false. Got: %s", instIsFreemium)
		}
		if instance.TotalDisksCapacity != 5 {
			instDiskSize := strconv.Itoa((int)(instance.TotalDisksCapacity))
			return fmt.Errorf("Bad Basic OCI disk capacity. Expected: 5. Got: %s", instDiskSize)
		}
		return nil
	}
}

func testAccCheckOktawaveOCI_RenameAndRetyped(instance *odk.Instance) resource.TestCheckFunc {
	return func(t *terraform.State) error {
		log.Printf("Inst name %s", instance.Name)
		if instance.Name != "my_instance1" {
			return fmt.Errorf("Bad  Updated OCI name. Expected: my_instance1. Got: %v", instance.Name)
		}

		if instance.Type_.Id != 1270 {
			instTypeId := strconv.Itoa((int)(instance.Type_.Id))
			return fmt.Errorf("Bad Updated OCI type id. Expected: 1270. Got: %s", instTypeId)
		}
		return nil
	}
}

func testAccCheckOktawaveOciDestroy(s *terraform.State) error {
	if os.Getenv("MOCK_STATUS") == "1" {
		httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, "Instance not found"))
	}

	client := testAccProvider.Meta().(*ClientConfig).oktaClient()
	auth := testAccProvider.Meta().(*ClientConfig).ctx
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oktawave_oci" {
			id, err := strconv.Atoi(rs.Primary.ID)
			if err != nil {
				return err
			}

			//finding instance
			_, resp, err := client.OCIApi.InstancesGet_2(*auth, (int32)(id), nil)
			if err != nil && resp.StatusCode != 404 {
				return fmt.Errorf("Error waitiing OCI to be destroyed: %s", rs.Primary.ID)
			}
			return nil
		}
	}
	return fmt.Errorf("Check destroy OCI  test function error: resource type was not found")
}

//
////Resource Config and eventually mock http operations for this resource
func testAccCheckOktawaveOCIConfig_basic(mockStatus string, token string) string {
	//	//Commented parts - for mock
	if mockStatus == "1" {
		name := "my_instance2"
		var ticketStatusId int32 = 136

		httpmock.Activate()
		mockOCIGet("my_instance2", 94, 4, 1268, 5, false)
		mockOCIPost(ticketStatusId, name)
		mockGetTicket(ticketStatusId, name, `https://api.oktawave.com/beta/tickets/1`)
		mockOCIGetDisk()
		mockOCIGetDisks(1658)
		mockOCIDelete(ticketStatusId, name)
		mockGetOciOpns()

	}

	return fmt.Sprintf(
		`provider "oktawave" {
		access_token="%s"
		api_url = "https://api.oktawave.com/beta"
	}


	resource "oktawave_oci" "my_oci" {
		disk_class =896
		init_disk_size = 5
		instance_name ="my_instance2"
		subregion_id =4
		template_id =94
		type_id = 1268
	}`, token)
}

//
func testAccCheckOktawaveOCIConig_RenameAndRetyped(mockStatus string, token string) string {
	//	//Commented parts - for mock
	//
	if mockStatus == "1" {
		name := "my_instance1"
		//var diskClass int32= 896
		var ticketStatusId int32 = 136

		mockOCIChangeNameTicket(ticketStatusId, name)
		mockOCIChangeTypeTicket(ticketStatusId, name)
		mockInitDiskUPD(ticketStatusId, name)
	}

	return fmt.Sprintf(`provider "oktawave" {
		access_token="%s"
		api_url = "https://api.oktawave.com/beta"
	}

	resource "oktawave_oci" "my_oci" {
		disk_class =896
		init_disk_size = 5
		instance_name ="my_instance1"
		subregion_id =4
		template_id =94
		type_id = 1270
	}`, token)
}

func mockOCIPost(int32ticketStatusId int32, name string) {
	httpmock.RegisterResponder("POST", "https://api.oktawave.com/beta/instances",
		func(req *http.Request) (*http.Response, error) {
			ticket := map[string]interface{}{
				"Id":            1,
				"CreationDate":  time.Now(),
				"CreationUser":  &(odk.UserResource{}),
				"EndDate":       time.Now(),
				"Status":        &(odk.DictionaryItem{int32ticketStatusId, "", &(odk.Resource{})}),
				"OperationType": &(odk.DictionaryItem{}),
				"ObjectId":      (int32)(1658),
				"ObjectType":    &(odk.DictionaryItem{}),
				"ObjectName":    name,
				"Progress":      (int32)(100),
			}
			if err := json.NewDecoder(req.Body).Decode(&ticket); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			resp, err := httpmock.NewJsonResponse(200, ticket)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		})
}

func mockGetTicket(int32ticketStatusId int32, name string, url string) {
	httpmock.RegisterResponder("GET", url,
		func(req *http.Request) (*http.Response, error) {

			return httpmock.NewJsonResponse(200, map[string]interface{}{
				"Id":            1,
				"CreationDate":  time.Now(),
				"CreationUser":  &(odk.UserResource{}),
				"EndDate":       time.Now(),
				"Status":        &(odk.DictionaryItem{int32ticketStatusId, "", &(odk.Resource{})}),
				"OperationType": &(odk.DictionaryItem{}),
				"ObjectId":      (int32)(1658),
				"ObjectType":    &(odk.DictionaryItem{}),
				"ObjectName":    name,
				"Progress":      (int32)(100),
			})
		})
}

func mockOCIGet(name string, template_id int32, subRegionId int32, type_id int32, diskSize int32, isFreemium bool) {
	apiLink := fmt.Sprintf("https://api.oktawave.com/beta/instances/1658")
	httpmock.RegisterResponder("GET", apiLink,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, map[string]interface{}{
				"Id":                 (int32)(1658),
				"Name":               name,
				"CreationDate":       time.Now(),
				"CreationUser":       &(odk.UserResource{}),
				"IsLocked":           false,
				"LockingDate":        time.Now(),
				"Template":           &(odk.BaseResource{template_id}),
				"Subregion":          &(odk.BaseResource{subRegionId}),
				"Type":               &(odk.DictionaryItem{type_id, "", &(odk.Resource{})}),
				"Status":             &(odk.DictionaryItem{100, "", &(odk.Resource{})}),
				"SystemCategory":     &(odk.DictionaryItem{100, "", &(odk.Resource{})}),
				"AutoscalingType":    &(odk.DictionaryItem{100, "", &(odk.Resource{})}),
				"VmWareToolsStatus":  &(odk.DictionaryItem{100, "", &(odk.Resource{})}),
				"MonitStatus":        &(odk.DictionaryItem{100, "", &(odk.Resource{})}),
				"TemplateType":       &(odk.DictionaryItem{template_id, "", &(odk.Resource{})}),
				"IpAddress":          "",
				"DnsAddress":         "",
				"PaymentType":        &(odk.DictionaryItem{template_id, "", &(odk.Resource{})}),
				"HealthCheck":        &(odk.BaseResource{}),
				"ScsiControllerType": &(odk.DictionaryItem{template_id, "", &(odk.Resource{})}),
				"CpuNumber":          (int32)(24),
				"RamMb":              (int32)(8196),
				"SupportType":        &(odk.Software{}),
				"TotalDisksCapacity": diskSize,
				"IsFreemium":         isFreemium,
			})
		})
}

func mockOCIDelete(int32ticketStatusId int32, name string) {

	httpmock.RegisterResponder(http.MethodDelete, "https://api.oktawave.com/beta/instances/1658",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, map[string]interface{}{
				"Id":            1,
				"CreationDate":  time.Now(),
				"CreationUser":  &(odk.UserResource{}),
				"EndDate":       time.Now(),
				"Status":        &(odk.DictionaryItem{int32ticketStatusId, "", &(odk.Resource{})}),
				"OperationType": &(odk.DictionaryItem{}),
				"ObjectId":      (int32)(1658),
				"ObjectType":    &(odk.DictionaryItem{}),
				"ObjectName":    name,
				"Progress":      (int32)(100),
			})
		})
}

func mockOCIGetDisks(instanceId int) {
	apiLink := fmt.Sprintf("https://api.oktawave.com/beta/instances/%s/disks", strconv.Itoa(instanceId))
	httpmock.RegisterResponder("GET", apiLink,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, map[string]interface{}{
				"Items": []odk.Disk{
					{
						Id:            1000,
						Name:          "",
						SpaceCapacity: 5,
						Tier: &odk.DictionaryItem{
							Id: 896,
						},
						CreationDate:   time.Now(),
						CreationUser:   &odk.UserResource{},
						IsShared:       false,
						SharedDiskType: &odk.DictionaryItem{},
						Subregion: &odk.BaseResource{
							Id: 4,
						},
						IsLocked:    false,
						LockingDate: time.Now(),
						Connections: []odk.DiskConnection{},
						IsFreemium:  false,
					},
				},
			})
		})
}

func mockOCIGetDisk() {
	httpmock.RegisterResponder("GET", "https://api.oktawave.com/beta/disks/1000",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, map[string]interface{}{
				"Id":            1000,
				"Name":          "",
				"SpaceCapacity": 5,
				"Tier": &odk.DictionaryItem{
					Id: 896,
				},
				"CreationDate":   time.Now(),
				"CreationUser":   &odk.UserResource{},
				"IsShared":       false,
				"SharedDiskType": &odk.DictionaryItem{},
				"Subregion": &odk.BaseResource{
					Id: 4,
				},
				"IsLocked":    false,
				"LockingDate": time.Now(),
				"Connections": []odk.DiskConnection{},
				"IsFreemium":  false,
			})
		})
}

func mockOCIChangeNameTicket(int32ticketStatusId int32, name string) {
	httpmock.RegisterResponder("POST", "https://api.oktawave.com/beta/instances/1658/change_name_ticket",
		func(req *http.Request) (*http.Response, error) {
			mockOCIGet(name, 94, 4, 1268, 5, false)
			ticket := map[string]interface{}{
				"Id":            1,
				"CreationDate":  time.Now(),
				"CreationUser":  &(odk.UserResource{}),
				"EndDate":       time.Now(),
				"Status":        &(odk.DictionaryItem{int32ticketStatusId, "", &(odk.Resource{})}),
				"OperationType": &(odk.DictionaryItem{}),
				"ObjectId":      (int32)(1658),
				"ObjectType":    &(odk.DictionaryItem{}),
				"ObjectName":    name,
				"Progress":      (int32)(100),
			}
			return httpmock.NewJsonResponse(200, ticket)
		})
}

func mockOCIChangeTypeTicket(int32ticketStatusId int32, name string) {
	httpmock.RegisterResponder("POST", "https://api.oktawave.com/beta/instances/1658/change_type_ticket",
		func(req *http.Request) (*http.Response, error) {
			mockOCIGet(name, 94, 4, 1270, 5, false)
			ticket := map[string]interface{}{
				"Id":            1,
				"CreationDate":  time.Now(),
				"CreationUser":  &(odk.UserResource{}),
				"EndDate":       time.Now(),
				"Status":        &(odk.DictionaryItem{int32ticketStatusId, "", &(odk.Resource{})}),
				"OperationType": &(odk.DictionaryItem{}),
				"ObjectId":      (int32)(1658),
				"ObjectType":    &(odk.DictionaryItem{}),
				"ObjectName":    name,
				"Progress":      (int32)(100),
			}
			return httpmock.NewJsonResponse(200, ticket)
		})
}

func mockInitDiskUPD(int32ticketStatusId int32, name string) {
	httpmock.RegisterResponder("PUT", "https://api.oktawave.com/beta/disks/1000",
		func(req *http.Request) (*http.Response, error) {
			ticket := map[string]interface{}{
				"Id":            1,
				"CreationDate":  time.Now(),
				"CreationUser":  &(odk.UserResource{}),
				"EndDate":       time.Now(),
				"Status":        &(odk.DictionaryItem{int32ticketStatusId, "", &(odk.Resource{})}),
				"OperationType": &(odk.DictionaryItem{}),
				"ObjectId":      (int32)(1658),
				"ObjectType":    &(odk.DictionaryItem{}),
				"ObjectName":    name,
				"Progress":      (int32)(100),
			}
			return httpmock.NewJsonResponse(200, ticket)
		})
}

func mockGetOciOpns() {
	httpmock.RegisterResponder(http.MethodGet, "https://api.oktawave.com/beta/instances/1658/opns",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, odk.ApiCollectionOpn{})
		})
}
