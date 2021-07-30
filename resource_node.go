package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
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
	log.Printf("[INFO] Resource Kubernetes Node. CREATE. Initializing")
	client := m.(*ClientConfig).oktaOKSClient()
	auth := m.(*ClientConfig).getOKSAuth()

	iaasClient := m.(*ClientConfig).oktaClient()
	iaasAuth := m.(*ClientConfig).ctx

	typeId := d.Get("type_id").(int)
	subregionId := d.Get("subregion_id").(int)
	clusterName := d.Get("cluster_name").(string)
	nodes := prepareNodeList(float32(subregionId), float32(typeId))
	operations, resp, err := client.ClustersApi.ClustersInstancesNamePost(*auth, swagger.K44SNodesSpecification{nodes}, clusterName)

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Resource Node. CREATE. Cluster with name %s was not found", clusterName)
		}
		return fmt.Errorf("Resource Node. CREATE. Error creating node: %s", err)
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

	// tasDone, err := evaluateTask(client, *auth, clusterName, operation[0])
	// if err != nil {
	// 	return fmt.Errorf("Resource Cluster Node. CREATE. %s", err)
	// }
	nodeList, resp, err := client.ClustersApi.ClustersInstancesNameGet(*auth, clusterName)

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Resource Node. READ. Cannot get node list: cluster by name %s was not found", clusterName)
		}
		return fmt.Errorf("Resource Node. READ. Error occured while retrieving clusters nodes: %s", err)
	}

	existedNode, err := retrieveNodeByName(nodeList, ticket.ObjectName)
	if err != nil {
		return fmt.Errorf("Resource Node. READ. Erroor occured while retrieving node from clusters node list: %s", err)
	}
	d.SetId(strconv.Itoa(int(existedNode.Id)))

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

	// _, err = evaluateTask(client, *auth, clusterName, task[0])
	// if err != nil {
	// 	return fmt.Errorf("Resource Node. DELETE. Error while waiting for node destroying: %s", err)
	// }

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

//prepares node list with given subregion and type
func prepareNodeList(subregionId float32, typeId float32) []swagger.Node {
	nodeList := make([]swagger.Node, 1)
	nodeList[0] = swagger.Node{float64(subregionId), float64(typeId)}

	return nodeList
}

// func evaluateTask(client swagger.APIClient, auth context.Context, clusterName string, task swagger.K44sTaskDto) (swagger.K44sTaskDto, error) {

// 	for {
// 		foundTask, resp, err := client.ClustersApi.ClustersInstancesNameTasksTaskIdGet(auth, task.TaskId, clusterName)
// 		if err != nil {
// 			if resp != nil && resp.StatusCode == 404 {
// 				return swagger.K44sTaskDto{}, fmt.Errorf("Cluster by name %s was not found", clusterName)
// 			}
// 			return swagger.K44sTaskDto{}, fmt.Errorf("Error occured while retrieving cluster tasks: %s", err)
// 		}
// 		if strings.ToLower(foundTask.Status) == "failed" {
// 			return swagger.K44sTaskDto{}, fmt.Errorf("Task status: FAILED")
// 		}
// 		log.Printf("[DEBUG] foundTask status, id, subregion and type %s, %s, %s, %s", strings.ToLower(foundTask.Status), strconv.Itoa(int(foundTask.InstanceId)), strconv.Itoa(int(foundTask.SubregionId)), strconv.Itoa(int(foundTask.TypeId)))
// 		if foundTask.Status == "Succeeded" {
// 			return foundTask, nil
// 		}
// 		time.Sleep(10 * time.Second)
// 	}
// }
