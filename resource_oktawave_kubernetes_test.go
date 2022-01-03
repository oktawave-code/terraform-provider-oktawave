package main

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/jarcoal/httpmock"
	swagger "github.com/oktawave-code/oks-sdk"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestAccOktawave_KubernetesCluster_Basic(t *testing.T) {
	var cluster swagger.K44SClusterDetailsDto
	mockStatus := os.Getenv("MOCK_STATUS")
	token := os.Getenv("TOKEN")
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOktawaveKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOktawaveKubernetesClusterConfig_basic(token, mockStatus),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOktawaveKubernetesClusterExists("oktawave_kubernetes_cluster.my_cluster", &cluster),
					testAccCheckOktawaveKubernetesClusterAttributes_basic(&cluster),
					resource.TestCheckResourceAttr("oktawave_kubernetes_cluster.my_cluster", "name", "tfclusr"),
					resource.TestCheckResourceAttr("oktawave_kubernetes_cluster.my_cluster", "version", "1.14.3"),
				),
				//ExpectError: regexp.MustCompile("errors during apply: Resource OCI. CREATE. Unable to create instance . Ticket status: 137"),
			},
		},
	})
	httpmock.DeactivateAndReset()
}

func testAccCheckOktawaveKubernetesClusterExists(name string, cluster *swagger.K44SClusterDetailsDto) resource.TestCheckFunc {
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
		id := rs.Primary.ID
		foundCluster, _, err := client.ClustersApi.ClustersNameGet(*auth, id)
		if err != nil {
			return fmt.Errorf("Cluster was not found by name: %s", rs.Primary.ID)
		}
		log.Printf("IP was found by id %s ", rs.Primary.ID)
		*cluster = foundCluster
		return nil
	}
}

func testAccCheckOktawaveKubernetesClusterAttributes_basic(cluster *swagger.K44SClusterDetailsDto) resource.TestCheckFunc {
	return func(t *terraform.State) error {
		if cluster.Version != "1.14.3" {
			return fmt.Errorf("Bad Cluster version. Expected: 1.14.3. Got: %v", cluster.Version)
		}
		return nil
	}
}

func testAccCheckOktawaveKubernetesClusterDestroy(s *terraform.State) error {
	if os.Getenv("MOCK_STATUS") == "1" {
		httpmock.RegisterNoResponder(httpmock.NewStringResponder(http.StatusNotFound, ""))
	}
	client := testAccProvider.Meta().(*ClientConfig).oktaOKSClient()
	auth := testAccProvider.Meta().(*ClientConfig).getOKSAuth()
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oktawave_kubernetes_cluster" {
			id := rs.Primary.ID
			_, resp, err := client.ClustersApi.ClustersNameGet(*auth, id)
			if err != nil && resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("Error waitiing IP to be destroyed. Error: %v", err.Error())
			}
			break
		}
	}
	return nil
}

func testAccCheckOktawaveKubernetesClusterConfig_basic(token string, mockStatus string) string {
	name := "tfclusr"
	version := "1.14.3"
	//comment := "example ip"
	if mockStatus == "1" {
		httpmock.Activate()
		mockClusterPost(name, version)
		mockClusterGet(name, version)
		mockClusterDelete(name, version)
	}
	return fmt.Sprintf(`
provider "oktawave" {
  access_token="%s"
  
  api_url = "https://api.oktawave.com/beta/"
}

resource "oktawave_kubernetes_cluster" "my_cluster" {
    name="%s"
    version="%s"
}`, token, name, version)
}

func mockClusterPost(name string, version string) {
	httpmock.RegisterResponder(http.MethodPost, fmt.Sprintf("https://k44s-api.i.k44s.oktawave.com/clusters/%s", name),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(http.StatusOK, swagger.K44SClusterDetailsDto{
				Name:         name,
				Version:      version,
				CreationDate: nil,
				Running:      true,
			})
		})
}

func mockClusterGet(name string, version string) {
	httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("https://k44s-api.i.k44s.oktawave.com/clusters/%s", name),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(http.StatusOK, swagger.K44SClusterDetailsDto{
				Name:         name,
				Version:      version,
				CreationDate: nil,
				Running:      true,
			})
		})
}

func mockClusterDelete(name string, version string) {
	httpmock.RegisterResponder(http.MethodDelete, fmt.Sprintf("https://k44s-api.i.k44s.oktawave.com/clusters/%s", name),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(http.StatusOK, swagger.K44SClusterDetailsDto{
				Name:         name,
				Version:      version,
				CreationDate: nil,
				Running:      true,
			})
		})
}
