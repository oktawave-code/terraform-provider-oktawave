package oktawave

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oks "github.com/oktawave-code/oks-sdk"
)

func resourceOksCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOksClusterCreate,
		ReadContext:   resourceOksClusterRead,
		DeleteContext: resourceOksClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_running": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},
	}
}

func resourceOksClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "creating oks cluster")

	client := m.(*ClientConfig).oksClient
	auth := m.(*ClientConfig).oksAuth

	createCmd := oks.K44sClusterCreateDto{
		Version: d.Get("version").(string),
	}
	tflog.Debug(ctx, "calling OKS ClustersApi.ClustersNamePost")
	cluster, _, err := client.ClustersApi.ClustersNamePost(*auth, createCmd, d.Get("name").(string))
	if err != nil {
		return diag.Errorf("OKS Error in ClustersApi.ClustersNamePost. Caused by: %s", err)
	}

	err = waitUntilClusterIsOperational(ctx, client, auth, cluster.Name)
	if err != nil {
		return diag.Errorf("OKS Error while waiting for cluster. %s", err)
	}

	tflog.Info(ctx, fmt.Sprintf("successfully created OKS cluster. id=%v", cluster.Name))
	d.SetId(cluster.Name)

	return resourceOksClusterRead(ctx, d, m)
}

func resourceOksClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "reading oks cluster")

	client := m.(*ClientConfig).oksClient
	auth := m.(*ClientConfig).oksAuth

	tflog.Debug(ctx, "calling OKS ClustersApi.ClustersNameGet")
	cluster, resp, err := client.ClustersApi.ClustersNameGet(*auth, d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
		}
		return diag.Errorf("ODK Error in ClustersApi.ClustersNameGet. %s", err)
	}

	return loadOksClusterData(ctx, d, m, cluster)
}

func resourceOksClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "deleting oks cluster")

	client := m.(*ClientConfig).oksClient
	auth := m.(*ClientConfig).oksAuth

	tflog.Debug(ctx, "calling OKS ClustersApi.ClustersNameDelete")
	_, _, err := client.ClustersApi.ClustersNameDelete(*auth, d.Id())
	if err != nil {
		return diag.Errorf("ODK Error in ClustersApi.ClustersNameDelete. %s", err)
	}

	d.SetId("")
	return nil
}

func loadOksClusterData(ctx context.Context, d *schema.ResourceData, m interface{}, cluster oks.K44SClusterDetailsDto) diag.Diagnostics {
	// Store everything
	// tflog.Debug(ctx, "Parsing returned data")
	// if d.Set("name", cluster.Name) != nil {
	// 	return diag.Errorf("Can't retrieve name")
	// }
	if d.Set("version", cluster.Version) != nil {
		return diag.Errorf("Can't retrieve version")
	}
	if d.Set("creation_date", cluster.CreationDate.String()) != nil {
		return diag.Errorf("Can't retrieve creation date")
	}
	if d.Set("is_running", cluster.Running) != nil {
		return diag.Errorf("Can't retrieve running state")
	}
	return nil
}

func waitUntilClusterIsOperational(ctx context.Context, client oks.APIClient, auth *context.Context, name string) error {
	tflog.Info(ctx, "Waiting for cluster")
	maxRetries := 18
	found := false
	first := true
	for !found {
		if !first {
			tflog.Debug(ctx, "Still waiting")
			time.Sleep(10 * time.Second)
		}
		first = false
		tflog.Debug(ctx, "calling OKS ClustersApi.ClustersNameGet")
		cluster, _, err := client.ClustersApi.ClustersNameGet(*auth, name)
		if err != nil {
			tflog.Warn(ctx, fmt.Sprintf("OKS Error in ClustersApi.ClustersNameGet. %v", err))
			if maxRetries <= 0 {
				return err
			}
			maxRetries--
		}
		found = cluster.Running
	}
	return nil
}
