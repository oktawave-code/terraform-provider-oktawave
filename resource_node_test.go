package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/jarcoal/httpmock"
	swagger "github.com/oktawave-code/oks-sdk"
)

func TestAccOktawave_Node_Basic(t *testing.T) {
	var node swagger.K44sInstance
	mockStatus := os.Getenv("MOCK_STATUS")
	token := os.Getenv("TOKEN")
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOktawaveNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOktawaveNodeConfig_basic(token, mockStatus),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOktawaveNodeExists("oktawave_kubernetes_node.my_node", &node),
					testAccCheckOktawaveNodeAttributes_basic(&node),
					resource.TestCheckResourceAttr("oktawave_kubernetes_node.my_node", "type_id", "1268"),
					resource.TestCheckResourceAttr("oktawave_kubernetes_node.my_node", "subregion_id", "4"),
				),
				//ExpectError: regexp.MustCompile("errors during apply: Resource OCI. CREATE. Unable to create instance . Ticket status: 137"),
			},
		},
	})
	httpmock.DeactivateAndReset()
}

func testAccCheckOktawaveNodeExists(name string, node *swagger.K44sInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ip ID is set")
		}

		client := testAccProvider.Meta().(*ClientConfig).oktaOKSClient()
		auth := testAccProvider.Meta().(*ClientConfig).getOKSAuth()
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error occured while converting id from string to int: %s", err)
		}
		clusterName := rs.Primary.Attributes["cluster_name"]
		nodes, _, err := client.ClustersApi.ClustersInstancesNameGet(*auth, clusterName)
		if err != nil {
			return fmt.Errorf("Cluster was not found by name: %s", clusterName)
		}

		foundInstance, err := retrieveNodeById(nodes, id)
		if err != nil {
			return fmt.Errorf(err.Error())
		}
		log.Printf("IP was found by id %s ", rs.Primary.ID)
		*node = foundInstance
		return nil
	}
}

func testAccCheckOktawaveNodeAttributes_basic(node *swagger.K44sInstance) resource.TestCheckFunc {
	return func(t *terraform.State) error {
		if node.Type_.Id != float64(1268) {
			return fmt.Errorf("Bad Node type id. Expected: 1268. Got: %v", strconv.Itoa(int(node.Type_.Id)))
		}
		if node.Subregion.Id != float64(4) {
			return fmt.Errorf("Bad Node type id. Expected: 4. Got: %v", strconv.Itoa(int(node.Subregion.Id)))
		}
		return nil
	}
}

func testAccCheckOktawaveNodeDestroy(s *terraform.State) error {
	if os.Getenv("MOCK_STATUS") == "1" {
		httpmock.RegisterNoResponder(httpmock.NewStringResponder(http.StatusNotFound, ""))
	}
	client := testAccProvider.Meta().(*ClientConfig).oktaOKSClient()
	auth := testAccProvider.Meta().(*ClientConfig).getOKSAuth()
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oktawave_kubernetes_node" {
			id, err := strconv.Atoi(rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("Error occured while converting id from string to int: %s", err)
			}
			clusterName := rs.Primary.Attributes["cluster_name"]
			nodes, _, err := client.ClustersApi.ClustersInstancesNameGet(*auth, clusterName)
			if err != nil {
				return fmt.Errorf("Cluster was not found by name: %s", clusterName)
			}

			_, err = retrieveNodeById(nodes, id)
			if err == nil {
				return fmt.Errorf("Error deleting instance: node still exists")
			}
			break
		}
	}
	return nil
}

func testAccCheckOktawaveNodeConfig_basic(token string, mockStatus string) string {
	cluster_name := "tfclusr"
	cluster_version := "1.15.0"
	node_type_id := 1268
	node_subregion_id := 4
	//comment := "example ip"
	if mockStatus == "1" {
		httpmock.Activate()
		mockClusterPost(cluster_name, cluster_version)
		mockClusterGet(cluster_name, cluster_version)
		mockClusterDelete(cluster_name, cluster_version)
		mockNodePost(cluster_name, float32(node_subregion_id), float32(node_type_id))
		// mockGetTask(cluster_name, "CREATE", float32(node_subregion_id), float32(node_type_id))
		mockGetTicket(TICKET_STATUS__SUCCESS, "name", `https://api.oktawave.com/beta/tickets/1`)
		mockGetNode(cluster_name, float64(node_subregion_id), float64(node_type_id))
		mockNodeDelete(cluster_name, float32(node_subregion_id), float32(node_type_id))
	}
	return fmt.Sprintf(`
provider "oktawave" {
  access_token="%s"
  
  api_url = "https://api.oktawave.com/beta"
}

resource "oktawave_kubernetes_cluster" "my_cluster" {
    name="%s"
    version="%s"
}

resource "oktawave_kubernetes_node" "my_node"{
	type_id=%s
	subregion_id=%s
	cluster_name=oktawave_kubernetes_cluster.my_cluster.id
	depends_on=[oktawave_kubernetes_cluster.my_cluster]
}`, token, cluster_name, cluster_version, strconv.Itoa(node_type_id), strconv.Itoa(node_subregion_id))
}

func mockNodePost(name string, subregionId float32, typeId float32) {
	httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("https://k44s-api.i.k44s.oktawave.com/clusters/instances/%s", name),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(http.StatusOK, []swagger.K44SOperation{
				swagger.K44SOperation{
					Ticket: &swagger.IaasTicket{
						Id:           float64(1),
						CreationDate: time.Now(),
						EndDate:      time.Now(),
						Status: &swagger.IaasDictionaryEntry{
							Id:    float64(TICKET_STATUS__SUCCESS),
							Label: "",
						},
						OperationType: &(swagger.IaasDictionaryEntry{}),
						ObjectId:      100,
						ObjectName:    "test",
						Progress:      100,
					},
					Error_: "",
				},
			})
		})
}

// func mockGetTask(cluster_name string, taskOperation string, subregionId float32, typeId float32) {
// 	httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("https://k44s-api.i.k44s.oktawave.com/clusters/instances/%s/tasks/1", cluster_name),
// 		func(req *http.Request) (*http.Response, error) {
// 			return httpmock.NewJsonResponse(http.StatusOK, swagger.K44sTaskDto{
// 				TaskId:       "1",
// 				Operation:    taskOperation,
// 				InstanceName: "cluster_node",
// 				SubregionId:  subregionId,
// 				TypeId:       typeId,
// 				InstanceId:   1,
// 				Status:       "Succeeded",
// 			})
// 		})
// }

func mockGetNode(name string, subregionId float64, typeId float64) {
	httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("https://k44s-api.i.k44s.oktawave.com/clusters/instances/%s", name),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(http.StatusOK, []swagger.K44sInstance{
				{
					Id:           1,
					Name:         "name",
					CreationDate: "",
					Subregion: &swagger.K44sInstanceSubregion{
						Id: subregionId,
					},
					Type_: &swagger.K44sInstanceType{
						Id: typeId,
					},
					Status:             &swagger.K44sInstanceStatus{},
					IpAddress:          "1.1.1.1",
					TotalDisksCapacity: float64(10),
					CpuNumber:          float64(10),
					RamMb:              float64(10),
				},
			})
		})
}

func mockNodeDelete(cluster_name string, subregionId float32, typeId float32) {
	httpmock.RegisterResponder(http.MethodDelete, fmt.Sprintf("https://k44s-api.i.k44s.oktawave.com/clusters/instances/%s", cluster_name),
		func(req *http.Request) (*http.Response, error) {
			// mockGetTask(cluster_name, "DELETE", subregionId, typeId)
			mockGetTicket(TICKET_STATUS__SUCCESS, "name", `https://api.oktawave.com/beta/tickets/1`)
			httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("https://k44s-api.i.k44s.oktawave.com/clusters/instances/%s", cluster_name),
				func(req *http.Request) (*http.Response, error) {
					return httpmock.NewJsonResponse(http.StatusOK, []swagger.K44sInstance{})
				})
			return httpmock.NewJsonResponse(http.StatusOK, []swagger.K44SOperation{
				{
					Ticket: &swagger.IaasTicket{
						Id:           1,
						CreationDate: time.Now(),
						EndDate:      time.Now(),
						Status: &swagger.IaasDictionaryEntry{
							Id:    float64(TICKET_STATUS__SUCCESS),
							Label: "",
						},
						OperationType: &(swagger.IaasDictionaryEntry{}),
						ObjectId:      100,
						ObjectName:    "test",
						Progress:      100,
					},
					Error_: "",
				},
			})
		})
}
