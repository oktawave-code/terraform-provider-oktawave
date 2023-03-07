package oktawave

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

func getGroupDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"affinity_rule_type_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"affinity_rule_type_label": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"is_load_balancer": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"instances_count": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"schedulers_count": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"autoscaling_type_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"autoscaling_type_label": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"last_change_date": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"creation_user_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupRead,
		Schema:      getGroupDataSourceSchema(),
	}
}

func dataSourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "reading group")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	id, ok := d.GetOk("id")
	if !ok {
		return diag.Errorf("Id must be specified.")
	}

	vars := map[string]interface{}{}
	group, _, err := client.OCIGroupsApi.GroupsGetGroup(*auth, int32(id.(int)), vars)
	if err != nil {
		return diag.Errorf("Group with id %d not found. Caused by: %s", id.(int), err)
	}

	return loadDataSourceGroupToSchema(d, group)
}

func loadDataSourceGroupToSchema(d *schema.ResourceData, group odk.Group) diag.Diagnostics {
	if err := d.Set("id", group.Id); err != nil {
		return diag.Errorf("Error setting id: %s", err)
	}
	if err := d.Set("name", group.Name); err != nil {
		return diag.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("affinity_rule_type_id", group.AffinityRuleType.Id); err != nil {
		return diag.Errorf("Error setting affinity rule type id: %s", err)
	}
	if err := d.Set("affinity_rule_type_label", group.AffinityRuleType.Label); err != nil {
		return diag.Errorf("Error setting affinity rule type label: %s", err)
	}
	if err := d.Set("is_load_balancer", group.IsLoadBalancer); err != nil {
		return diag.Errorf("Error setting is load balancer: %s", err)
	}
	if err := d.Set("instances_count", group.InstancesCount); err != nil {
		return diag.Errorf("Error setting instances count: %s", err)
	}
	if err := d.Set("schedulers_count", group.SchedulersCount); err != nil {
		return diag.Errorf("Error setting schedulers count: %s", err)
	}
	if err := d.Set("autoscaling_type_id", group.AutoscalingType.Id); err != nil {
		return diag.Errorf("Error setting autoscaling type id: %s", err)
	}
	if err := d.Set("autoscaling_type_label", group.AutoscalingType.Label); err != nil {
		return diag.Errorf("Error setting autoscaling type label: %s", err)
	}
	if err := d.Set("last_change_date", group.LastChangeDate.String()); err != nil {
		return diag.Errorf("Error setting last change date: %s", err)
	}
	if err := d.Set("creation_user_id", group.CreationUser.Id); err != nil {
		return diag.Errorf("Error setting creation user name: %s", err)
	}

	d.SetId(strconv.Itoa(int(group.Id)))
	return nil
}
