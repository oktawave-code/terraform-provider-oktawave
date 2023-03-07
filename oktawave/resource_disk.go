package oktawave

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

func resourceDisk() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDiskCreate,
		ReadContext:   resourceDiskRead,
		UpdateContext: resourceDiskUpdate,
		DeleteContext: resourceDiskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Required
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Disk name",
			},
			"tier_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Defines disk performance class. Value from dictionary #17",
			},
			"subregion_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID from subregions resource. Migration between subregions is possible if disk is not attached to instance.",
			},
			// Optional
			"capacity": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "Disk size in GB. At least 5 GB. Disk capacity can be only scaled up.",
			},
			"shared_disk_type_id": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "Type of disk sharing. Value from dictionary #162",
			},
			"instance_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "List of instances ids that are connected to this disk.",
			},
			// Computed
			"creation_user_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Id of user who created disk.",
			},
			"creation_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date of disk creation.",
			},
			"is_shared": {
				Type:       schema.TypeBool,
				Computed:   true,
				Deprecated: "Allows disk to be shared amongst multiple instances.",
			},
			"is_locked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Tells if disk is currently locked.",
			},
			"locking_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Disk lock date.",
			},
			"is_freemium": {
				Type:       schema.TypeBool,
				Computed:   true,
				Deprecated: "Tells if this disk is a freemium disk.",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},
		Description: "Oktawave Volume Storage(OVS) service provides block storage disks.",
	}
}

func resourceDiskCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "creating disk")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	createCommand := odk.CreateDiskCommand{
		DiskName:      d.Get("name").(string),
		SpaceCapacity: int32(d.Get("capacity").(int)),
		TierId:        int32(d.Get("tier_id").(int)),
		SubregionId:   int32(d.Get("subregion_id").(int)),
	}

	tflog.Debug(ctx, "calling OVSApi.DisksPost")
	ticket, _, err := client.OVSApi.DisksPost(*auth, createCommand)
	if err != nil {
		return diag.Errorf("ODK Error in OVSApi.DisksPost. %s", err)
	}

	createTicket, err := waitForTicket(client, auth, ticket)
	if err != nil {
		return diag.Errorf("ODK Error in TicketsApi.TicketsGet. %s", err)
	}
	if createTicket.Status.Id != DICT_TICKET_SUCCEED {
		return diag.Errorf("Unable to create OVS. Ticket status=%v", createTicket.Status.Id)
	}

	tflog.Info(ctx, fmt.Sprintf("successfully created OVS. id=%v", createTicket.ObjectId))
	d.SetId(strconv.Itoa(int(createTicket.ObjectId)))

	return resourceDiskRead(ctx, d, m)
}

func resourceDiskRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "reading disk")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	diskId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Invalid OVS id: %v %s", d.Id(), err)
	}

	tflog.Debug(ctx, "calling ODK OVSApi.DisksGet")
	disk, resp, err := client.OVSApi.DisksGet(*auth, int32(diskId), nil)
	if err != nil {
		if resp != nil && (resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusForbidden) { // api returns 403 on missing disk
			d.SetId("")
			return diag.Errorf("OVS %v not found", diskId)
		}
		return diag.Errorf("ODK Error in OVSApi.DisksGet. %s", err)
	}

	return loadDiskData(ctx, d, m, disk)
}

func resourceDiskUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "updating disk")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	diskId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Invalid OVS id: %v %s", d.Id(), err)
	}

	d.Partial(true)

	updateDiskCmd := odk.UpdateDiskCommand{
		DiskName:      d.Get("name").(string),
		SpaceCapacity: int32(d.Get("capacity").(int)),
		TierId:        int32(d.Get("tier_id").(int)),
	}

	if d.HasChange("subregion_id") {
		tflog.Info(ctx, "subregion id change detected")
		oldSubregion, newSubregion := d.GetChange("subregion_id")
		// detach from instances in old subregion
		updateDiskCmd.SubregionId = int32(oldSubregion.(int))
		if err := updateDisk(ctx, client, auth, int32(diskId), updateDiskCmd); err != nil {
			return err
		}
		// change subregion
		updateDiskCmd.SubregionId = int32(newSubregion.(int))
		if err := updateDisk(ctx, client, auth, int32(diskId), updateDiskCmd); err != nil {
			return err
		}
	}

	updateDiskCmd.SubregionId = int32(d.Get("subregion_id").(int))
	if err := updateDisk(ctx, client, auth, int32(diskId), updateDiskCmd); err != nil {
		return err
	}

	d.Partial(false)
	return resourceDiskRead(ctx, d, m)
}

func resourceDiskDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "deleting disk")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	diskId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Invalid OVS id: %v %s", d.Id(), err)
	}

	detachCommand := odk.UpdateDiskCommand{
		DiskName:        d.Get("name").(string),
		SpaceCapacity:   int32(d.Get("capacity").(int)),
		TierId:          int32(d.Get("tier_id").(int)),
		SubregionId:     int32(d.Get("subregion_id").(int)),
		InstanceIdsList: nil,
	}
	tflog.Debug(ctx, "calling ODK OVSApi.DisksPut")
	ticket, _, err := client.OVSApi.DisksPut(*auth, int32(diskId), detachCommand)
	if err != nil {
		return diag.Errorf("ODK Error in OVSApi.DisksPut. %s", err)
	}

	ticket, err = waitForTicket(client, auth, ticket)
	if err != nil {
		return diag.Errorf("ODK Error in TicketsApi.TicketsGet. %s", err)
	}
	if ticket.Status.Id != DICT_TICKET_SUCCEED {
		return diag.Errorf("Unable to detach OVS. Ticket status=%v", ticket.Status.Id)
	}

	tflog.Debug(ctx, "calling ODK OVSApi.DisksDelete")
	ticket, _, err = client.OVSApi.DisksDelete(*auth, int32(diskId))
	if err != nil {
		return diag.Errorf("ODK Error in OVSApi.DisksDelete. %s", err)
	}

	ticket, err = waitForTicket(client, auth, ticket)
	if err != nil {
		return diag.Errorf("ODK Error in TicketsApi.TicketsGet. %s", err)
	}
	if ticket.Status.Id != DICT_TICKET_SUCCEED {
		return diag.Errorf("Unable to delete OVS. Ticket status=%v", ticket.Status.Id)
	}

	d.SetId("")
	return nil
}

func loadDiskData(ctx context.Context, d *schema.ResourceData, m interface{}, disk odk.Disk) diag.Diagnostics {
	// Store everything
	tflog.Debug(ctx, "Parsing returned data")
	if d.Set("name", disk.Name) != nil {
		return diag.Errorf("Can't retrieve disk name")
	}
	if d.Set("tier_id", disk.Tier.Id) != nil {
		return diag.Errorf("Can't retrieve tier id")
	}
	if d.Set("subregion_id", disk.Subregion.Id) != nil {
		return diag.Errorf("Can't retrieve subregion id")
	}
	if d.Set("capacity", disk.SpaceCapacity) != nil {
		return diag.Errorf("Can't retrieve capacity")
	}
	if d.Set("creation_user_id", disk.CreationUser.Id) != nil {
		return diag.Errorf("Can't retrieve creation user id")
	}
	if d.Set("creation_date", disk.CreationDate.String()) != nil {
		return diag.Errorf("Can't retrieve creation date")
	}
	if d.Set("is_shared", disk.IsShared) != nil {
		return diag.Errorf("Can't retrieve shared state")
	}
	if disk.SharedDiskType != nil {
		if d.Set("shared_disk_type_id", disk.SharedDiskType.Id) != nil {
			return diag.Errorf("Can't retrieve shared type id")
		}
	}
	if d.Set("is_locked", disk.IsLocked) != nil {
		return diag.Errorf("Can't retrieve locked state")
	}
	if d.Set("locking_date", disk.LockingDate.String()) != nil {
		return diag.Errorf("Can't retrieve locking date")
	}
	if len(disk.Connections) > 0 {
		if d.Set("instance_ids", getConnectionInstanceIds(disk.Connections)) != nil {
			return diag.Errorf("Can't retrieve instances list")
		}
	}
	if d.Set("is_freemium", disk.IsFreemium) != nil {
		return diag.Errorf("Can't retrieve freemium state")
	}
	return nil
}

func updateDisk(ctx context.Context, client odk.APIClient, auth *context.Context, diskId int32, updateCmd odk.UpdateDiskCommand) diag.Diagnostics {
	tflog.Debug(ctx, "calling ODK OVSApi.DisksPut")
	ticket, _, err := client.OVSApi.DisksPut(*auth, diskId, updateCmd)
	if err != nil {
		return diag.Errorf("ODK Error in OVSApi.DisksPut. %s", err)
	}

	ticket, err = waitForTicket(client, auth, ticket)
	if err != nil {
		return diag.Errorf("ODK Error in TicketsApi.TicketsGet. %s", err)
	}
	if ticket.Status.Id != DICT_TICKET_SUCCEED {
		return diag.Errorf("Unable to modify OVS. Ticket status=%v", ticket.Status.Id)
	}

	return nil
}
