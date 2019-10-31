package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	swagger "github.com/oktawave-code/oks-sdk"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
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

	typeId := d.Get("type_id").(int)
	subregionId := d.Get("subregion_id").(int)
	clusterName := d.Get("cluster_name").(string)
	nodes := prepareNodeList(float32(subregionId), float32(typeId))
	details, resp, err := client.ClustersApi.ClustersInstancesNamePost(*auth, swagger.K44SNodesSpecification{nodes}, clusterName)

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Resource Node. CREATE. Cluster by name %s was not found", clusterName)
		}
		return fmt.Errorf("Resource Node. CREATE. Error creating node: %s", err)
	}

	tasDone, err := evaluateTask(client, *auth, clusterName, details[0])
	if err != nil {
		return fmt.Errorf("Resource Cluster Node. CREATE. %s", err)
	}
	nodeList, resp, err := client.ClustersApi.ClustersInstancesNameGet(*auth, clusterName)

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Resource Node. READ. Cannot get node list: cluster by name %s was not found", clusterName)
		}
		return fmt.Errorf("Resource Node. READ. Error occured while retrieving clusters nodes: %s", err)
	}

	existedNode, err := retrieveNodeByName(nodeList, tasDone.InstanceName)
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
	clusterName := d.Get("cluster_name").(string)
	id, err := strconv.Atoi(d.Id())
	nodeList := swagger.K44SNodesList{[]float32{float32(id)}}
	task, resp, err := client.ClustersApi.ClustersInstancesNameDelete(*auth, nodeList, clusterName)

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Resource Node. DELETE. Error deleting Nodes: cluster by name %s was not found", clusterName)
		}
		return fmt.Errorf("Resource Node. DELETE. Error deleting Nodes: %s", err)
	}
	_, err = evaluateTask(client, *auth, clusterName, task[0])
	if err != nil {
		return fmt.Errorf("Resource Node. DELETE. Error while waiting for node destroying: %s", err)
	}
	d.SetId("")
	return nil
}

//prepares node list with given subregion and type
func prepareNodeList(subregionId float32, typeId float32) []swagger.Node {
	nodeList := make([]swagger.Node, 1)
	nodeList[0] = swagger.Node{subregionId, typeId}

	return nodeList
}

func evaluateTask(client swagger.APIClient, auth context.Context, clusterName string, task swagger.K44sTaskDto) (swagger.K44sTaskDto, error) {

	for {
		foundTask, resp, err := client.ClustersApi.ClustersInstancesNameTasksTaskIdGet(auth, task.TaskId, clusterName)
		if err != nil {
			if resp != nil && resp.StatusCode == 404 {
				return swagger.K44sTaskDto{}, fmt.Errorf("Cluster by name %s was not found", clusterName)
			}
			return swagger.K44sTaskDto{}, fmt.Errorf("Error occured while retrieving cluster tasks: %s", err)
		}
		if strings.ToLower(foundTask.Status) == "failed" {
			return swagger.K44sTaskDto{}, fmt.Errorf("Task status: FAILED")
		}
		log.Printf("[DEBUG] foundTask status, id, subregion and type %s, %s, %s, %s", strings.ToLower(foundTask.Status), strconv.Itoa(int(foundTask.InstanceId)), strconv.Itoa(int(foundTask.SubregionId)), strconv.Itoa(int(foundTask.TypeId)))
		if foundTask.Status == "Succeeded" {
			return foundTask, nil
		}
		time.Sleep(10 * time.Second)
	}
}
