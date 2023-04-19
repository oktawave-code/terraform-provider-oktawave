package oktawave

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

func resourceOpn() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpnCreate,
		ReadContext:   resourceOpnRead,
		UpdateContext: resourceOpnUpdate,
		DeleteContext: resourceOpnDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of OPN",
			},
			"creation_user_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Id of user who created this resource.",
			},
			"creation_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation date",
			},
			"last_change_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last change date",
			},
			"instance_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "List of instance ids in this OPN.",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},
		Description: "Oktawave Private Network(OPN) is a conterpart of typical VLAN network.",
	}
}

func resourceOpnCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "creating opn")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	opnName := d.Get("name").(string)
	createCommand := odk.CreateOpnCommand{
		OpnName: opnName,
	}

	tflog.Debug(ctx, "calling ODK NetworkingApi.OpnsPost")
	ticket, _, err := client.NetworkingApi.OpnsPost(*auth, createCommand)
	if err != nil {
		return diag.Errorf("ODK Error in NetworkingApi.OpnsPost. %s", err)
	}

	createTicket, err := waitForTicket(client, auth, ticket)
	if err != nil {
		return diag.Errorf("ODK Error in TicketsApi.TicketsGet. %s", err)
	}
	if createTicket.Status.Id != DICT_TICKET_SUCCEED {
		return diag.Errorf("Unable to create OPN. Ticket status=%v", createTicket.Status.Id)
	}

	tflog.Info(ctx, fmt.Sprintf("successfully created OPN. id=%v", createTicket.ObjectId))
	d.SetId(strconv.Itoa(int(createTicket.ObjectId)))

	return resourceOpnRead(ctx, d, m)
}

func resourceOpnRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "reading opn")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	opnId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Invalid OPN id: %v %s", d.Id(), err)
	}

	tflog.Debug(ctx, "calling ODK NetworkingApi.OpnsGet_1")
	params := map[string]interface{}{
		"pageSize": int32(math.MaxInt16),
	}
	opn, resp, err := client.NetworkingApi.OpnsGet_1(*auth, int32(opnId), params)
	if err != nil {
		if resp != nil && (resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusForbidden) { // api returns 403 on missing opn
			d.SetId("")
		}
		return diag.Errorf("Error while retrieving OPN %v (possibly \"not found\" error): %s", opnId, err)
	}

	return loadOpnData(ctx, d, m, opn)
}

func resourceOpnUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "updating opn")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	opnId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Invalid OPN id: %v %s", d.Id(), err)
	}

	if d.HasChange("name") {
		tflog.Info(ctx, "opn name change detected")
		opnName := d.Get("name").(string)
		updateCommand := odk.UpdateOpnCommand{
			OpnName: opnName,
		}
		tflog.Debug(ctx, "calling ODK NetworkingApi.OpnsPut")
		_, _, err := client.NetworkingApi.OpnsPut(*auth, int32(opnId), updateCommand)
		if err != nil && err.Error() != "EOF" { // "EOF" condition is a patch for ODK 1.4 bug: it reports error when API returns empty body
			return diag.Errorf("Error while updating OPN %v: %s", int32(opnId), err)
		}
	}

	return resourceOpnRead(ctx, d, m)
}

func resourceOpnDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "deleting opn")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	opnId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Invalid OPN id: %v %s", d.Id(), err)
	}

	tflog.Debug(ctx, "calling ODK NetworkingApi.OpnsGet")
	params := map[string]interface{}{
		"pageSize": int32(math.MaxInt16),
	}
	opn, _, err := client.NetworkingApi.OpnsGet_1(*auth, int32(opnId), params)
	if err != nil {
		return diag.Errorf("ODK Error in NetworkingApi.OpnsGet. %s", err)
	}

	var instanceIds []int32
	for _, privateIp := range opn.PrivateIps {
		instanceIds = append(instanceIds, privateIp.Instance.Id)
	}

	err = detachInstancesFromOpn(client, auth, instanceIds, int32(opnId))
	if err != nil {
		return diag.Errorf("%s", err)
	}

	tflog.Debug(ctx, "calling ODK NetworkingApi.OpnsDelete")
	ticket, _, err := client.NetworkingApi.OpnsDelete(*auth, int32(opnId))
	if err != nil {
		return diag.Errorf("ODK Error in NetworkingApi.OpnsDelete. %s", err)
	}

	deleteTicket, err := waitForTicket(client, auth, ticket)
	if err != nil {
		return diag.Errorf("ODK Error in TicketsApi.TicketsGet. %s", err)
	}
	if deleteTicket.Status.Id != DICT_TICKET_SUCCEED {
		return diag.Errorf("Unable to delete OPN. Ticket status=%v", deleteTicket.Status.Id)
	}

	d.SetId("")
	return nil
}

func loadOpnData(ctx context.Context, d *schema.ResourceData, m interface{}, opn odk.Opn) diag.Diagnostics {
	// Load instances list
	instances := make([]int32, 0)
	for _, ip := range opn.PrivateIps {
		instances = append(instances, ip.Instance.Id)
	}

	// Store everything
	tflog.Debug(ctx, "Parsing returned data")
	if d.Set("name", opn.Name) != nil {
		return diag.Errorf("Can't retrieve opn name")
	}
	if d.Set("creation_user_id", opn.CreationUser.Id) != nil {
		return diag.Errorf("Can't retrieve creation user id")
	}
	if d.Set("creation_date", opn.CreationDate.String()) != nil {
		return diag.Errorf("Can't retrieve creation date")
	}
	if d.Set("last_change_date", opn.LastChangeDate.String()) != nil {
		return diag.Errorf("Can't retrieve last change date")
	}
	if d.Set("instance_ids", instances) != nil {
		return diag.Errorf("Can't retrieve instances list")
	}
	return nil
}

func detachInstancesFromOpn(client odk.APIClient, auth *context.Context, instancesIds []int32, opnId int32) error {
	for _, instanceId := range instancesIds {
		detachCommand := odk.DetachInstanceFromOpnCommand{
			OpnId: opnId,
		}
		tflog.Debug(context.Background(), "calling ODK OCIInterfacesApi.InstancesDetachFromOpn")
		ticket, resp, err := client.OCIInterfacesApi.InstancesDetachFromOpn(*auth, instanceId, detachCommand)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("instance id %v or OPN id %v was not found", instanceId, opnId)
			}
			return fmt.Errorf("ODK Error in OCIInterfacesApi.InstancesDetachFromOpn. %s", err)
		}
		detachTicket, err := waitForTicket(client, auth, ticket)
		if err != nil {
			return fmt.Errorf("ODK Error in TicketsApi.TicketsGet. %s", err)
		}
		if detachTicket.Status.Id != DICT_TICKET_SUCCEED {
			return fmt.Errorf("unable to detach instance. Ticket status=%v", detachTicket.Status.Id)
		}
	}
	return nil
}
