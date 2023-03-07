package oktawave

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

func getInstanceDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"name": {
			Type:        schema.TypeString,
			Description: "Name in https://api.oktawave.com/beta/docs/index#!/OCI/Instances_Get",
			Computed:    true,
		},
		"subregion_id": {
			Type:        schema.TypeInt,
			Description: "Subregion ID https://api.oktawave.com/beta/docs/index#!/Subregions/Subregions_Get",
			Computed:    true,
		},
		"template_id": {
			Type:        schema.TypeInt,
			Description: "Template ID https://api.oktawave.com/beta/docs/index#!/OCI_Templates/Templates_Get",
			Computed:    true,
		},
		"type_id": {
			Type:        schema.TypeInt,
			Description: "Instance type https://api.oktawave.com/beta/docs/index#!/OCI/Instances_GetInstancesTypes",
			Computed:    true,
		},
		"type_label": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"is_freemium": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"creation_date": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"creation_user_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		// tod
		// "creation_user_first_name": {
		// 	Type:     schema.TypeString,
		// 	Computed: true,
		// },
		// "creation_user_first_last_name": {
		// 	Type:     schema.TypeString,
		// 	Computed: true,
		// },
		"is_locked": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"locking_date": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ip_address": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"private_ip_address": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"dns_address": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"status_label": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"system_category_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"system_category_label": {
			Type:     schema.TypeString,
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
		"vmware_tools_status_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"vmware_tools_status_label": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"monit_status_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"monit_status_label": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"template_type_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"template_type_label": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"payment_type_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"payment_type_label": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"scsi_controller_type_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"scsi_controller_type_label": {
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
		"health_check_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"support_type_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"support_type_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		// todo
		// "is_hot_plug_enabled": {
		// 	Type:     schema.TypeBool,
		// 	Computed: true,
		// },
	}
}

func dataSourceInstance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceRead,
		Schema:      getInstanceDataSourceSchema(),
	}
}

func dataSourceInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "reading instance")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	id, ok := d.GetOk("id")
	if !ok {
		return diag.Errorf("Id must be specified.")
	}

	vars := map[string]interface{}{}
	instance, _, err := client.OCIApi.InstancesGet_2(*auth, int32(id.(int)), vars)
	if err != nil {
		return diag.Errorf("Instance with id %d not found. Caused by: %s", id.(int), err)
	}

	return loadDataSourceInstanceToSchema(d, instance)
}

func loadDataSourceInstanceToSchema(d *schema.ResourceData, instance odk.Instance) diag.Diagnostics {
	if err := d.Set("id", instance.Id); err != nil {
		return diag.Errorf("Error setting id: %s", err)
	}

	if err := d.Set("name", instance.Name); err != nil {
		return diag.Errorf("Error setting name: %s", err)
	}

	if err := d.Set("subregion_id", instance.Subregion.Id); err != nil {
		return diag.Errorf("Error setting subregion_id: %s", err)
	}

	if err := d.Set("template_id", instance.Template.Id); err != nil {
		return diag.Errorf("Error setting template_id: %s", err)
	}

	if err := d.Set("type_id", instance.Type_.Id); err != nil {
		return diag.Errorf("Error setting type_id: %s", err)
	}

	if err := d.Set("type_label", instance.Type_.Label); err != nil {
		return diag.Errorf("Error setting type_label: %s", err)
	}

	if err := d.Set("is_freemium", instance.IsFreemium); err != nil {
		return diag.Errorf("Error setting is_freemium: %s", err)
	}

	if err := d.Set("creation_date", instance.CreationDate.String()); err != nil {
		return diag.Errorf("Error setting creation_date: %s", err)
	}

	if err := d.Set("creation_user_id", instance.CreationUser.Id); err != nil {
		return diag.Errorf("Error setting creation_user_id: %s", err)
	}

	// if err := d.Set("creation_user_first_name", instance.CreationUser.FirstName); err != nil {
	// 	return diag.Errorf("Error setting creation_user_first_name: %s", err)
	// }

	// if err := d.Set("creation_user_first_last_name", instance.Crea); err != nil {
	// 	return diag.Errorf("Error setting creation_user_first_last_name: %s", err)
	// }

	if err := d.Set("is_locked", instance.IsLocked); err != nil {
		return diag.Errorf("Error setting is_locked: %s", err)
	}

	if err := d.Set("locking_date", instance.LockingDate.String()); err != nil {
		return diag.Errorf("Error setting locking_date: %s", err)
	}

	if err := d.Set("ip_address", instance.IpAddress); err != nil {
		return diag.Errorf("Error setting ip_address: %s", err)
	}

	if err := d.Set("private_ip_address", instance.PrivateIpAddress); err != nil {
		return diag.Errorf("Error setting private_ip_address: %s", err)
	}

	if err := d.Set("dns_address", instance.DnsAddress); err != nil {
		return diag.Errorf("Error setting dns_address: %s", err)
	}

	if err := d.Set("status_id", instance.Status.Id); err != nil {
		return diag.Errorf("Error setting status_id: %s", err)
	}

	if err := d.Set("status_label", instance.Status.Label); err != nil {
		return diag.Errorf("Error setting status_label: %s", err)
	}

	if err := d.Set("system_category_id", instance.SystemCategory.Id); err != nil {
		return diag.Errorf("Error setting system_category_id: %s", err)
	}

	if err := d.Set("system_category_label", instance.SystemCategory.Label); err != nil {
		return diag.Errorf("Error setting system_category_label: %s", err)
	}

	if err := d.Set("autoscaling_type_id", instance.AutoscalingType.Id); err != nil {
		return diag.Errorf("Error setting autoscaling_type_id: %s", err)
	}

	if err := d.Set("autoscaling_type_label", instance.AutoscalingType.Label); err != nil {
		return diag.Errorf("Error setting autoscaling_type_label: %s", err)
	}

	if err := d.Set("vmware_tools_status_id", instance.VmWareToolsStatus.Id); err != nil {
		return diag.Errorf("Error setting vmware_tools_status_id: %s", err)
	}

	if err := d.Set("vmware_tools_status_label", instance.VmWareToolsStatus.Label); err != nil {
		return diag.Errorf("Error setting vmware_tools_status_label: %s", err)
	}

	if err := d.Set("monit_status_id", instance.MonitStatus.Id); err != nil {
		return diag.Errorf("Error setting monit_status_id: %s", err)
	}

	if err := d.Set("monit_status_label", instance.MonitStatus.Label); err != nil {
		return diag.Errorf("Error setting monit_status_label: %s", err)
	}

	if err := d.Set("template_type_id", instance.TemplateType.Id); err != nil {
		return diag.Errorf("Error setting template_type_id: %s", err)
	}

	if err := d.Set("template_type_label", instance.TemplateType.Label); err != nil {
		return diag.Errorf("Error setting template_type_label: %s", err)
	}

	if err := d.Set("payment_type_id", instance.PaymentType.Id); err != nil {
		return diag.Errorf("Error setting payment_type_id: %s", err)
	}

	if err := d.Set("payment_type_label", instance.PaymentType.Label); err != nil {
		return diag.Errorf("Error setting payment_type_label: %s", err)
	}

	if err := d.Set("scsi_controller_type_id", instance.ScsiControllerType.Id); err != nil {
		return diag.Errorf("Error setting scsi_controller_type_id: %s", err)
	}

	if err := d.Set("scsi_controller_type_label", instance.ScsiControllerType.Label); err != nil {
		return diag.Errorf("Error setting scsi_controller_type_label: %s", err)
	}

	if err := d.Set("total_disks_capacity", instance.TotalDisksCapacity); err != nil {
		return diag.Errorf("Error setting total_disks_capacity: %s", err)
	}

	if err := d.Set("cpu_number", instance.CpuNumber); err != nil {
		return diag.Errorf("Error setting cpu_number: %s", err)
	}

	if err := d.Set("ram_mb", instance.RamMb); err != nil {
		return diag.Errorf("Error setting ram_mb: %s", err)
	}

	if instance.HealthCheck != nil {
		if err := d.Set("health_check_id", instance.HealthCheck.Id); err != nil {
			return diag.Errorf("Error setting health_check_id: %s", err)
		}
	}

	if instance.SupportType != nil {
		if err := d.Set("support_type_id", instance.SupportType.Id); err != nil {
			return diag.Errorf("Error setting support_type_id: %s", err)
		}

		if err := d.Set("support_type_name", instance.SupportType.Name); err != nil {
			return diag.Errorf("Error setting support_type_name: %s", err)
		}
	}

	d.SetId(strconv.Itoa(int(instance.Id)))
	return nil
}
