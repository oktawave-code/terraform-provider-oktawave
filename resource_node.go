package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
	swagger "github.com/oktawave-code/oks-sdk"
)

func resourceNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceNodeCreate,
		Read:   resourceNodeRead,
		Delete: resourceNodeDelete,
		Schema: map[string]*schema.Schema{
			"type_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"subregion_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"node_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},
	}
}

func resourceNodeCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[TRACE] New OKS worker node requested.")
	client := m.(*ClientConfig).oktaOKSClient()
	auth := m.(*ClientConfig).getOKSAuth()

	iaasClient := m.(*ClientConfig).oktaClient()
	iaasAuth := m.(*ClientConfig).ctx

	typeId := d.Get("type_id").(int)
	subregionId := d.Get("subregion_id").(int)
	clusterName := d.Get("cluster_name").(string)

	nodes := make([]swagger.Node, 1)
	nodes[0] = swagger.Node{Subregion: float64(subregionId), Type_: float64(typeId)}

	log.Printf("[TRACE] Calling API for new OKS worker node. Cluster name: [%s], Subregion id: [%d], Type id: [%d], Nodes number: [%d]", clusterName, subregionId, typeId, len(nodes))
	operations, resp, err := client.ClustersApi.ClustersInstancesNamePost(*auth, swagger.K44SNodesSpecification{nodes}, clusterName)

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Could not create requested OKS worker node. Cluster with name [%s] was not found", clusterName)
		}
		return fmt.Errorf("Could not create requested OKS worker node. Cluster name: [%s], Unexpected error: %s", clusterName, err)
	}

	if len(operations) < 1 || len(operations[0].Error_) > 0 {
		return fmt.Errorf("Could not create requested OKS worker node. Cluster name: [%s], Unexpected error: %s", clusterName, operations[0].Error_)
	}

	ticketId := operations[0].Ticket.Id
	ticket := odk.Ticket{
		Id: int64(ticketId),
	}

	ticket, err = evaluateTicket(iaasClient, iaasAuth, ticket)
	if err != nil {
		return fmt.Errorf("Could not create requested OKS worker node. Ticket poll failed. Returned error: [%s]", err)
	}

	if ticket.Status.Id == TICKET_STATUS__ERROR {
		return fmt.Errorf("Failed to create OKS worker node. Ticket operation failed. Ticket id: [%d], Ticket status id: [%d]", ticket.Id, ticket.Status.Id)
	}
	d.SetId(strconv.Itoa(int(ticket.ObjectId)))

	return resourceNodeRead(d, m)
}

func resourceNodeRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource Kubernetes Node. Read. Initializing")
	client := m.(*ClientConfig).oktaOKSClient()
	auth := m.(*ClientConfig).getOKSAuth()

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Resource Kubernetes Node.  READ. Error occured while converting id from string to int: %s", err)
	}
	clusterName := d.Get("cluster_name").(string)
	nodeList, resp, err := client.ClustersApi.ClustersInstancesNameGet(*auth, clusterName)

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Resource Node. READ. Cannot get node list: cluster by name %s was not found", clusterName)
		}
		return fmt.Errorf("Resource Node. READ. Error occured while retrieving clusters nodes: %s", err)
	}

	existedNode, err := retrieveNodeById(nodeList, id)
	if err != nil {
		return fmt.Errorf("Resource Node. READ. Erroor occured while retrieving node from clusters node list: %s", err)
	}

	if err := d.Set("type_id", existedNode.Type_.Id); err != nil {
		return fmt.Errorf("Resource Kubernetes Node. Read. Error occured while refreshing local state of type id: %s", err)
	}
	if err := d.Set("subregion_id", existedNode.Subregion.Id); err != nil {
		return fmt.Errorf("Resource Kubernetes Node. Read. Error occured while refreshing local state of subregion id: %s", err)
	}
	if err := d.Set("node_name", existedNode.Name); err != nil {
		return fmt.Errorf("Resource Kubernetes Node. Read. Error occured while refreshing local state of node name: %s", err)
	}

	return nil
}

func resourceNodeDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource Kubernetes Node. Delete. Initializing")
	client := m.(*ClientConfig).oktaOKSClient()
	auth := m.(*ClientConfig).getOKSAuth()
	iaasClient := m.(*ClientConfig).oktaClient()
	iaasAuth := m.(*ClientConfig).ctx
	clusterName := d.Get("cluster_name").(string)
	id, err := strconv.Atoi(d.Id())
	nodeList := swagger.K44SNodesList{[]float64{float64(id)}}
	operations, resp, err := client.ClustersApi.ClustersInstancesNameDelete(*auth, nodeList, clusterName)

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Resource Node. DELETE. Error deleting Nodes: cluster by name %s was not found", clusterName)
		}
		return fmt.Errorf("Resource Node. DELETE. Error deleting Nodes: %s", err)
	}

	if len(operations[0].Error_) > 0 {
		return fmt.Errorf("Resource Node. CREATE. Error creating node: %s", operations[0].Error_)
	}

	ticketId := operations[0].Ticket.Id
	ticket := odk.Ticket{
		Id: int64(ticketId),
	}

	log.Printf("[INFO] Resource Node. CREATE. Post instance ticket was successfull created.")
	log.Printf("[INFO] Resource Node. CREATE. Waiting for ticket progress = 100..")
	//waiting for ticket status
	ticket, err = evaluateTicket(iaasClient, iaasAuth, ticket)
	if err != nil {
		return fmt.Errorf("Resource Node. CREATE. Ticket retrieval error. %s", err)
	}

	if ticket.Status.Id == TICKET_STATUS__ERROR {
		apiTicketStatusId := int(ticket.Status.Id)
		return fmt.Errorf("Resource Node. CREATE. Unable to create instance. Ticket status: %s", strconv.Itoa(apiTicketStatusId))
	}

	d.SetId("")
	return nil
}
