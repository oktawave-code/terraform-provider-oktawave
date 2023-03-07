package oktawave

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
	oks "github.com/oktawave-code/oks-sdk"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OKTAWAVE_ACCESS_TOKEN", nil),
			},
			"dc": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OKTAWAVE_DC", nil),
			},
			"odk_api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OKTAWAVE_ODK_API_URL", nil),
			},
			"odk_api_skip_tls": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: envBoolDefaultFunc("OKTAWAVE_ODK_API_SKIP_TLS", false),
			},
			"oks_api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OKTAWAVE_OKS_API_URL", nil),
			},
			"oks_api_skip_tls": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: envBoolDefaultFunc("OKTAWAVE_OKS_API_SKIP_TLS", false),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"oktawave_instance":      resourceInstance(),
			"oktawave_template":      resourceTemplate(),
			"oktawave_disk":          resourceDisk(),
			"oktawave_opn":           resourceOpn(),
			"oktawave_ip":            resourceIpAddress(),
			"oktawave_group":         resourceGroup(),
			"oktawave_load_balancer": resourceLoadBalancer(),
			"oktawave_ssh_key":       resourceSshKey(),
			"oktawave_oks_cluster":   resourceOksCluster(),
			"oktawave_oks_node":      resourceOksNode(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"oktawave_instance":       dataSourceInstance(),
			"oktawave_instances":      dataSourceInstances(),
			"oktawave_template":       dataSourceTemplate(),
			"oktawave_templates":      dataSourceTemplates(),
			"oktawave_disk":           dataSourceDisk(),
			"oktawave_disks":          dataSourceDisks(),
			"oktawave_opn":            dataSourceOpn(),
			"oktawave_opns":           dataSourceOpns(),
			"oktawave_ip":             dataSourceIp(),
			"oktawave_ips":            dataSourceIps(),
			"oktawave_group":          dataSourceGroup(),
			"oktawave_groups":         dataSourceGroups(),
			"oktawave_load_balancer":  dataSourceLoadBalancer(),
			"oktawave_load_balancers": dataSourceLoadBalancers(),
			"oktawave_ssh_key":        dataSourceSshKey(),
			"oktawave_ssh_keys":       dataSourceSshKeys(),
			"oktawave_oks_cluster":    dataSourceOksCluster(),
			"oktawave_oks_clusters":   dataSourceOksClusters(),
			"oktawave_oks_node":       dataSourceOksNode(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	tflog.Info(ctx, "Initializing Oktawave provider")

	access_token := d.Get("access_token").(string)
	odkAuth := context.WithValue(context.Background(), odk.ContextAccessToken, access_token)
	oksAuth := context.WithValue(context.Background(), oks.ContextAccessToken, access_token)

	var odkUrl string = ""
	var oksUrl string = ""
	dcId, dcIdSet := d.GetOk("dc")
	if dcIdSet {
		cfg, ok := dcConfigs[dcId.(string)]
		if ok {
			odkUrl = cfg.odkApiUrl
			oksUrl = cfg.oksApiUrl
		} else {
			tflog.Error(ctx, "Unknown DC id")
		}
	}
	odkApiUrl, odkApiSet := d.GetOk("odk_api_url")
	if odkApiSet && (odkApiUrl != "") {
		odkUrl = odkApiUrl.(string)
	}
	oksApiUrl, oksApiSet := d.GetOk("oks_api_url")
	if oksApiSet && (oksApiUrl != "") {
		oksUrl = oksApiUrl.(string)
	}

	odkCfg := odk.NewConfiguration()
	if odkUrl != "" {
		odkCfg.BasePath = odkUrl
		tflog.Info(ctx, fmt.Sprintf("ODK API url was set to \"%v\"", odkCfg.BasePath))
	}
	if d.Get("odk_api_skip_tls").(bool) {
		tflog.Info(ctx, "Disabling ODK API certificate verification test")
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		odkCfg.HTTPClient = &http.Client{Transport: tr}
	}

	oksCfg := oks.NewConfiguration()
	oksCfg.BasePath = "https://k44s-api.i.k44s.oktawave.com" // Default is not provided by library
	if oksUrl != "" {
		oksCfg.BasePath = oksUrl
		tflog.Info(ctx, fmt.Sprintf("OKS API url was set to \"%v\"", oksCfg.BasePath))
	}
	if d.Get("oks_api_skip_tls").(bool) {
		tflog.Info(ctx, "Disabling OKS API certificate verification test")
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		oksCfg.HTTPClient = &http.Client{Transport: tr}
	}

	tflog.Trace(ctx, "Connection settings", map[string]interface{}{
		"ODK_url": odkCfg.BasePath,
		"OKS_url": oksCfg.BasePath,
	})

	odkClient := odk.NewAPIClient(odkCfg)
	oksClient := oks.NewAPIClient(oksCfg)

	client := ClientConfig{
		odkAuth:   &odkAuth,
		odkClient: *odkClient,
		oksAuth:   &oksAuth,
		oksClient: *oksClient,
	}
	tflog.Debug(ctx, "Oktawave provider initialized")
	return &client, *new(diag.Diagnostics)
}

func envBoolDefaultFunc(k string, dv interface{}) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		if v := os.Getenv(k); v != "" {
			v = strings.ToLower(v)
			return (v == "yes") || (v == "y") || (v == "true") || (v == "t"), nil
		}
		return dv, nil
	}
}
