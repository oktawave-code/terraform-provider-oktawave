package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

//TODO: Implement subregion migration operation. For now - unsupported operation
//TODO: Implement retries in CRUD - unsupported
func resourceOci() *schema.Resource {
	return &schema.Resource{
		Create: resourceOciCreate,
		Read:   resourceOciRead,
		Update: resourceOciUpdate,
		Delete: resourceOciDelete,

		Schema: map[string]*schema.Schema{
			"authorization_method_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1399,
			},
			"init_disk_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"opn_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"disk_class": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"init_disk_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instances_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"without_public_ip": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"init_ip_address_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ssh_keys_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"ip_address_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"subregion_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"template_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"type_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"isfreemium": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"init_script": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"opn_mac": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Computed: true,
					Type:     schema.TypeString,
				},
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"islocked": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"creation_userid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"init_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},
	}
}

func resourceOciCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource OCI. CREATE. Initializing.")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx
	authorizationMethod := (int32)(d.Get("authorization_method_id").(int))

	log.Printf("[INFO] Resource OCI. CREATE. setting create instance OCI command parameters")
	createCommand := odk.CreateInstanceCommand{
		AuthorizationMethodId: authorizationMethod,
		DiskClass:             (int32)(d.Get("disk_class").(int)),
		DiskSize:              (int32)(d.Get("init_disk_size").(int)),
		InstanceName:          d.Get("instance_name").(string),
		InstancesCount:        (int32)(d.Get("instances_count").(int)),
		IPAddressId:           (int32)(d.Get("init_ip_address_id").(int)),
		SubregionId:           (int32)(d.Get("subregion_id").(int)),
		TemplateId:            (int32)(d.Get("template_id").(int)),
		TypeId:                (int32)(d.Get("type_id").(int)),
		Freemium:              d.Get("isfreemium").(bool),
		InitScript:            d.Get("init_script").(string),
		WithoutPublicIp:       (bool)(d.Get("without_public_ip").(bool)),
		OpnsIds:               retrieve_ids(d.Get("opn_ids").(*schema.Set).List()),
	}

	if authorizationMethod == 1398 {
		sshKeyInterface := d.Get("ssh_keys_ids")
		sshKeyIds := retrieve_ids(sshKeyInterface.(*schema.Set).List())
		if len(sshKeyIds) == 0 {
			return fmt.Errorf("Resource OCI. CREATE. Error occured while configuring OCI creation: " +
				"When authorization method == ssh keys, non empty ssh keys ids list expected")
		}
		sshKeys, _, err := client.AccountApi.AccountGetSshKeys(*auth, nil)
		if err != nil {
			return fmt.Errorf("Resource OCI. CREATE. Error occured while getting users list of ssh keys: %s", err)
		}
		if err := checkSshKeysList(sshKeyIds, sshKeys.Items); err != nil {
			return err
		}

		createCommand.SshKeysIds = sshKeyIds
	}

	log.Printf("[INFO] Resource OCI. CREATE. trying to post instance")
	ticket, _, err := client.OCIApi.InstancesPost(*auth, createCommand)
	if err != nil {
		return fmt.Errorf("Resource OCI. CREATE. Error: cannot post instance. Msg: %s", err)
	}
	log.Printf("[INFO] Resource. OCI. CREATE. Post instance ticket was successfull created.")
	log.Printf("[INFO] Resource.OCI. CREATE. Waiting for ticket progress = 100..")
	//waiting for ticket status
	createTicket, err := evaluateTicket(client, auth, ticket)
	if err != nil {
		return fmt.Errorf("Resource OCI. CREATE. Can't get create ticket %s", err)
	}
	switch createTicket.Status.Id {
	//error case
	case TICKET_STATUS__ERROR:
		apiTicketStatusId := int(createTicket.Status.Id)
		return fmt.Errorf("Resource OCI. CREATE. Unable to create instance . Ticket progress: %s", strconv.Itoa(apiTicketStatusId))
	//success case
	case TICKET_STATUS__SUCCESS:
		log.Printf("[INFO]Resource OCI. CREATE. Ticket was returned with good status.")
		d.SetId(strconv.Itoa(int(createTicket.ObjectId)))
		log.Printf("[INFO] Resource OCI. CREATE. Resource Id was set - %s", d.Id())
	}

	// if opnIdSet, opnIsSet := d.GetOk("opn_ids"); opnIsSet {
	// 	opnIds := retrieve_ids(opnIdSet.(*schema.Set).List())
	// 	err := attachInstanceToOpns(client, auth, opnIds, createTicket.ObjectId)
	// 	if err != nil {
	// 		return fmt.Errorf("Resource OCI. CREATE. %s", err)
	// 	}
	// }

	if ipIdSet, ipIsSet := d.GetOk("ip_address_ids"); ipIsSet {
		ipIds := retrieve_ids(ipIdSet.(*schema.Set).List())
		err := attachInstanceToIps(client, auth, ipIds, createTicket.ObjectId)
		if err != nil {
			return fmt.Errorf("Resource OCI. CREATE. %s", err)
		}
	}

	return resourceOciRead(d, m)

}

func resourceOciRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource OCI. READ. Initializing.")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Resource OCI. READ. Invalid instance id: %v", err)
	}

	disk, resp, err := client.OCIApi.InstancesGetDisks(*auth, int32(id), nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return fmt.Errorf("Resource OCI. READ. Disk by OCI id=%s was not found", strconv.Itoa(id))
		}
		return fmt.Errorf("Resource OCI. READ. Error retrieving disk: %s", err)
	}

	instance, resp, err := client.OCIApi.InstancesGet_2(*auth, (int32)(id), nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			//recreate
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Resource OCI. READ. Error retrieving instance: %s", err)
	}

	opnMacMap, err := getOpnMacMap(client, *auth, int32(id))
	if err != nil {
		return fmt.Errorf("Resource OCI. READ. %s", err)
	}

	log.Printf("[DEBUG] Resource OCI. READ. Trying to retrieve instance name")
	if err := d.Set("instance_name", instance.Name); err != nil {
		return fmt.Errorf("Error: can't retrieve instance name")
	}

	log.Printf("[DEBUG] Resource OCI. READ. Trying to retrieve instance subregion id")
	if err := d.Set("subregion_id", int(instance.Subregion.Id)); err != nil {
		return fmt.Errorf("Resource OCI. READ. Error: can't retrieve instance subregion id")
	}

	log.Printf("[DEBUG] Resource OCI. READ. Trying to retrieve instance template id")
	if err := d.Set("template_id", instance.Template.Id); err != nil {
		return fmt.Errorf("Resource OCI. READ. Error: can't retrieve instance template id")
	}

	log.Printf("[DEBUG] Resource OCI. READ. Trying to retrieve instance type id")
	if err := d.Set("type_id", instance.Type_.Id); err != nil {
		return fmt.Errorf("Error: can't retrieve instance type id")
	}

	log.Printf("[DEBUG] Resource OCI. READ. Trying to retrieve instance freemium option")
	if err := d.Set("isfreemium", instance.IsFreemium); err != nil {
		return fmt.Errorf("Resource OCI. READ. Error: can't retrieve instance freemium option")
	}

	if err := d.Set("init_disk_id", int(disk.Items[0].Id)); err != nil {
		return fmt.Errorf("Resource OCI. CREATE. Error: can't retrieve instance initial disk")
	}

	log.Printf("[DEBUG] Resource OCI. READ. Trying to retrieve instance disks capacitu")
	if err := d.Set("init_disk_size", int(disk.Items[0].SpaceCapacity)); err != nil {
		return fmt.Errorf("Resource OCI. READ. Error: can't retrieve instance initial disk size")
	}

	log.Printf("[DEBUG] Resource OCI. READ. Trying to retrieve instance creation date")
	if err := d.Set("creation_date", instance.CreationDate.String()); err != nil {
		return fmt.Errorf("Resource OCI. READ. Error: can't retrieve instance creation date")
	}

	log.Printf("[DEBUG] Resource OCI. READ. Trying to retrieve instance lock option")
	if err := d.Set("islocked", instance.IsLocked); err != nil {
		return fmt.Errorf("Resource OCI. READ. Error: can't retrieve instance locked option")
	}

	log.Printf("[DEBUG] Resource OCI. READ. Trying to retrieve instance creation user id")
	if err := d.Set("creation_userid", instance.CreationUser.Id); err != nil {
		return fmt.Errorf("Resource OCI. READ. Error: can't retrieve instance creation user id")
	}

	log.Printf("[DEBUG] Resource OCI. READ. Trying to retrieve instance ip address")
	if err := d.Set("init_ip_address", instance.IpAddress); err != nil {
		return fmt.Errorf("Resource OCI. READ. Error: can't retrieve instance ip address")
	}

	log.Printf("[DEBUG] Resource OCI. READ. Trying to retrieve instance dns address")
	if err := d.Set("dns_address", instance.DnsAddress); err != nil {
		return fmt.Errorf("Resource OCI. READ. Error: can't retrieve instance dns ip address")
	}

	log.Printf("[DEBUG] Resource OCI. READ. Trying to retrieve instance initial disk class")
	if err := d.Set("disk_class", disk.Items[0].Tier.Id); err != nil {
		return fmt.Errorf("Resource OCI. READ. Error: can't retrieve instance initial disk class")
	}

	log.Printf("[DEBUG] Resource OCI. READ. Trying to retrieve instance OPN connections")

	opnIds, resp, err := client.OCIInterfacesApi.InstancesGetOpns(*auth, int32(id), nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Resource OCI. READ. Error occured while retrieving OPN connections: " +
				"Instance not found")
		}
		return fmt.Errorf("Resource OCI. READ. Error occured while retrieving instance Opns: %s", err)
	}
	opnIds_int := getOpnIds(opnIds.Items)
	if err := d.Set("opn_ids", opnIds_int); err != nil {
		return fmt.Errorf("Resource OCI. READ. Error: can't set instance opn connections")
	}

	if err := d.Set("opn_mac", opnMacMap); err != nil {
		return fmt.Errorf("Resource OCI. READ. Error occured while setting local state of opn macs. Error: %s", err)
	}
	return nil

}

func resourceOciUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource OCI. UPDATE. Initializing")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx
	instanceId, err := strconv.Atoi(d.Id())
	initDiskId := d.Get("init_disk_id").(int)
	if err != nil {
		return fmt.Errorf("Resource OCI. UPDATE. Invalid OCI id: %v", err)
	}

	log.Printf("[INFO] Resource OCI. Checking resource update status")
	if d.HasChange("instance_name") {
		log.Printf("[DEBUG] Resource OCI. UPDATE. Instance name. Mentioned change in config file")
		newName := d.Get("instance_name").(string)
		updateTicket, resp, err := client.OCIApi.InstancesChangeName(*auth, (int32)(instanceId), newName)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("Resource OCI. UPDATE. Can't find instance. %s", err)
			}
			return fmt.Errorf("Resource OCI. UPDATE. Error occured while updating instance name %s", err)
		}

		updateTicketId, err := evaluateTicket(client, auth, updateTicket)
		if err != nil {
			return fmt.Errorf("Resource OCI. UPDATE INSTANCE NAME.Can't get ticket %s", err)
		}
		switch updateTicketId.Status.Id {
		//error case
		case TICKET_STATUS__ERROR:
			apiTicketStatusId := int(updateTicketId.Status.Id)
			return fmt.Errorf("Resource OCI. UPDATE. Unable to UPDATE instance name. Ticket status: %s", strconv.Itoa(apiTicketStatusId))
		//succeed case
		case TICKET_STATUS__SUCCESS:
			log.Printf("[INFO] Resource OCI. UPDATE. Instance name. Real state was updated successfully")
		}
	}
	//not supported by now
	if d.HasChange("subregion_id") {
		//	//	newSubregionId := d.Get("subregion_id").(int)
		//	//	subregionUpdateCommand := odk.ChangeInstanceSubregionCommand{(int32)(newSubregionId)}
		//	//	updateTicket, _, err:=client.OCIApi.InstancesChangeSubregion(*auth, (int32)(instanceId), subregionUpdateCommand)
		//	//	_,updateTicketId,err := evaluateTicket(client, auth, updateTicket)
		//	//	if err!=nil{
		//	//		fmt.Errorf("Resource OCI. UPDATE. Can't get ticket %s", err)
		//	//	}
		//	//	switch updateTicketId {
		//	//	//error case
		//	//	case 137:
		//	//		apiTicketStatusId := int(updateTicketId)
		//	//		return fmt.Errorf("Resource OCI. DELETE. Unable to UPDATE instance subregion. Ticket status: %s", strconv.Itoa(apiTicketStatusId))
		//	//	//succeed case
		//	//	case 136:
		//	//
		//	//	}
		return fmt.Errorf("Resource OCI. UPDATE. Subregion changes are not supported for now")
	}

	if d.HasChange("type_id") {
		log.Printf("[DEBUG] Resource OCI. UPDATE. Type id. Mentioned change in config file")
		newTypeId := d.Get("type_id").(int)
		updateTicket, _, err := client.OCIApi.InstancesChangeType_1(*auth, (int32)(instanceId), (int32)(newTypeId))
		updateTicketId, err := evaluateTicket(client, auth, updateTicket)
		if err != nil {
			return fmt.Errorf("Resource OCI. UPDATE TYPE ID. Can't get ticket %s", err)
		}
		switch updateTicketId.Status.Id {
		//error case
		case TICKET_STATUS__ERROR:
			apiTicketStatusId := int(updateTicketId.Id)
			return fmt.Errorf("Resource OCI. UPDATE. Unable to UPDATE instance type. Ticket status: %s", strconv.Itoa(apiTicketStatusId))
		//succeed case
		case TICKET_STATUS__SUCCESS:
			log.Printf("[INFO] Resource OCI. UPDATE. Type id. Real state was updated successfully")
		}
	}
	if d.HasChange("init_ip_address_id") {
		oldIpId, newIpId := d.GetChange("init_ip_address_id")
		oldIpId_int32 := int32(oldIpId.(int))
		newIpId_int32 := int32(newIpId.(int))
		if oldIpId != 0 {
			ticket, _, err := detachIpById(client, auth, int32(instanceId), oldIpId_int32)
			if err != nil {
				return fmt.Errorf("Resource OCI. UPDATE. %s", err)
			}
			switch ticket.Status.Id {
			//error case
			case TICKET_STATUS__ERROR:
				apiTicketStatusId := int(ticket.Status.Id)
				return fmt.Errorf("Resource OCI. UPDATE. Unable to detach ip. Ticket progress: %s", strconv.Itoa(apiTicketStatusId))
			//succeed case
			case TICKET_STATUS__SUCCESS:
				log.Printf("[INFO] Resource OCI. UPDATE. Ip detaching. remote state was updated successfully")
			}
		}
		if newIpId != 0 {
			ticket, _, err := attachIpById(client, auth, int32(instanceId), newIpId_int32)
			if err != nil {

				return fmt.Errorf("Resource OCI. UPDATE. %s", err)
			}
			switch ticket.Status.Id {
			//error case
			case TICKET_STATUS__ERROR:
				apiTicketStatusId := int(ticket.Status.Id)
				return fmt.Errorf("Resource OCI. UPDATE. Unable to attach ip. Ticket status: %s", strconv.Itoa(apiTicketStatusId))
			//succeed case
			case TICKET_STATUS__SUCCESS:
				log.Printf("[INFO] Resource OCI. UPDATE. Ip attaching. remote state was updated successfully")
			}
		}
	}

	disk, resp, err := client.OVSApi.DisksGet(*auth, int32(initDiskId), nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Resource OCI. UPDATE. Disk by id %s was not found", strconv.Itoa(initDiskId))
		}
		return fmt.Errorf("Resource OCI. UPDATE. Error occured while retrieving initial instance disk: %s", err)
	}

	if d.HasChange("opn_ids") {
		oldOpnId, newOpnId := d.GetChange("opn_ids")
		oldOpnIds := oldOpnId.(*schema.Set).List()
		newOpnIds := newOpnId.(*schema.Set).List()

		oldOpnIds_int32 := retrieve_ids(oldOpnIds)
		newOpnIds_int32 := retrieve_ids(newOpnIds)

		opnIdsToDetach := getIdListToDetach(oldOpnIds_int32, newOpnIds_int32)
		opnIdsToAttach := getIdListsToAttach(oldOpnIds_int32, newOpnIds_int32)

		if err := detachInstanceFromOpns(client, auth, opnIdsToDetach, int32(instanceId)); err != nil {
			return fmt.Errorf("Resource OCI. UPDATE. %s", err)
		}

		if err := attachInstanceToOpns(client, auth, opnIdsToAttach, int32(instanceId)); err != nil {
			return fmt.Errorf("Resource OCI. UPDATE. %s", err)
		}
	}

	connectionIdList := getConnectionsInstancesIds_int32(disk.Connections)
	updCmd := odk.UpdateDiskCommand{
		DiskName:        disk.Name,
		SpaceCapacity:   disk.SpaceCapacity,
		TierId:          disk.Tier.Id,
		SubregionId:     disk.Subregion.Id,
		InstanceIdsList: connectionIdList,
	}
	if d.HasChange("init_disk_size") {
		_, newDiskSize := d.GetChange("init_disk_size")
		updCmd.SpaceCapacity = int32(newDiskSize.(int))
	}

	if d.HasChange("disk_class") {
		_, newDiskClass := d.GetChange("disk_class")
		updCmd.TierId = int32(newDiskClass.(int))
	}

	ticket, resp, err := client.OVSApi.DisksPut(*auth, disk.Id, updCmd)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Resource OCI. UPDATE. Disk by id=%s was not found", strconv.Itoa(int(disk.Id)))
		}
		return fmt.Errorf("Resource OCI. UPDATE. Error occured while updating disk class and/or size: %s", err)
	}

	respTicket, err := evaluateTicket(client, auth, ticket)
	if err != nil {
		return fmt.Errorf("Resource OCI. UPDATE. Can't get create ticket %s", err)
	}
	switch respTicket.Status.Id {
	case TICKET_STATUS__ERROR:
		apiTicketStatusId := int(respTicket.Status.Id)
		return fmt.Errorf("Resource OCI. UPDATE. Unable to update disk class and/or size. Ticket status: %s", strconv.Itoa(apiTicketStatusId))
	}

	if d.HasChange("ip_address_ids") {
		oldIpId, newIpId := d.GetChange("ip_address_ids")
		oldIpIds := oldIpId.(*schema.Set).List()
		newIpIds := newIpId.(*schema.Set).List()
		//detaching if there  was old one and attaching if there are new one

		oldIpIdList := retrieve_ids(oldIpIds)
		newIpIdList := retrieve_ids(newIpIds)

		ipIdListToDetach := getIdListToDetach(oldIpIdList, newIpIdList)
		ipIdListToAttach := getIdListsToAttach(oldIpIdList, newIpIdList)

		if err := detachInstanceFromIps(client, auth, ipIdListToDetach, int32(instanceId)); err != nil {
			return fmt.Errorf("Resource OCI. UPDATE. %s", err)
		}

		if err := attachInstanceToIps(client, auth, ipIdListToAttach, int32(instanceId)); err != nil {
			return fmt.Errorf("Resource OCI. UPDATE. %s", err)
		}

	}
	return resourceOciRead(d, m)
}

func resourceOciDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource OCI. DELETE. Initializing.")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Resource OCI. DELETE. Invalid OCI id: %v", err)
	}

	log.Printf("[INFO]Resource OCI. DELETE. Trying to create delete ticket")
	deleteTicket, _, err := client.OCIApi.InstancesDelete(*auth, (int32)(id), nil)
	if err != nil {
		return fmt.Errorf("Resource OCI. DELETE. Instance was not found to delete. Error %s", err)
	}

	log.Printf("Resource OCI. DELETE. Delete ticket was created successfully.")
	log.Printf("Resource OCI. DELETE. Waiting for delete ticket progress")
	deleteTicketId, err := evaluateTicket(client, auth, deleteTicket)
	if err != nil {
		return fmt.Errorf("Resource OCI. DELETE. Can't get delete ticket %s", err)
	}

	switch deleteTicketId.Status.Id {
	//error case
	case TICKET_STATUS__ERROR:
		apiTicketStatusId := int(deleteTicketId.Id)
		return fmt.Errorf("Resource OCI. DELETE. Unable to delete instance . Ticket status: %s", strconv.Itoa(apiTicketStatusId))
	//succeed case
	case TICKET_STATUS__SUCCESS:
		log.Printf("Resource OCI. DELETE. Instance was successfully deleted")
		d.SetId("")
	}
	return nil
}

func attachInstanceToOpns(client odk.APIClient, auth *context.Context, opnIds []int32, instanceId int32) error {
	for _, opnId := range opnIds {
		attachOpnCmd := odk.AttachInstanceToOpnCommand{
			OpnId: opnId,
		}
		ticket, resp, err := client.OCIInterfacesApi.InstancesAttachOpn(*auth, instanceId, attachOpnCmd)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("Opn by id %s or instance by id %s was not found", strconv.Itoa(int(opnId)), strconv.Itoa(int(instanceId)))
			}
			return fmt.Errorf("Error occured while creating initial instance attach opn ticket: %s", err)
		}
		respTicket, err := evaluateTicket(client, auth, ticket)
		if err != nil {
			return fmt.Errorf("Can't get attach opn ticket %s", err)
		}
		switch respTicket.Status.Id {
		case TICKET_STATUS__ERROR:
			apiTicketStatusId := int(respTicket.Status.Id)
			return fmt.Errorf("Unable to attach instance to opn. Ticket status: %s", strconv.Itoa(apiTicketStatusId))
		case TICKET_STATUS__SUCCESS:
			log.Printf("[INFO]Resource OCI. CREATE. Ticket was returned with good status.")
		}
	}

	return nil
}

func attachInstanceToIps(client odk.APIClient, auth *context.Context, ipIds []int32, instanceId int32) error {
	for _, ipId := range ipIds {
		localOptions := map[string]interface{}{
			"ipId": ipId,
		}
		ticket, resp, err := client.OCIInterfacesApi.InstancesPostAttachIpTicket(*auth, instanceId, localOptions)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("Instance Id %s or Ip Id %s was not found", strconv.Itoa(int(instanceId)), strconv.Itoa(int(ipId)))
			}
			return fmt.Errorf("Error occured while creating ip attachment ticket: %s", err)
		}
		resultTicket, err := evaluateTicket(client, auth, ticket)
		if err != nil {
			return err
		}
		if resultTicket.Status.Id == TICKET_STATUS__ERROR {
			return fmt.Errorf("Unable to attach ip to instance. Ticket status %s", strconv.Itoa(int(resultTicket.Status.Id)))
		}
	}
	return nil
}

func detachInstanceFromIps(client odk.APIClient, auth *context.Context, ipIds []int32, instanceId int32) error {
	for idx := 0; idx < len(ipIds); idx++ {
		ticket, resp, err := client.OCIInterfacesApi.InstancesPostDetachIpTicket(*auth, instanceId, ipIds[idx])
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("Instance Id %s or Ip Id %s was not found", strconv.Itoa(int(instanceId)), strconv.Itoa(int(ipIds[idx])))
			}
			return fmt.Errorf("Error occured while creating ip detachment ticket: %s", err)
		}
		resultTicket, err := evaluateTicket(client, auth, ticket)
		if err != nil {
			return err
		}
		if resultTicket.Status.Id == TICKET_STATUS__ERROR {
			return fmt.Errorf("Unable to detach ip from instance. Ticket status %s", strconv.Itoa(int(resultTicket.Status.Id)))
		}
	}
	return nil
}

func detachInstanceFromOpns(client odk.APIClient, auth *context.Context, opnIds []int32, instanceId int32) error {
	for _, opnId := range opnIds {
		detachOpnCmd := odk.DetachInstanceFromOpnCommand{
			OpnId: opnId,
		}
		ticket, resp, err := client.OCIInterfacesApi.InstancesDetachFromOpn(*auth, instanceId, detachOpnCmd)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("Opn by id %s or instance by id %s was not found", strconv.Itoa(int(opnId)), strconv.Itoa(int(instanceId)))
			}
			return fmt.Errorf("Error occured while creating initial instance attach opn ticket: %s", err)
		}
		respTicket, err := evaluateTicket(client, auth, ticket)
		if err != nil {
			return fmt.Errorf("Can't get attach opn ticket %s", err)
		}
		switch respTicket.Status.Id {
		case TICKET_STATUS__ERROR:
			apiTicketStatusId := int(respTicket.Status.Id)
			return fmt.Errorf("Unable to attach instance to opn. Ticket status: %s", strconv.Itoa(apiTicketStatusId))
		case TICKET_STATUS__SUCCESS:
			log.Printf("[INFO]Resource OCI. CREATE. Ticket was returned with good status.")
		}
	}

	return nil
}

//looking, if list of ssh keys ids from config is sublist of list of ssh keys of this user
func checkSshKeysList(sshKyesIds []int32, userSshKeys []odk.SshKey) error {
	for _, sshKey := range sshKyesIds {
		isPresent := false
		for _, userSshKey := range userSshKeys {
			if sshKey == userSshKey.Id {
				isPresent = true
				break
			}
		}
		if !isPresent {
			return fmt.Errorf("Resource OCI. CREATE. Error occured while configuring OCI creation: "+
				"ssh key by id %s is not present on the users ssh keys list", strconv.Itoa(int(sshKey)))
		}
	}
	return nil
}

func getIdListsToAttach(olds []int32, news []int32) []int32 {
	attachList := make([]int32, 0)

	for _, newId := range news {
		toAttach := true
		for _, oldId := range olds {
			if newId == oldId {
				toAttach = false
			}
		}
		if toAttach {
			attachList = append(attachList, newId)
		}
	}

	return attachList
}

func getIdListToDetach(olds []int32, news []int32) []int32 {
	detachList := make([]int32, 0)

	for _, oldId := range olds {
		toDetach := true
		for _, newId := range news {
			if oldId == newId {
				toDetach = false
				break
			}
		}
		if toDetach {
			detachList = append(detachList, oldId)
		}
	}

	return detachList
}

func getOpnIds(opns []odk.Opn) []int {
	opnIds := make([]int, len(opns))
	for i, opn := range opns {
		opnIds[i] = int(opn.Id)
	}
	return opnIds
}

func getOpnMacMap(client odk.APIClient, auth context.Context, instanceId int32) (map[string]string, error) {
	apiOpn, response, err := client.OCIInterfacesApi.InstancesGetOpns(auth, instanceId, nil)
	if err != nil {
		if response != nil && response.StatusCode != http.StatusNotFound {
			return nil, fmt.Errorf("Instance by id %s is not exist", strconv.Itoa(int(instanceId)))
		}
		return nil, fmt.Errorf("Error occured while retrieving instances(ID = %s) opn list: %s", strconv.Itoa(int(instanceId)), err)
	}
	opns := apiOpn.Items
	opnMacMap := make(map[string]string)
	for _, opn := range opns {
		for i, ip := range opn.PrivateIps {
			if ip.Instance.Id == instanceId {
				opnMacMap[opn.Name] = opn.PrivateIps[i].MacAddress
				break
			}
		}
	}
	return opnMacMap, nil
}
