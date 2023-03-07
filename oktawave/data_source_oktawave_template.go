package oktawave

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

func getTemplateDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"version": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"creation_date": {
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
		"default_instance_type_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Value from dictionary #12",
		},
		"default_instance_type_label": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"minimum_instance_type_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Value from dictionary #12",
		},
		"minimum_instance_type_label": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ethernet_controllers_number": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"ethernet_controllers_type_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Value from dictionary #167",
		},
		"ethernet_controllers_type_label": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"system_category_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Value from dictionary #70",
		},
		"system_category_label": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"owner_account_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"publication_status_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Value from dictionary #140",
		},
		"publication_status_label": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"disks": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"space_capacity": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"tier_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"tier_label": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"creation_date": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"controller": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"slot": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"is_system_disk": {
						Type:     schema.TypeBool,
						Computed: true,
					},
				},
			},
		},
		"software": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
		"template_type_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Value from dictionary #52",
		},
		"template_type_label": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourceTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTemplateRead,
		Schema:      getTemplateDataSourceSchema(),
	}
}

func dataSourceTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "reading template")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	id, ok := d.GetOk("id")
	if !ok {
		return diag.Errorf("Id must be specified.")
	}

	vars := map[string]interface{}{}
	key, _, err := client.OCITemplatesApi.TemplatesGet_1(*auth, int32(id.(int)), vars)
	if err != nil {
		return diag.Errorf("Template with id %d not found. Caused by: %s", id.(int), err)
	}

	return loadDataSourceTemplateToSchema(d, key)
}

func loadDataSourceTemplateToSchema(d *schema.ResourceData, template odk.Template) diag.Diagnostics {
	disks := make([]map[string]interface{}, len(template.Disks))
	for idx, disk := range template.Disks {
		disks[idx] = makeTemplateDiskMap(disk)
	}

	softwares := make([]map[string]interface{}, len(template.Software))
	for idx, software := range template.Software {
		softwares[idx] = makeTemplateSoftwareMap(software)
	}

	if err := d.Set("id", template.Id); err != nil {
		return diag.Errorf("Error setting id: %s", err)
	}

	if err := d.Set("name", template.Name); err != nil {
		return diag.Errorf("Error setting name: %s", err)
	}

	if err := d.Set("description", template.Description); err != nil {
		return diag.Errorf("Error setting description: %s", err)
	}

	if err := d.Set("version", template.Version); err != nil {
		return diag.Errorf("Error setting version: %s", err)
	}

	if err := d.Set("creation_date", template.CreationDate.String()); err != nil {
		return diag.Errorf("Error setting creation_date: %s", err)
	}

	if err := d.Set("last_change_date", template.LastChangeDate.String()); err != nil {
		return diag.Errorf("Error setting last_change_date: %s", err)
	}

	if err := d.Set("creation_user_id", template.CreationUser.Id); err != nil {
		return diag.Errorf("Error setting creation_user_id: %s", err)
	}

	if err := d.Set("default_instance_type_id", template.DefaultInstanceType.Id); err != nil {
		return diag.Errorf("Error setting default_instance_type_id: %s", err)
	}

	if err := d.Set("default_instance_type_label", template.DefaultInstanceType.Label); err != nil {
		return diag.Errorf("Error setting default_instance_type_label: %s", err)
	}

	if err := d.Set("minimum_instance_type_id", template.MinimumInstanceType.Id); err != nil {
		return diag.Errorf("Error setting minimum_instance_type_id: %s", err)
	}

	if err := d.Set("minimum_instance_type_label", template.MinimumInstanceType.Label); err != nil {
		return diag.Errorf("Error setting minimum_instance_type_label: %s", err)
	}

	if err := d.Set("ethernet_controllers_number", template.EthernetControllersNumber); err != nil {
		return diag.Errorf("Error setting ethernet_controllers_number: %s", err)
	}

	if err := d.Set("ethernet_controllers_type_id", template.EthernetControllersType.Id); err != nil {
		return diag.Errorf("Error setting ethernet_controllers_type_id: %s", err)
	}

	if err := d.Set("ethernet_controllers_type_label", template.EthernetControllersType.Label); err != nil {
		return diag.Errorf("Error setting ethernet_controllers_type_label: %s", err)
	}

	if err := d.Set("system_category_id", template.SystemCategory.Id); err != nil {
		return diag.Errorf("Error setting system_category_id: %s", err)
	}

	if err := d.Set("system_category_label", template.SystemCategory.Label); err != nil {
		return diag.Errorf("Error setting system_category_label: %s", err)
	}

	if template.OwnerAccount != nil {
		if err := d.Set("owner_account_id", template.OwnerAccount.Id); err != nil {
			return diag.Errorf("Error setting owner_account_id: %s", err)
		}
	}

	if template.PublicationStatus != nil {
		if err := d.Set("publication_status_id", template.PublicationStatus.Id); err != nil {
			return diag.Errorf("Error setting publication_status_id: %s", err)
		}

		if err := d.Set("publication_status_label", template.PublicationStatus.Label); err != nil {
			return diag.Errorf("Error setting publication_status_label: %s", err)
		}
	}

	if err := d.Set("disks", disks); err != nil {
		return diag.Errorf("Error setting disks: %s", err)
	}

	if err := d.Set("software", softwares); err != nil {
		return diag.Errorf("Error setting software: %s", err)
	}

	if err := d.Set("template_type_id", template.TemplateType.Id); err != nil {
		return diag.Errorf("Error setting template_type_id: %s", err)
	}

	if err := d.Set("template_type_label", template.TemplateType.Label); err != nil {
		return diag.Errorf("Error setting template_type_label: %s", err)
	}

	d.SetId(strconv.Itoa(int(template.Id)))
	return nil
}

func makeTemplateDiskMap(disk odk.TemplateDisk) map[string]interface{} {
	return map[string]interface{}{
		"id":             disk.Id,
		"name":           disk.Name,
		"space_capacity": disk.SpaceCapacity,
		"tier_id":        disk.Tier.Id,
		"tier_label":     disk.Tier.Label,
		"creation_date":  disk.CreationDate.String(),
		"controller":     disk.Controller,
		"slot":           disk.Slot,
		"is_system_disk": disk.IsSystemDisk,
	}
}

func makeTemplateSoftwareMap(software odk.Software) map[string]interface{} {
	return map[string]interface{}{
		"id":   software.Id,
		"name": software.Name,
	}
}
