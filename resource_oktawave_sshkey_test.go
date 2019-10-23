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

func TestAccOktawave_SshKey_Basic(t *testing.T) {
	var sshKey odk.SshKey
	token := os.Getenv("TOKEN")
	mockStatus := os.Getenv("MOCK_STATUS")
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOktawaveSshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOktawaveSshKeyConfig_basic(token, mockStatus),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOktawaveSshKeyExists("oktawave_sshKey.my_key", &sshKey),
					testAccCheckOktawaveSshKeyAttributes_basic(&sshKey),
					resource.TestCheckResourceAttr("oktawave_sshKey.my_key", "ssh_key_name", "my_sshKey"),
					resource.TestCheckResourceAttr("oktawave_sshKey.my_key", "ssh_key_value", "ssh-rsa ABBAB3NzaC1yc2EAAAABIwAAAgEAwrr66r8n6B8Y0zMF3dOpXEapIQD9DiYQ6D6/zwor9o39jSkHNiMMER/GETBbzP83LOcekm02aRjo55ArO7gPPVvCXbrirJu9pkm4AC4BBre5xSLS7soyzwbigFruM8G63jSXqpHqJ/ooi168sKMC2b0Ncsi+JlTfNYlDXJVLKEeZgZOInQyMmtisaDTUQWTIv1snAizf4iIYENuAkGYGNCL77u5Y5VOu5eQipvFajTnps9QvUx/zdSFYn9e2sulWM3Bxc/S4IJ67JWHVRpfJxGi3hinRBH8WQdXuUwdJJTiJHKPyYrrM7Q6Xq4TOMFtcRuLDC6u3BXM1L0gBvHPNOnD5l2Lp5EjUkQ9CBf2j4A4gfH+iWQZyk08esAG/iwArAVxkl368+dkbMWOXL8BN4x5zYgdzoeypQZZ2RKH780MCTSo4WQ19DP8pw+9q3bSFC9H3xYAxrKAJNWjeTUJOTrTe+mWXXU770gYyQTxa2ycnYrlZucn1S3vsvn6eq7NZZ8NRbyv1n15Ocg+nHK4fuKOrwPhU3NbKQwtjb0Wsxx1gAmQqIOLTpAdsrAauPxC7TPYA5qQVCphvimKuhQM/1gMV225JrnjspVlthCzuFYUjXOKC3wxz6FFEtwnXu3uC5bVVkmkNadJmD21gD23yk4BraGXVYpRMIB+X+OTUUI8= dhopson@VMUbuntu-DSH"),
				),
				//ExpectError: regexp.MustCompile("errors during apply: Resource OCI. CREATE. Unable to create instance . Ticket status: 137"),
			},
		},
	})
}

func TestAccOktawave_SshKey_TierUpdate(t *testing.T) {
	var sshKey odk.SshKey
	token := os.Getenv("TOKEN")
	mockStatus := os.Getenv("MOCK_STATUS")
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOktawaveSshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOktawaveSshKeyConfig_basic(token, mockStatus),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOktawaveSshKeyExists("oktawave_sshKey.my_key", &sshKey),
					testAccCheckOktawaveSshKeyAttributes_basic(&sshKey),
					resource.TestCheckResourceAttr("oktawave_sshKey.my_key", "ssh_key_name", "my_sshKey"),
					resource.TestCheckResourceAttr("oktawave_sshKey.my_key", "ssh_key_value", "ssh-rsa ABBAB3NzaC1yc2EAAAABIwAAAgEAwrr66r8n6B8Y0zMF3dOpXEapIQD9DiYQ6D6/zwor9o39jSkHNiMMER/GETBbzP83LOcekm02aRjo55ArO7gPPVvCXbrirJu9pkm4AC4BBre5xSLS7soyzwbigFruM8G63jSXqpHqJ/ooi168sKMC2b0Ncsi+JlTfNYlDXJVLKEeZgZOInQyMmtisaDTUQWTIv1snAizf4iIYENuAkGYGNCL77u5Y5VOu5eQipvFajTnps9QvUx/zdSFYn9e2sulWM3Bxc/S4IJ67JWHVRpfJxGi3hinRBH8WQdXuUwdJJTiJHKPyYrrM7Q6Xq4TOMFtcRuLDC6u3BXM1L0gBvHPNOnD5l2Lp5EjUkQ9CBf2j4A4gfH+iWQZyk08esAG/iwArAVxkl368+dkbMWOXL8BN4x5zYgdzoeypQZZ2RKH780MCTSo4WQ19DP8pw+9q3bSFC9H3xYAxrKAJNWjeTUJOTrTe+mWXXU770gYyQTxa2ycnYrlZucn1S3vsvn6eq7NZZ8NRbyv1n15Ocg+nHK4fuKOrwPhU3NbKQwtjb0Wsxx1gAmQqIOLTpAdsrAauPxC7TPYA5qQVCphvimKuhQM/1gMV225JrnjspVlthCzuFYUjXOKC3wxz6FFEtwnXu3uC5bVVkmkNadJmD21gD23yk4BraGXVYpRMIB+X+OTUUI8= dhopson@VMUbuntu-DSH"),
				),
				//ExpectError: regexp.MustCompile("errors during apply: Resource OCI. CREATE. Unable to create instance . Ticket status: 137"),
			},
			{
				Config: testAccCheckOktawaveSshKeyConfig_Updated(token, mockStatus),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOktawaveSshKeyExists("oktawave_sshKey.my_key", &sshKey),
					testAccCheckOktawaveSshKeyAttributes_update(&sshKey),
					resource.TestCheckResourceAttr("oktawave_sshKey.my_key", "ssh_key_name", "my_sshKey1"),
					resource.TestCheckResourceAttr("oktawave_sshKey.my_key", "ssh_key_value", "ssh-rsa AABAB3NzaC1yc2EAAAABIwAAAgEAwrr66r8n6B8Y0zMF3dOpXEapIQD9DiYQ6D6/zwor9o39jSkHNiMMER/GETBbzP83LOcekm02aRjo55ArO7gPPVvCXbrirJu9pkm4AC4BBre5xSLS7soyzwbigFruM8G63jSXqpHqJ/ooi168sKMC2b0Ncsi+JlTfNYlDXJVLKEeZgZOInQyMmtisaDTUQWTIv1snAizf4iIYENuAkGYGNCL77u5Y5VOu5eQipvFajTnps9QvUx/zdSFYn9e2sulWM3Bxc/S4IJ67JWHVRpfJxGi3hinRBH8WQdXuUwdJJTiJHKPyYrrM7Q6Xq4TOMFtcRuLDC6u3BXM1L0gBvHPNOnD5l2Lp5EjUkQ9CBf2j4A4gfH+iWQZyk08esAG/iwArAVxkl368+dkbMWOXL8BN4x5zYgdzoeypQZZ2RKH780MCTSo4WQ19DP8pw+9q3bSFC9H3xYAxrKAJNWjeTUJOTrTe+mWXXU770gYyQTxa2ycnYrlZucn1S3vsvn6eq7NZZ8NRbyv1n15Ocg+nHK4fuKOrwPhU3NbKQwtjb0Wsxx1gAmQqIOLTpAdsrAauPxC7TPYA5qQVCphvimKuhQM/1gMV225JrnjspVlthCzuFYUjXOKC3wxz6FFEtwnXu3uC5bVVkmkNadJmD21gD23yk4BraGXVYpRMIB+X+OTUUI8= dhopson@VMUbuntu-DSH"),
				),
			},
		},
	})
}

func testAccCheckOktawaveSshKeyExists(name string, sshKey *odk.SshKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ssh key ID is set")
		}

		client := testAccProvider.Meta().(*ClientConfig).oktaClient()
		auth := testAccProvider.Meta().(*ClientConfig).ctx
		id, _ := strconv.Atoi(rs.Primary.ID)
		foundSshKey, _, err := client.AccountApi.AccountGetSshKey(*auth, int32(id), nil)
		if err != nil {
			return fmt.Errorf("Ssh key was not found by id: %s", rs.Primary.ID)
		}
		log.Printf("Ssh key was found by id %s and name %v", rs.Primary.ID, foundSshKey.Name)
		*sshKey = foundSshKey
		return nil
	}
}

func testAccCheckOktawaveSshKeyAttributes_basic(sshKey *odk.SshKey) resource.TestCheckFunc {
	return func(t *terraform.State) error {
		if sshKey.Name != "my_sshKey" {
			return fmt.Errorf("Bad Ssh key name. Expected: my_sshKey. Got: %v", sshKey.Name)
		}
		if sshKey.Value != "ssh-rsa ABBAB3NzaC1yc2EAAAABIwAAAgEAwrr66r8n6B8Y0zMF3dOpXEapIQD9DiYQ6D6/zwor9o39jSkHNiMMER/GETBbzP83LOcekm02aRjo55ArO7gPPVvCXbrirJu9pkm4AC4BBre5xSLS7soyzwbigFruM8G63jSXqpHqJ/ooi168sKMC2b0Ncsi+JlTfNYlDXJVLKEeZgZOInQyMmtisaDTUQWTIv1snAizf4iIYENuAkGYGNCL77u5Y5VOu5eQipvFajTnps9QvUx/zdSFYn9e2sulWM3Bxc/S4IJ67JWHVRpfJxGi3hinRBH8WQdXuUwdJJTiJHKPyYrrM7Q6Xq4TOMFtcRuLDC6u3BXM1L0gBvHPNOnD5l2Lp5EjUkQ9CBf2j4A4gfH+iWQZyk08esAG/iwArAVxkl368+dkbMWOXL8BN4x5zYgdzoeypQZZ2RKH780MCTSo4WQ19DP8pw+9q3bSFC9H3xYAxrKAJNWjeTUJOTrTe+mWXXU770gYyQTxa2ycnYrlZucn1S3vsvn6eq7NZZ8NRbyv1n15Ocg+nHK4fuKOrwPhU3NbKQwtjb0Wsxx1gAmQqIOLTpAdsrAauPxC7TPYA5qQVCphvimKuhQM/1gMV225JrnjspVlthCzuFYUjXOKC3wxz6FFEtwnXu3uC5bVVkmkNadJmD21gD23yk4BraGXVYpRMIB+X+OTUUI8= dhopson@VMUbuntu-DSH" {
			return fmt.Errorf("Bad Ssh key value. Expected: ssh-rsa ABBAB3NzaC1yc2EAAAABIwAAAg"+
				"EAwrr66r8n6B8Y0zMF3dOpXEapIQD9DiYQ6D6/zwor9o39jSkHNiMMER"+
				"/GETBbzP83LOcekm02aRjo55ArO7gPPVvCXbrirJu9pkm4AC4BBre5xSLS7soyzwbigFruM8G63jSXqpHqJ"+
				"/ooi168sKMC2b0Ncsi+JlTfNYlDXJVLKEeZgZOInQyMmtisaDTUQWTIv1snAizf4iIYENuAkGYGNCL77u5Y5VOu5eQipvFajTnps9QvUx/zdSFYn9e2sulWM3Bxc"+
				"/S4IJ67JWHVRpfJxGi3hinRBH8WQdXuUwdJJTiJHKPyYrrM7Q6Xq4TOMFtcRuLDC6u3BXM1L0gBvHPNOnD5l2Lp5EjUkQ9CBf2j4A4gfH+iWQZyk08esAG/"+
				"iwArAVxkl368+dkbMWOXL8BN4x5zYgdzoeypQZZ2RKH780MCTSo4WQ19DP8pw+9q3bSFC9H3xYAxrKAJNWjeTUJOTrTe+mWXXU770gYyQTxa2ycnYrlZucn1S3"+
				"vsvn6eq7NZZ8NRbyv1n15Ocg+nHK4fuKOrwPhU3NbKQwtjb0Wsxx1gAmQqIOLTpAdsrAauPxC7TPYA5qQVCphvimKuhQM/1gMV225JrnjspVlthCzuFYUjXOKC3wxz"+
				"6FFEtwnXu3uC5bVVkmkNadJmD21gD23yk4BraGXVYpRMIB+X+OTUUI8= dhopson@VMUbuntu-DSH. Got: %v", sshKey.Value)
		}
		return nil
	}
}

func testAccCheckOktawaveSshKeyAttributes_update(sshKey *odk.SshKey) resource.TestCheckFunc {
	return func(t *terraform.State) error {
		if sshKey.Name != "my_sshKey1" {
			return fmt.Errorf("Bad Ssh key name. Expected: my_sshKey1. Got: %v", sshKey.Name)
		}
		if sshKey.Value != "ssh-rsa AABAB3NzaC1yc2EAAAABIwAAAgEAwrr66r8n6B8Y0zMF3dOpXEapIQD9DiYQ6D6/zwor9o39jSkHNiMMER/GETBbzP83LOcekm02aRjo55ArO7gPPVvCXbrirJu9pkm4AC4BBre5xSLS7soyzwbigFruM8G63jSXqpHqJ/ooi168sKMC2b0Ncsi+JlTfNYlDXJVLKEeZgZOInQyMmtisaDTUQWTIv1snAizf4iIYENuAkGYGNCL77u5Y5VOu5eQipvFajTnps9QvUx/zdSFYn9e2sulWM3Bxc/S4IJ67JWHVRpfJxGi3hinRBH8WQdXuUwdJJTiJHKPyYrrM7Q6Xq4TOMFtcRuLDC6u3BXM1L0gBvHPNOnD5l2Lp5EjUkQ9CBf2j4A4gfH+iWQZyk08esAG/iwArAVxkl368+dkbMWOXL8BN4x5zYgdzoeypQZZ2RKH780MCTSo4WQ19DP8pw+9q3bSFC9H3xYAxrKAJNWjeTUJOTrTe+mWXXU770gYyQTxa2ycnYrlZucn1S3vsvn6eq7NZZ8NRbyv1n15Ocg+nHK4fuKOrwPhU3NbKQwtjb0Wsxx1gAmQqIOLTpAdsrAauPxC7TPYA5qQVCphvimKuhQM/1gMV225JrnjspVlthCzuFYUjXOKC3wxz6FFEtwnXu3uC5bVVkmkNadJmD21gD23yk4BraGXVYpRMIB+X+OTUUI8= dhopson@VMUbuntu-DSH" {
			return fmt.Errorf("Bad Ssh key value. Expected: ssh-rsa AABAB3NzaC1yc2EAAAABIwAAAg"+
				"EAwrr66r8n6B8Y0zMF3dOpXEapIQD9DiYQ6D6/zwor9o39jSkHNiMMER"+
				"/GETBbzP83LOcekm02aRjo55ArO7gPPVvCXbrirJu9pkm4AC4BBre5xSLS7soyzwbigFruM8G63jSXqpHqJ"+
				"/ooi168sKMC2b0Ncsi+JlTfNYlDXJVLKEeZgZOInQyMmtisaDTUQWTIv1snAizf4iIYENuAkGYGNCL77u5Y5VOu5eQipvFajTnps9QvUx/zdSFYn9e2sulWM3Bxc"+
				"/S4IJ67JWHVRpfJxGi3hinRBH8WQdXuUwdJJTiJHKPyYrrM7Q6Xq4TOMFtcRuLDC6u3BXM1L0gBvHPNOnD5l2Lp5EjUkQ9CBf2j4A4gfH+iWQZyk08esAG/"+
				"iwArAVxkl368+dkbMWOXL8BN4x5zYgdzoeypQZZ2RKH780MCTSo4WQ19DP8pw+9q3bSFC9H3xYAxrKAJNWjeTUJOTrTe+mWXXU770gYyQTxa2ycnYrlZucn1S3"+
				"vsvn6eq7NZZ8NRbyv1n15Ocg+nHK4fuKOrwPhU3NbKQwtjb0Wsxx1gAmQqIOLTpAdsrAauPxC7TPYA5qQVCphvimKuhQM/1gMV225JrnjspVlthCzuFYUjXOKC3wxz"+
				"6FFEtwnXu3uC5bVVkmkNadJmD21gD23yk4BraGXVYpRMIB+X+OTUUI8= dhopson@VMUbuntu-DSH. Got: %v", sshKey.Value)
		}
		return nil
	}
}

func testAccCheckOktawaveSshKeyDestroy(s *terraform.State) error {
	if os.Getenv("MOCK_STATUS") == "1" {
		httpmock.RegisterNoResponder(httpmock.NewStringResponder(404, "Ssh key was not found"))
	}
	client := testAccProvider.Meta().(*ClientConfig).oktaClient()
	auth := testAccProvider.Meta().(*ClientConfig).ctx
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oktawave_sshKey" {
			id, _ := strconv.Atoi(rs.Primary.ID)
			_, resp, err := client.AccountApi.AccountGetSshKey(*auth, int32(id), nil)
			if err != nil && resp.StatusCode == 404 {
				return fmt.Errorf("Error waitiing Ssh key to be destroyed. Error: %v", err.Error())
			}
			break
		}
	}
	return nil
}

func testAccCheckOktawaveSshKeyConfig_basic(token string, mockStatus string) string {
	sshKeyName := "my_sshKey"
	sshKeyValue := "ssh-rsa ABBAB3NzaC1yc2EAAAABIwAAAgEAwrr66r8n6B8Y0zMF3dOpXEapIQD9DiYQ6D6/zwor9o39jSkHNiMMER/GETBbzP83LOcekm02aRjo55ArO7gPPVvCXbrirJu9pkm4AC4BBre5xSLS7soyzwbigFruM8G63jSXqpHqJ/ooi168sKMC2b0Ncsi+JlTfNYlDXJVLKEeZgZOInQyMmtisaDTUQWTIv1snAizf4iIYENuAkGYGNCL77u5Y5VOu5eQipvFajTnps9QvUx/zdSFYn9e2sulWM3Bxc/S4IJ67JWHVRpfJxGi3hinRBH8WQdXuUwdJJTiJHKPyYrrM7Q6Xq4TOMFtcRuLDC6u3BXM1L0gBvHPNOnD5l2Lp5EjUkQ9CBf2j4A4gfH+iWQZyk08esAG/iwArAVxkl368+dkbMWOXL8BN4x5zYgdzoeypQZZ2RKH780MCTSo4WQ19DP8pw+9q3bSFC9H3xYAxrKAJNWjeTUJOTrTe+mWXXU770gYyQTxa2ycnYrlZucn1S3vsvn6eq7NZZ8NRbyv1n15Ocg+nHK4fuKOrwPhU3NbKQwtjb0Wsxx1gAmQqIOLTpAdsrAauPxC7TPYA5qQVCphvimKuhQM/1gMV225JrnjspVlthCzuFYUjXOKC3wxz6FFEtwnXu3uC5bVVkmkNadJmD21gD23yk4BraGXVYpRMIB+X+OTUUI8= dhopson@VMUbuntu-DSH"

	if mockStatus == "1" {
		httpmock.Activate()
		mockPostSshKey(sshKeyName, sshKeyValue)
		mockGetSshKey(sshKeyName, sshKeyValue)
		mockDeleteSshKey(sshKeyName, sshKeyValue)
	}

	return fmt.Sprintf(`
provider "oktawave" {
		access_token="%s"
		api_url = "https://api.oktawave.com/beta"
}

resource "oktawave_sshKey" "my_key"{
	ssh_key_name="%s"
	ssh_key_value="%s"
}`, token, sshKeyName, sshKeyValue)
}

func testAccCheckOktawaveSshKeyConfig_Updated(token string, mockStatus string) string {
	sshKeyName := "my_sshKey1"
	sshKeyValue := "ssh-rsa AABAB3NzaC1yc2EAAAABIwAAAgEAwrr66r8n6B8Y0zMF3dOpXEapIQD9DiYQ6D6/zwor9o39jSkHNiMMER/GETBbzP83LOcekm02aRjo55ArO7gPPVvCXbrirJu9pkm4AC4BBre5xSLS7soyzwbigFruM8G63jSXqpHqJ/ooi168sKMC2b0Ncsi+JlTfNYlDXJVLKEeZgZOInQyMmtisaDTUQWTIv1snAizf4iIYENuAkGYGNCL77u5Y5VOu5eQipvFajTnps9QvUx/zdSFYn9e2sulWM3Bxc/S4IJ67JWHVRpfJxGi3hinRBH8WQdXuUwdJJTiJHKPyYrrM7Q6Xq4TOMFtcRuLDC6u3BXM1L0gBvHPNOnD5l2Lp5EjUkQ9CBf2j4A4gfH+iWQZyk08esAG/iwArAVxkl368+dkbMWOXL8BN4x5zYgdzoeypQZZ2RKH780MCTSo4WQ19DP8pw+9q3bSFC9H3xYAxrKAJNWjeTUJOTrTe+mWXXU770gYyQTxa2ycnYrlZucn1S3vsvn6eq7NZZ8NRbyv1n15Ocg+nHK4fuKOrwPhU3NbKQwtjb0Wsxx1gAmQqIOLTpAdsrAauPxC7TPYA5qQVCphvimKuhQM/1gMV225JrnjspVlthCzuFYUjXOKC3wxz6FFEtwnXu3uC5bVVkmkNadJmD21gD23yk4BraGXVYpRMIB+X+OTUUI8= dhopson@VMUbuntu-DSH"
	if mockStatus == "1" {
		mockDeleteSshKey(sshKeyName, sshKeyValue)
	}
	return fmt.Sprintf(`
provider "oktawave" {
		access_token="%s"
		api_url = "https://api.oktawave.com/beta"
}

resource "oktawave_sshKey" "my_key"{
	ssh_key_name="%s"
	ssh_key_value="%s"
}`, token, sshKeyName, sshKeyValue)
}

func mockPostSshKey(sshKeyName string, sshKeyValue string) {
	httpmock.RegisterResponder(http.MethodPost, "https://api.oktawave.com/beta/account/sshkeys",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, odk.SshKey{
				Id:           1658,
				Name:         sshKeyName,
				Value:        sshKeyValue,
				OwnerUser:    &odk.UserResource{},
				CreationDate: time.Now(),
			})
		})
}

func mockGetSshKey(sshKeyName string, sshKeyValue string) {
	httpmock.RegisterResponder(http.MethodGet, "https://api.oktawave.com/beta/account/sshkeys/1658",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewJsonResponse(200, odk.SshKey{
				Id:           1658,
				Name:         sshKeyName,
				Value:        sshKeyValue,
				OwnerUser:    &odk.UserResource{},
				CreationDate: time.Now(),
			})
		})
}

func mockDeleteSshKey(sshKeyName string, sshKeyValue string) {
	httpmock.RegisterResponder(http.MethodDelete, "https://api.oktawave.com/beta/account/sshkeys/1658",
		func(req *http.Request) (*http.Response, error) {
			mockGetSshKey(sshKeyName, sshKeyValue)
			return httpmock.NewJsonResponse(204, "")
		})
}
