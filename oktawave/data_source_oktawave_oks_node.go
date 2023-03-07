package oktawave

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/oktawave-code/oks-sdk"
)

func getOksNodeDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cluster_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"subregion_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "ID from subregions resource",
		},
		"type_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Value from dictionary #12",
		},
		"type_label": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Value from dictionary #12",
		},
		"status_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Value from dictionary #27",
		},
		"status_lavel": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Value from dictionary #27",
		},
		"creation_date": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ip_address": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"total_disks_capacity": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"cpu_number": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"ram_mb": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func dataSourceOksNode() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOksNodeRead,
		Schema:      getOksNodeDataSourceSchema(),
	}
}

func dataSourceOksNodeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "reading oks node")

	client := m.(*ClientConfig).oksClient
	auth := m.(*ClientConfig).oksAuth

	cluster_id, ok := d.GetOk("cluster_id")
	if !ok {
		return diag.Errorf("Cluster id must be specified.")
	}
	id, ok := d.GetOk("id")
	if !ok {
		return diag.Errorf("Node id must be specified.")
	}

	instances, _, err := client.ClustersApi.ClustersInstancesNameGet(*auth, cluster_id.(string))
	if err != nil {
		return diag.Errorf("Oks cluster with id %s not found. Caused by: %s", cluster_id.(string), err)
	}

	var instance *swagger.K44sInstance = nil
	for _, i := range instances {
		if int(i.Id) == id.(int) {
			instance = &i
			break
		}
	}

	if instance == nil {
		return diag.Errorf("Oks node with id %d for cluster %s not found. Caused by: %s", id, cluster_id.(string), err)
	}

	if err := d.Set("cluster_id", cluster_id); err != nil {
		return diag.Errorf("Error setting cluster_id: %s", err)
	}

	return loadDataSourceOksNodeToSchema(d, *instance)
}

func loadDataSourceOksNodeToSchema(d *schema.ResourceData, node swagger.K44sInstance) diag.Diagnostics {
	if err := d.Set("id", node.Id); err != nil {
		return diag.Errorf("Error setting id: %s", err)
	}

	if err := d.Set("name", node.Name); err != nil {
		return diag.Errorf("Error setting name: %s", err)
	}

	if err := d.Set("subregion_id", node.Subregion.Id); err != nil {
		return diag.Errorf("Error setting subregion_id: %s", err)
	}

	if err := d.Set("type_id", node.Type_.Id); err != nil {
		return diag.Errorf("Error setting type_id: %s", err)
	}

	if err := d.Set("type_label", node.Type_.Label); err != nil {
		return diag.Errorf("Error setting type_label: %s", err)
	}

	if err := d.Set("status_id", node.Status.Id); err != nil {
		return diag.Errorf("Error setting status_id: %s", err)
	}

	if err := d.Set("status_lavel", node.Status.Label); err != nil {
		return diag.Errorf("Error setting status_lavel: %s", err)
	}

	if err := d.Set("creation_date", node.CreationDate.String()); err != nil {
		return diag.Errorf("Error setting creation_date: %s", err)
	}

	if err := d.Set("ip_address", node.IpAddress); err != nil {
		return diag.Errorf("Error setting ip_address: %s", err)
	}

	if err := d.Set("total_disks_capacity", node.TotalDisksCapacity); err != nil {
		return diag.Errorf("Error setting total_disks_capacity: %s", err)
	}

	if err := d.Set("cpu_number", node.CpuNumber); err != nil {
		return diag.Errorf("Error setting cpu_number: %s", err)
	}

	if err := d.Set("ram_mb", node.RamMb); err != nil {
		return diag.Errorf("Error setting ram_mb: %s", err)
	}

	d.SetId(strconv.Itoa(int(node.Id)))
	return nil
}
