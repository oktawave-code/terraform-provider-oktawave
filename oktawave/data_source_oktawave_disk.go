package oktawave

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

func getDiskDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tier_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Value from dictionary #17",
		},
		"tier_label": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"subregion_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "ID from subregions resource",
		},
		"capacity": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "In GB",
		},
		"shared_disk_type_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Value from dictionary #162",
		},
		"shared_disk_type_label": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"connections": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"instance_id": {
						Type:     schema.TypeInt,
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
		"creation_user_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"creation_date": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"is_shared": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"is_locked": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"locking_date": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"is_freemium": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	}
}

func dataSourceDisk() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDiskRead,
		Schema:      getDiskDataSourceSchema(),
	}
}

func dataSourceDiskRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "reading disk")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	id, ok := d.GetOk("id")
	if !ok {
		return diag.Errorf("Id must be specified.")
	}

	vars := map[string]interface{}{}
	disk, _, err := client.OVSApi.DisksGet(*auth, int32(id.(int)), vars)
	if err != nil {
		return diag.Errorf("Disk with id %d not found. Caused by: %s", id.(int), err)
	}

	return loadDataSourceDiskToSchema(d, disk)
}

func loadDataSourceDiskToSchema(d *schema.ResourceData, disk odk.Disk) diag.Diagnostics {
	if err := d.Set("id", disk.Id); err != nil {
		return diag.Errorf("Error setting id: %s", err)
	}

	if err := d.Set("name", disk.Name); err != nil {
		return diag.Errorf("Error setting name: %s", err)
	}

	if err := d.Set("tier_id", disk.Tier.Id); err != nil {
		return diag.Errorf("Error setting tier_id: %s", err)
	}

	if err := d.Set("tier_label", disk.Tier.Label); err != nil {
		return diag.Errorf("Error setting tier_label: %s", err)
	}

	if err := d.Set("subregion_id", disk.Subregion.Id); err != nil {
		return diag.Errorf("Error setting subregion_id: %s", err)
	}

	if err := d.Set("capacity", disk.SpaceCapacity); err != nil {
		return diag.Errorf("Error setting capacity: %s", err)
	}

	if disk.SharedDiskType != nil {
		if err := d.Set("shared_disk_type_id", disk.SharedDiskType.Id); err != nil {
			return diag.Errorf("Error setting shared_disk_type_id: %s", err)
		}

		if err := d.Set("shared_disk_type_label", disk.SharedDiskType.Label); err != nil {
			return diag.Errorf("Error setting shared_disk_type_label: %s", err)
		}
	}

	if err := d.Set("creation_user_id", disk.CreationUser.Id); err != nil {
		return diag.Errorf("Error setting creation_user_id: %s", err)
	}

	if err := d.Set("creation_date", disk.CreationDate.String()); err != nil {
		return diag.Errorf("Error setting creation_date: %s", err)
	}

	if err := d.Set("is_shared", disk.IsShared); err != nil {
		return diag.Errorf("Error setting is_shared: %s", err)
	}

	if err := d.Set("is_locked", disk.IsLocked); err != nil {
		return diag.Errorf("Error setting is_locked: %s", err)
	}

	if err := d.Set("locking_date", disk.LockingDate.String()); err != nil {
		return diag.Errorf("Error setting locking_date: %s", err)
	}

	if err := d.Set("is_freemium", disk.IsFreemium); err != nil {
		return diag.Errorf("Error setting is_freemium: %s", err)
	}

	connections := make([]map[string]interface{}, len(disk.Connections))
	for idx, connection := range disk.Connections {
		connections[idx] = connectionToMap(connection)
	}

	if err := d.Set("connections", connections); err != nil {
		return diag.Errorf("Error setting connections: %s", err)
	}

	d.SetId(strconv.Itoa(int(disk.Id)))
	return nil
}

func connectionToMap(connection odk.DiskConnection) map[string]interface{} {
	return map[string]interface{}{
		"instance_id":    connection.Instance.Id,
		"controller":     connection.Controller,
		"slot":           connection.Slot,
		"is_system_disk": connection.IsSystemDisk,
	}
}
