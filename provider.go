package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
	"github.com/oktawave-code/oks-sdk"
	"log"
)

func Provider() *schema.Provider {
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
			"oktawave_kubernetes_cluster": resourceKubernetesCluster(),
			"oktawave_kubernetes_node":    resourceNode(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"oktawave_oci": dataSourceOci(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	// TODO access token not set throw error
	log.Printf("[INFO] Initializing Oktawave provider.")
	auth := context.WithValue(context.Background(), odk.ContextAccessToken, d.Get("access_token"))
	authOKS := context.WithValue(context.Background(), swagger.ContextAccessToken, d.Get("access_token"))
	config := odk.NewConfiguration()
	oksCfg := swagger.NewConfiguration()

	odkApiUrl, odkApiSet := d.GetOk("api_url")

	if odkApiSet && odkApiUrl != "" {
		config.BasePath = odkApiUrl.(string)
		log.Printf("[DEBUG] ODK API url was changed to: %s", odkApiUrl.(string))
	}
	oksCfg.BasePath = "https://k44s-api.i.k44s.oktawave.com"

	oksClient := swagger.NewAPIClient(oksCfg)
	odkClient := odk.NewAPIClient(config)

	client := ClientConfig{
		oksCtx:    &authOKS,
		ctx:       &auth,
		odkClient: *odkClient,
		oksClient: *oksClient,
	}
	return &client, nil
}
