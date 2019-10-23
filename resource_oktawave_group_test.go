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

func TestAccOktawave_Group_Basic(t *testing.T) {
	var group odk.Group
	token := os.Getenv("TOKEN")
	mockStatus := os.Getenv("MOCK_STATUS")
	//Environment variable to manage sticket status
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOktawaveGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOktawaveGroupConfig_basic(token, mockStatus),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOktawaveGroupExists("oktawave_group.my_group", &group),
					testAccCheckOktawaveGroupAttributes_basic(&group),
					resource.TestCheckResourceAttr("oktawave_group.my_group", "group_name", "my_group"),
					resource.TestCheckResourceAttr("oktawave_group.my_group", "affinity_rule_type_id", "1403"),
					resource.TestCheckResourceAttr("oktawave_group.my_group", "instances_count", "0"),
				),
				//ExpectError: regexp.MustCompile("errors during apply: Resource OCI. CREATE. Unable to create instance . Ticket status: 137"),
			},
		},
	})
	httpmock.DeactivateAndReset()
}

func TestAccOktawave_Group_Updated(t *testing.T) {
	var group odk.Group
	token := os.Getenv("TOKEN")
	mockStatus := os.Getenv("MOCK_STATUS")
	//Environment variable to manage sticket status
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOktawaveGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOktawaveGroupConfig_basic(token, mockStatus),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOktawaveGroupExists("oktawave_group.my_group", &group),
					testAccCheckOktawaveGroupAttributes_basic(&group),
					resource.TestCheckResourceAttr("oktawave_group.my_group", "group_name", "my_group"),
					resource.TestCheckResourceAttr("oktawave_group.my_group", "affinity_rule_type_id", "1403"),
					resource.TestCheckResourceAttr("oktawave_group.my_group", "instances_count", "0"),
				),
				//ExpectError: regexp.MustCompile("errors during apply: Resource OCI. CREATE. Unable to create instance . Ticket status: 137"),
			},
			{
				Config: testAccCheckOktawaveGroupConfig_updated(token, mockStatus),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOktawaveGroupExists("oktawave_group.my_group", &group),
					testAccCheckOktawaveGroupAttributes_updated(&group),
					resource.TestCheckResourceAttr("oktawave_group.my_group", "group_name", "my_group1"),
					resource.TestCheckResourceAttr("oktawave_group.my_group", "affinity_rule_type_id", "1404"),
					resource.TestCheckResourceAttr("oktawave_group.my_group", "instances_count", "0"),
				),
				//ExpectError: regexp.MustCompile("errors during apply: Resource OCI. CREATE. Unable to create instance . Ticket status: 137"),
			},
		},
	})
	httpmock.DeactivateAndReset()
}

func testAccCheckOktawaveGroupExists(name string, group *odk.Group) resource.TestCheckFunc {
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
		foundInstance, _, err := client.OCIGroupsApi.GroupsGetGroup(*auth, (int32)(id), nil)
		if err != nil {
			return fmt.Errorf("Instance not found %s ", err)
		}
		log.Printf("[INFO]Test Group. Found group name: %s", foundInstance.Name)
		*group = foundInstance
		return nil
	}
}

func testAccCheckOktawaveGroupAttributes_basic(group *odk.Group) resource.TestCheckFunc {
	return func(t *terraform.State) error {
		if group.Name != "my_group" {
			return fmt.Errorf("Bad Group name. Expected: my_group. Got: %v", group.Name)
		}
		if group.AffinityRuleType.Id != 1403 {
			groupAffinityRuleType := strconv.Itoa((int)(group.AffinityRuleType.Id))
			return fmt.Errorf("Bad Group affinity rule type id. Expected: 1403. Got: %s", groupAffinityRuleType)
		}
		return nil
	}
}

func testAccCheckOktawaveGroupAttributes_updated(group *odk.Group) resource.TestCheckFunc {
	return func(t *terraform.State) error {
		if group.Name != "my_group1" {
			return fmt.Errorf("Bad Group name. Expected: my_group1. Got: %v", group.Name)
		}
		if group.AffinityRuleType.Id != 1404 {
			groupAffinityRuleType := strconv.Itoa((int)(group.AffinityRuleType.Id))
			return fmt.Errorf("Bad Group affinity rule type id. Expected: 1404. Got: %s", groupAffinityRuleType)
		}
		return nil
	}
}

func testAccCheckOktawaveGroupDestroy(s *terraform.State) error {
	if os.Getenv("MOCK_STATUS") == "1" {
		httpmock.RegisterNoResponder(httpmock.NewStringResponder(200, ""))
	}
	client := testAccProvider.Meta().(*ClientConfig).oktaClient()
	auth := testAccProvider.Meta().(*ClientConfig).ctx
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oktawave_group" {
			id, _ := strconv.Atoi(rs.Primary.ID)
			_, resp, err := client.OCIGroupsApi.GroupsGetGroup(*auth, int32(id), nil)
			if err != nil && resp.StatusCode != 200 {
				return fmt.Errorf("Error waitiing OVS to be destroyed. Error: %v", err.Error())
			}
			return nil
		}

	}
	return fmt.Errorf("Check destroy Group  test function error: resource type was not found")
}

func testAccCheckOktawaveGroupConfig_basic(token string, mockStatus string) string {
	name := "my_group"
	isLoadBalancer := false
	var affinitRuleTypeId int32 = 1403

	if mockStatus == "1" {
		httpmock.Activate()
		mockGroupPost(name, isLoadBalancer, affinitRuleTypeId)
		mockChangeGroupAssignments_empty()
		mockGetGroup(name, isLoadBalancer, affinitRuleTypeId)
		mockGetGroupAssignments_empty()
		mockDeleteGroup()
	}

	return fmt.Sprintf(`
provider "oktawave" {
 access_token="%s"

 api_url = "https://api.oktawave.com/beta/"
}


resource "oktawave_group" "my_group"{
	group_name="%s"
	affinity_rule_type_id=%s
}`, token, name, strconv.Itoa(int(affinitRuleTypeId)))
}

func testAccCheckOktawaveGroupConfig_updated(token string, mockStatus string) string {
	name := "my_group1"
	isLoadBalancer := false
	var affinityRuleTypeId int32 = 1404
	if os.Getenv("MOCK_STATUS") == "1" {
		mockGroupPut(name, isLoadBalancer, affinityRuleTypeId)
	}
	return fmt.Sprintf(`
provider "oktawave" {
 access_token="%s"

 api_url = "https://api.oktawave.com/beta/"
}


resource "oktawave_group" "my_group"{
	group_name="my_group1"
	affinity_rule_type_id=1404
}`, token)
}

func mockGroupPost(name string, isLoadBalancer bool, affinityRuleTypeId int32) {
	httpmock.RegisterResponder(http.MethodPost, "https://api.oktawave.com/beta//groups",
		func(req *http.Request) (*http.Response, error) {
			group := odk.Group{
				Id:              1658,
				Name:            name,
				IsLoadBalancer:  isLoadBalancer,
				InstancesCount:  0,
				SchedulersCount: 0,
				AffinityRuleType: &odk.DictionaryItem{
					Id: affinityRuleTypeId,
				},
				AutoscalingType: &odk.DictionaryItem{},
				LastChangeDate:  time.Now(),
				CreationUser:    &odk.UserResource{},
			}
			return httpmock.NewJsonResponse(200, group)
		})
}

func mockChangeGroupAssignments_empty() {
	httpmock.RegisterResponder(http.MethodPut, "https://api.oktawave.com/beta//groups/1658/assignments",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, odk.ApiCollectionGroupAssignment{
				Items: []odk.GroupAssignment{},
				Meta:  &odk.ApiCollectionMetadata{},
			})
		})
}

func mockGetGroupAssignments_empty() {
	httpmock.RegisterResponder(http.MethodGet, "https://api.oktawave.com/beta//groups/1658/assignments",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, odk.ApiCollectionGroupAssignment{
				Items: []odk.GroupAssignment{},
				Meta:  &odk.ApiCollectionMetadata{},
			})
		})
}

func mockGetGroup(name string, isLoadBalancer bool, affinityRuleTypeId int32) {
	httpmock.RegisterResponder(http.MethodGet, "https://api.oktawave.com/beta//groups/1658",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, odk.Group{
				Id:              1658,
				Name:            name,
				IsLoadBalancer:  isLoadBalancer,
				InstancesCount:  0,
				SchedulersCount: 0,
				AffinityRuleType: &odk.DictionaryItem{
					Id: affinityRuleTypeId,
				},
				AutoscalingType: &odk.DictionaryItem{},
				LastChangeDate:  time.Now(),
				CreationUser:    &odk.UserResource{},
			})
		})
}

func mockDeleteGroup() {
	httpmock.RegisterResponder(http.MethodDelete, "https://api.oktawave.com/beta//groups/1658",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, odk.Object{})
		})
}

func mockGroupPut(name string, isLoadBalancer bool, affinityRuleTypeId int32) {
	httpmock.RegisterResponder(http.MethodPut, "https://api.oktawave.com/beta//groups/1658",
		func(req *http.Request) (*http.Response, error) {
			mockGetGroup(name, isLoadBalancer, affinityRuleTypeId)
			group := odk.Group{
				Id:              1658,
				Name:            name,
				IsLoadBalancer:  isLoadBalancer,
				InstancesCount:  0,
				SchedulersCount: 0,
				AffinityRuleType: &odk.DictionaryItem{
					Id: affinityRuleTypeId,
				},
				AutoscalingType: &odk.DictionaryItem{},
				LastChangeDate:  time.Now(),
				CreationUser:    &odk.UserResource{},
			}
			return httpmock.NewJsonResponse(200, group)
		})
}
