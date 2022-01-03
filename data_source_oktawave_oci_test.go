package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/oktawave-code/odk"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestAccDataSourceOktawave_Oci_Basic(t *testing.T) {
	//var instance odk.Instance
	mockStatus := os.Getenv("MOCK_STATUS")
	token := os.Getenv("TOKEN")
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceOktawaveOCIConfig_basic(mockStatus, token),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.oktawave_oci.oci",
						"id",
						regexp.MustCompile("^\\d+$"),
					),
				),
			},
		},
	})
}

func testAccCheckDataSourceOktawaveOCIConfig_basic(mockStatus string, token string) string {
	if mockStatus == "1" {
		httpmock.Activate()
		mockOCIGetAll("my_instance3", 4, 94, 1268)
	}

	return fmt.Sprintf(
		`provider "oktawave" {
		access_token="%s"
		api_url = "https://api.oktawave.com/beta"
	}

	data "oktawave_oci" "oci" {
		instance_name = "my_instance3"
		subregion_id = 4
		template_id = 94
		type_id = 1268
	}
	`, token)
}

func mockOCIGetAll(instance_name string, subregion_id int32, template_id int32, type_id int32) {
	apiLink := fmt.Sprintf("https://api.oktawave.com/beta/instances")
	httpmock.RegisterResponder("GET", apiLink,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, map[string]interface{}{
				"Items": []odk.Instance{
					{
						Name: instance_name,
						Subregion: &odk.BaseResource{
							Id: subregion_id,
						},
						Template: &(odk.BaseResource{template_id}),
						Type_:    &(odk.DictionaryItem{type_id, "", &(odk.Resource{})}),
					},
				},
			})
		})
}
