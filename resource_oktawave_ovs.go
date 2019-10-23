package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/oktawave-code/odk"
	"log"
	"net/http"
	"strconv"
	"time"
)

func resourceOvs() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvsCreate,
		Read:   resourceOvsRead,
		Update: resourceOvsUpdate,
		Delete: resourceOvsDelete,
		Schema: map[string]*schema.Schema{
			"disk_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			//ask macrin about default
			"space_capacity": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"tier_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"is_shared": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"shared_disk_type_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"subregion_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"is_locked": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"connections_with_instanceids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"isfreemium": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},
	}
}

func resourceOvsCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource OVS. CREATE. Initializing")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx

	log.Printf("[DEBUG] Resource OVS. CREATE. Retrieving attributes from config file")

	diskName := d.Get("disk_name").(string)
	spaceCapacity := (int32)(d.Get("space_capacity").(int))
	tierId := (int32)(d.Get("tier_id").(int))
	isShared := d.Get("is_shared").(bool)
	sharedDiskTypeId := (int32)(d.Get("shared_disk_type_id").(int))
	subregionId := (int32)(d.Get("subregion_id").(int))
	idConnectedInstances, idsExist := d.GetOk("connections_with_instanceids")
	idConnectedInstancesInt32 := retrieve_ids(idConnectedInstances.(*schema.Set).List())

	log.Printf("[INFO] Resource OVS. CREATE. Setting up command for creating OVS")
	createCommand := odk.CreateDiskCommand{
		DiskName:         diskName,
		SpaceCapacity:    spaceCapacity,
		TierId:           tierId,
		IsShared:         isShared,
		SharedDiskTypeId: sharedDiskTypeId,
		SubregionId:      subregionId,
	}

	if idsExist {
		if !isShared && len(idConnectedInstancesInt32) > 1 {
			return fmt.Errorf("Resource OVS. CREATE. Unshared disk can be attached only to 1 or 0 volumes")
		}
		createCommand.InstanceIdsList = idConnectedInstancesInt32
	}
	if isShared && !(sharedDiskTypeId == 1411 || sharedDiskTypeId == 1412) {
		return fmt.Errorf("Resource OVS. CREATE. If shared status is set true, shared disk type id should be set to 1411 or 1412")
	}

	instanceToOn, err := powerOffInstances(client, auth, idConnectedInstancesInt32)
	if err != nil {
		return fmt.Errorf("Resource OVS. CREATE. %s", err)
	}
	log.Printf("[DEBUG] Resource OVS. CREATE. Trying to post disk")
	ticket, _, err := client.OVSApi.DisksPost(*auth, createCommand)
	if err != nil {
		return fmt.Errorf("Resource OVS. CREATE. Cannot post new disk. %v", err)
	}

	log.Printf("[DEBUG] Resource OVS. CREATE. Waiting for ticket..")
	createTicket, err := evaluateTicket(client, auth, ticket)
	if err != nil {
		return fmt.Errorf("Resource OVS. CREATE. Can't get create ticket %s", err)
	}
	switch createTicket.Status.Id {
	//error case
	case TICKET_STATUS__ERROR:
		apiTicketStatusId := int(createTicket.Status.Id)
		return fmt.Errorf("Resource OVS. CREATE. Unable to create disk. Check, whether you set shared type id if shared status set true. Ticket status: %s", strconv.Itoa(apiTicketStatusId))
	//success case
	case TICKET_STATUS__SUCCESS:
		log.Printf("[DEBUG] Resource OVS. CREATE. Disk was created")
		d.SetId(strconv.Itoa(int(createTicket.ObjectId)))
	}
	err = powerOnInstances(client, auth, instanceToOn)
	if err != nil {
		return fmt.Errorf("Resource OVS. CREATE. %s", err)
	}
	return resourceOvsRead(d, m)
}

func resourceOvsRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Resource OVS. READ. Initializing")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Resource OVS. READ. Invalid volume id: %v", err)
	}

	log.Print("[DEBUG] Resource OVS. READ. Context ", *auth, " Id: ", id)
	disk, resp, err := client.OVSApi.DisksGet(*auth, (int32)(id), nil)

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return fmt.Errorf("Resource OVS. READ. Resource OVS by id=%s was not found", strconv.Itoa(id))
		}
		return fmt.Errorf("Resource OVS. READ. Error retrieving volume: %s", err)
	}

	log.Printf("[DEBUG] Resource OVS. READ. Trying to retrieve OVS disk name")
	if err := d.Set("disk_name", disk.Name); err != nil {
		return fmt.Errorf("Resource OVS. READ. Error: can't retrieve OVS disk name")
	}

	log.Printf("[DEBUG] Resource OVS. READ. Trying to retrieve OVS disk space capacity")
	if err := d.Set("space_capacity", disk.SpaceCapacity); err != nil {
		return fmt.Errorf("Resource OVS. READ. Error: can't retrieve OVS disk space capacity")
	}

	log.Printf("[DEBUG] Resource OVS. READ. Trying to retrieve OVS disk tier id")
	if err := d.Set("tier_id", disk.Tier.Id); err != nil {
		return fmt.Errorf("Resource OVS. READ. Error: can't retrieve OVS disk tier id")
	}

	log.Printf("[DEBUG] Resource OVS. READ. Trying to retrieve OVS disk shared option")
	if err := d.Set("is_shared", disk.IsShared); err != nil {
		return fmt.Errorf("Resource OVS. READ. Error: can't retrieve OVS disk shared option")
	}

	//shared disk type is nullable, so we need to check, whether this attribute was set via terraform before
	if _, ok := d.GetOk("shared_disk_type_id"); ok {
		log.Printf("[DEBUG] Resource OVS. READ. Trying to retrieve OVS disk shared disk type id")
		if err := d.Set("shared_disk_type_id", disk.SharedDiskType.Id); err != nil && disk.SharedDiskType.Label != "" {
			return fmt.Errorf("Resource OVS. READ. Error: can't retrieve OVS shared disk type id")
		}
	}

	log.Printf("[DEBUG] Resource OVS. READ. Trying to retrieve OVS disk subregion id")
	if err := d.Set("subregion_id", disk.Subregion.Id); err != nil {
		return fmt.Errorf("Resource OVS. READ. Error: can't retrieve OVS disk subregion id")
	}

	if err := d.Set("is_locked", disk.IsLocked); err != nil {
		return fmt.Errorf("Resource OVS. READ. Error: can't retrieve OVS disk locked option")
	}

	log.Printf("[DEBUG] Resource OVS. READ. Trying to retrieve OVS connected to volume instance ids")
	if err := d.Set("connections_with_instanceids", getConnectionsInstancesIds(disk.Connections)); err != nil {
		return fmt.Errorf("Resource OVS. READ. Error: can't retrieve instance ids connected with OVS disk. %s", err)
	}

	log.Printf("[DEBUG] Resource OVS. READ. Trying to retrieve OVS volume subregion option")
	if err := d.Set("isfreemium", disk.IsFreemium); err != nil {
		return fmt.Errorf("Resource OVS. READ. Error: can't retrieve OVS disk freemium option")
	}

	return nil
}
func resourceOvsUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource OVS. Update. Initializaing")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx
	diskId, _ := strconv.Atoi(d.Id())

	d.Partial(true)

	newConnections := d.Get("connections_with_instanceids")
	connectionsId_Int32 := retrieve_ids(newConnections.(*schema.Set).List())
	updateDiskCmd := odk.UpdateDiskCommand{
		DiskName:      d.Get("disk_name").(string),
		SpaceCapacity: int32(d.Get("space_capacity").(int)),
		TierId:        int32(d.Get("tier_id").(int)),
	}

	if d.HasChange("subregion_id") {
		oldSubregion, newSubregion := d.GetChange("subregion_id")
		updateDiskCmd.SubregionId = int32(oldSubregion.(int))
		updateDiskCmd.InstanceIdsList = nil
		if err := updateDisk(client, auth, int32(diskId), updateDiskCmd); err != nil {
			return err
		}
		d.SetPartial("subregion_id")
		updateDiskCmd.SubregionId = int32(newSubregion.(int))
		if err := updateDisk(client, auth, int32(diskId), updateDiskCmd); err != nil {
			return err
		}
	}

	updateDiskCmd.SubregionId = int32(d.Get("subregion_id").(int))
	updateDiskCmd.InstanceIdsList = connectionsId_Int32
	if err := updateDisk(client, auth, int32(diskId), updateDiskCmd); err != nil {
		return err
	}

	d.SetPartial("connections_with_instanceids")
	d.SetPartial("disk_name")
	d.SetPartial("space_capacity")
	d.SetPartial("tier_id")
	d.Partial(false)
	return resourceOvsRead(d, m)
}
func resourceOvsDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource OVS. DELETE. Initializing")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Resource OVS. DELETE. Invalid resource ID: %s", strconv.Itoa(id))
	}

	log.Printf("[DEBUG] Resource OVS. DELETE. Trying to get disk by id %s", strconv.Itoa(id))
	deleteCommand := odk.UpdateDiskCommand{
		DiskName:        d.Get("disk_name").(string),
		SpaceCapacity:   int32(d.Get("space_capacity").(int)),
		TierId:          int32(d.Get("tier_id").(int)),
		SubregionId:     int32(d.Get("subregion_id").(int)),
		InstanceIdsList: nil,
	}

	log.Printf("[DEBUG] Resource OVS. DELETE. Trying to create detach ticket")
	detachTicket, resp, err := client.OVSApi.DisksPut(*auth, int32(id), deleteCommand)
	if err != nil {

		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Resource OVS. DELETE. Volume was not found. %s", err)
		}

		return fmt.Errorf("Resource OVS. DELETE. Error occured while trying create detach volume ticket. %s", err)
	}

	log.Printf("[DEBUG]Resource OVS. DELETE. detach ticket was created successfully.")
	log.Printf("[DEBUG]Resource OVS. DELETE. Waiting for detach ticket status")
	detachTicket, err = evaluateTicket(client, auth, detachTicket)
	switch detachTicket.Status.Id {
	//error case
	case TICKET_STATUS__ERROR:
		apiTicketStatusId := int(detachTicket.Status.Id)
		return fmt.Errorf("Resource OVS. DELETE. Unable to detach volume. Ticket status id is: %s", strconv.Itoa(apiTicketStatusId))
	}
	log.Printf("[DEBUG] Resource OVS. DELETE. Trying to create delete ticket")
	deleteVolumeTicket, resp, err := client.OVSApi.DisksDelete(*auth, int32(id))
	if err != nil {

		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Resource OVS. DELETE. Volume was not found. %s", err)
		}

		return fmt.Errorf("Resource OVS. DELETE. Error occured while trying create delete volume ticket. %s", err)
	}

	log.Printf("[DEBUG]Resource OVS. DELETE. Delete ticket was created successfully.")
	log.Printf("[DEBUG]Resource OVS. DELETE. Waiting for delete ticket status")
	deleteVolumeTicket, err = evaluateTicket(client, auth, deleteVolumeTicket)
	switch deleteVolumeTicket.Status.Id {
	//error case
	case TICKET_STATUS__ERROR:
		apiTicketStatusId := int(deleteVolumeTicket.Status.Id)
		return fmt.Errorf("Resource OVS. DELETE. Unable to delete volume. Ticket status id is: %s", strconv.Itoa(apiTicketStatusId))
		//success case
	}
	log.Printf("[INFO]Resource OVS. DELETE. OVS delete was successfull")
	return nil
}

func updateDisk(client odk.APIClient, auth *context.Context, diskId int32, updateCmd odk.UpdateDiskCommand) error {
	updateDiskCmd := updateCmd
	ticket, _, err := client.OVSApi.DisksPut(*auth, diskId, updateDiskCmd)
	if err != nil {
		return fmt.Errorf("Resource OVS. UPDATE. Cannot update connected disks. %v", err)
	}

	log.Printf("[DEBUG] Resource OVS. UPDATE. Waiting for ticket..")
	updateTicket, err := evaluateTicket(client, auth, ticket)
	if err != nil {
		return fmt.Errorf("Resource OVS. UPDATE. Can't get update ticket %s", err)
	}
	switch updateTicket.Status.Id {
	//error case
	case TICKET_STATUS__ERROR:
		apiTicketStatusId := int(updateTicket.Status.Id)
		return fmt.Errorf("Resource OVS. UPDATE. Unable to update disk. Ticket status: %s", strconv.Itoa(apiTicketStatusId))
	//success case
	case TICKET_STATUS__SUCCESS:
		log.Printf("[DEBUG] Resource OVS. UPDATE. Disk was updated")
	}
	return nil
}

func powerOffInstances(client odk.APIClient, auth *context.Context, idConnectedInstancesInt32 []int32) ([]int32, error) {
	idInstancesToPowerOn := make([]int32, 0)
	for _, id := range idConnectedInstancesInt32 {
		instance, resp, err := client.OCIApi.InstancesGet_2(*auth, id, nil)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return idInstancesToPowerOn, fmt.Errorf("Instance by id %s was not found", strconv.Itoa(int(id)))
			}
			return idInstancesToPowerOn, fmt.Errorf("Error occured whilel retrieving instance: %s", err)
		}
		if instance.MonitStatus.Id == 839 {
			ticket, resp, err := client.OCIApi.InstancesPowerOff(*auth, id)
			if err != nil {
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					return idInstancesToPowerOn, fmt.Errorf("Instance by id %s was not found", strconv.Itoa(int(id)))
				}
				return idInstancesToPowerOn, fmt.Errorf("Error occured whilel powering off instance: %s", err)
			}
			apiResponseTicket, err := evaluateTicket(client, auth, ticket)
			if err != nil {
				return idInstancesToPowerOn, fmt.Errorf("Error occured while evaluating api ticket: %s", err)
			}
			switch apiResponseTicket.Status.Id {
			case TICKET_STATUS__ERROR:
				statusCode := strconv.Itoa(int(apiResponseTicket.Status.Id))
				return idInstancesToPowerOn, fmt.Errorf("Unable to power off instance by id %s. Ticket status code id %s", err, statusCode)
			}

			idInstancesToPowerOn = append(idInstancesToPowerOn, id)
		}
	}

	return idInstancesToPowerOn, nil
}

func powerOnInstances(client odk.APIClient, auth *context.Context, instancesToOn []int32) error {
	for _, id := range instancesToOn {
		ticket, resp, err := client.OCIApi.InstancesPowerOff(*auth, id)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("Instance by id %s was not found", strconv.Itoa(int(id)))
			}
			return fmt.Errorf("Error occured whilel powering off instance: %s", err)
		}
		apiResponseTicket, err := evaluateTicket(client, auth, ticket)
		if err != nil {
			return fmt.Errorf("Error occured while evaluating api ticket: %s", err)
		}
		switch apiResponseTicket.Status.Id {
		case TICKET_STATUS__ERROR:
			statusCode := strconv.Itoa(int(apiResponseTicket.Status.Id))
			return fmt.Errorf("Unable to power off instance by id %s. Ticket status code id %s", err, statusCode)
		}
	}
	return nil
}
