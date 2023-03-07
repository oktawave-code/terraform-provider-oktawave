package oktawave

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/oktawave-code/oks-sdk"
)

func getOksClusterDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"version": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"creation_date": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"is_running": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	}
}

func dataSourceOksCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOksClusterRead,
		Schema:      getOksClusterDataSourceSchema(),
	}
}

func dataSourceOksClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "reading oks cluster")

	client := m.(*ClientConfig).oksClient
	auth := m.(*ClientConfig).oksAuth

	name, ok := d.GetOk("name")
	if !ok {
		return diag.Errorf("Name must be specified.")
	}

	cluster, _, err := client.ClustersApi.ClustersNameGet(*auth, name.(string))
	if err != nil {
		return diag.Errorf("Oks cluster with name %s not found. Caused by: %s", name.(string), err)
	}

	return loadDataSourceOksClusterToSchema(d, cluster)
}

func loadDataSourceOksClusterToSchema(d *schema.ResourceData, cluster swagger.K44SClusterDetailsDto) diag.Diagnostics {
	if err := d.Set("name", cluster.Name); err != nil {
		return diag.Errorf("Error setting name: %s", err)
	}

	if err := d.Set("version", cluster.Version); err != nil {
		return diag.Errorf("Error setting version: %s", err)
	}

	if err := d.Set("creation_date", cluster.CreationDate.String()); err != nil {
		return diag.Errorf("Error setting creation_date: %s", err)
	}

	if err := d.Set("is_running", cluster.Running); err != nil {
		return diag.Errorf("Error setting is_running: %s", err)
	}

	d.SetId(cluster.Name)
	return nil
}
