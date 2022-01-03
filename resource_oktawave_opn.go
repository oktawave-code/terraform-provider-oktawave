package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/oktawave-code/odk"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func resourceOpn() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpnCreate,
		Read:   resourceOpnRead,
		Update: resourceOpnUpdate,
		Delete: resourceOpnDelete,
		Schema: map[string]*schema.Schema{
			"opn_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"creation_user_id": {
				Type:     schema.TypeInt,
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
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},
	}
}

func resourceOpnCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource OPN. CREATE. Initializing")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx

	log.Printf("[DEBUG] Resource OPN. CREATE. Retrieving attributes from config file")
	opnName := d.Get("opn_name").(string)
	createOpnCmd := odk.CreateOpnCommand{
		OpnName: opnName,
	}
	log.Printf("[DEBUG] Resource OPN. CREATE. Trying to create post OPN ticket")
	ticket, _, err := client.NetworkingApi.OpnsPost(*auth, createOpnCmd)
	if err != nil {
		return fmt.Errorf("Resource OPN. CREATE. Cannot post new OPN. Error: %s", err)
	}
	log.Printf("[DEBUG] Resource OPN. CREATE. Post ticket was created")
	log.Printf("[DEBUG] Resource OPN. CREATE. Waiting for ticket status")
	createTicket, err := evaluateTicket(client, auth, ticket)
	if err != nil {
		return fmt.Errorf("Resource OPN. CREATE. Can't get create ticket %s", err)
	}
	switch createTicket.Status.Id {
	//error case
	case TICKET_STATUS__ERROR:
		apiTicketStatusId := int(createTicket.Status.Id)
		return fmt.Errorf("Resource OPN. CREATE. Unable to post OPN . Ticket status: %s", strconv.Itoa(apiTicketStatusId))
	//success case
	case TICKET_STATUS__SUCCESS:
		log.Printf("[INFO]Resource OPN. CREATE. Ticket was returned with good status.")
		d.SetId(strconv.Itoa(int(createTicket.ObjectId)))
		log.Printf("[INFO] Resource OPN. CREATE. Resource Id was set - %s", d.Id())
	}

	return resourceOpnRead(d, m)
}

func resourceOpnRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource OPN. READ. Initializing.")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Resource OPN. READ. Invalid OPN id: %v", err)
	}

	log.Printf("[DEBUG] Resource OPN. READ. Trying to get list of OPNs")
	opnCollection, _, err := client.NetworkingApi.OpnsGet(*auth, nil)
	if err != nil {
		return fmt.Errorf("Resource OPN. READ. Error occured while getting list of OPNS : %s", err)
	}

	log.Printf("[DEBUG] Resource OPN. READ. Trying to find OPN by id %s on list of OPNs", d.Id())
	opn, err := findOpnById(int32(id), opnCollection)
	if err != nil {
		return fmt.Errorf("Resource OPN. READ. %s", err)
	}

	log.Printf("[DEBUG] Resource OPN. READ. Appropriate OPN was found. Trying to synchronize remote and local state..")
	if err := d.Set("opn_name", opn.Name); err != nil {
		return fmt.Errorf("Resource OPN. READ. Error occured while retrieving OPNs name %s", err)
	}

	if err := d.Set("creation_user_id", opn.CreationUser.Id); err != nil {
		return fmt.Errorf("Resource OPN. READ. Error occured while retrieving OPNs creation user id %s", err)
	}

	if err := d.Set("creation_date", opn.CreationDate.String()); err != nil {
		return fmt.Errorf("Resource OPN. READ. Error occured while retrieving OPNs creation date %s", err)
	}

	if err := d.Set("last_change_date", opn.LastChangeDate.String()); err != nil {
		return fmt.Errorf("Resource OPN. READ. Error occured while retrieving OPNs last change date %s", err)
	}

	return nil
}

func resourceOpnUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource OPN. UPDATE. Initializing.")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Resource OPN. UPDATE. Invalid OPN id: %v", err)
	}
	d.Partial(true)
	if d.HasChange("opn_name") {
		log.Printf("Resource OPN. UPDATE. Opn name attribute change detected. Updating..")
		opnName := d.Get("opn_name").(string)
		updOPNCmd := odk.UpdateOpnCommand{
			OpnName: opnName,
		}

		_, resp, err := client.NetworkingApi.OpnsPut(*auth, int32(id), updOPNCmd)
		if err != nil && !strings.Contains(err.Error(), "EOF") {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("Resource OPN. UPDATE. OPN was not found by id %s", d.Id())
			}
			return fmt.Errorf("Resource OPN. UPDATE. Error occured while updating opn name: %s", err)
		}
		d.SetPartial("opn_name")
	}

	d.Partial(false)
	return resourceOpnRead(d, m)
}

//TODO: OPN won't delete if there is instances attached to it - find the way, how to firstly detach instances from OPN and after that delete OPN
func resourceOpnDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource OPN. DELETE. Initializing.")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Resource OPN. DELETE. Invalid OPN id: %v", err)
	}

	log.Printf("[DEBUG] Looking for opns instances attachment to detach")
	opnCollection, _, err := client.NetworkingApi.OpnsGet(*auth, nil)
	if err != nil {
		return fmt.Errorf("Resource OPN. DELETE. Error occured while getting list of OPNS : %s", err)
	}

	log.Printf("[DEBUG] Resource OPN. DELETE. Trying to find OPN by id %s on list of OPNs", d.Id())
	opn, err := findOpnById(int32(id), opnCollection)
	if err != nil {
		return fmt.Errorf("Resource OPN. DELETE. %s", err)
	}
	instanceIds := retrieveInstanceFromPrivateIps_int32(opn.PrivateIps)
	err = detachInstancesFromOpn(client, auth, instanceIds, int32(id))
	if err != nil {
		return fmt.Errorf("Resource OPN. DELETE. %s", err)
	}

	log.Printf("[DEBUG] Resource OPN. DELETE. Trying to create OPN delete ticket")
	ticket, resp, err := client.NetworkingApi.OpnsDelete(*auth, int32(id))

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Resource OPN. DELETE. OPN by id %s was not found", d.Id())
		}
		return fmt.Errorf("Resource OPN. DELETE. Error occured while trying to post OPN delete ticket: %s", err)
	}

	responseTicket, err := evaluateTicket(client, auth, ticket)

	if err != nil {
		return fmt.Errorf("Resource OPN. DELETE. Can't get delete ticket %s", err)
	}

	switch responseTicket.Status.Id {
	//error case
	case TICKET_STATUS__ERROR:
		apiTicketStatusId := int(responseTicket.Id)
		return fmt.Errorf("Resource OPN. DELETE. Unable to delete OPN . Ticket status: %s", strconv.Itoa(apiTicketStatusId))
	//succeed case
	case TICKET_STATUS__SUCCESS:
		log.Printf("Resource OPN. DELETE. OPN was successfully deleted")
		d.SetId("")
	}
	return nil
}

func findOpnById(id int32, collection odk.ApiCollectionOpn) (odk.Opn, error) {
	opnCollection := collection.Items
	//iterating over list of opns and when find appropriate - return it.
	for i := 0; i < len(opnCollection); i++ {
		foundOpn := opnCollection[i]
		if foundOpn.Id == id {
			return foundOpn, nil
		}
	}
	//if opn was not found on list - return empty structure, and error
	return *(new(odk.Opn)), fmt.Errorf("Resource OPN. OPN by id %s is not present on list of OPNs", strconv.Itoa(int(id)))
}

func detachInstancesFromOpn(client odk.APIClient, auth *context.Context, instancesIds []int32, opnId int32) error {

	//iterating over list of instance ids and detach instances to opn by its id
	for _, instanceId := range instancesIds {
		detachCmd := odk.DetachInstanceFromOpnCommand{
			OpnId: opnId,
		}
		ticket, resp, err := client.OCIInterfacesApi.InstancesDetachFromOpn(*auth, instanceId, detachCmd)
		instanceId_string := strconv.Itoa(int(instanceId))
		opnId_string := strconv.Itoa(int(opnId))
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("OPN by id %s or instance by id %s was not found", opnId_string, instanceId_string)
			}
			return fmt.Errorf("Error occured while creating ticket for detaching instance to opn %s", err)
		}
		respTicket, err := evaluateTicket(client, auth, ticket)
		switch respTicket.Status.Id {
		case TICKET_STATUS__ERROR:
			ticketStatus := strconv.Itoa(int(respTicket.Status.Id))
			return fmt.Errorf("Unable to detach instance by id %s to opn by id %s. Ticket status is: %s", instanceId_string, opnId_string, ticketStatus)
		}
	}
	return nil
}

//iterating over list of private ips, retrieved from opn and retrieving it's instance id
func retrieveInstanceFromPrivateIps_int32(privateIps []odk.PrivateIp) []int32 {
	var instanceIds []int32
	for _, privateIp := range privateIps {
		instanceIds = append(instanceIds, privateIp.Instance.Id)
	}
	return instanceIds
}
