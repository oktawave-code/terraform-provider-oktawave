package main

import (
	"context"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/oktawave-code/odk"
	"log"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_token": {
				Required: true,
				Type:     schema.TypeString,
			},
			"api_url": {
				Optional: true,
				Type:     schema.TypeString,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"oktawave_oci":                resourceOci(),
			"oktawave_ovs":                resourceOvs(),
			"oktawave_sshKey":             resourceSshKey(),
			"oktawave_group":              resourceGroup(),
			"oktawave_ip_address":         resourceIpAddress(),
			"oktawave_opn":                resourceOpn(),
			"oktawave_load_balancer":      resourceLoadBalancer(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	// TODO access token not set throw error
	log.Printf("[INFO] Initializing Oktawave provider.")
	auth := context.WithValue(context.Background(), odk.ContextAccessToken, d.Get("access_token"))
	config := odk.NewConfiguration()

	odkApiUrl, odkApiSet := d.GetOk("api_url")

	if odkApiSet && odkApiUrl != "" {
		config.BasePath = odkApiUrl.(string)
		log.Printf("[DEBUG] ODK API url was changed to: %s", odkApiUrl.(string))
	}

	odkClient := odk.NewAPIClient(config)

	client := ClientConfig{
		ctx:       &auth,
		odkClient: *odkClient,
	}
	return &client, nil
}
