package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/oktawave-code/oks-sdk"
	"log"
	"net/http"
	"time"
)

func resourceKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesClusterCreate,
		Read:   resourceKubernetesClusterRead,
		Delete: resourceKubernetesClusterDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"is_running": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},
	}
}

func resourceKubernetesClusterCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource Kubernetes Cluster. CREATE. Initializing")
	client := m.(*ClientConfig).oktaOKSClient()
	auth := m.(*ClientConfig).getOKSAuth()

	version := d.Get("version").(string)
	name := d.Get("name").(string)
	createClusterCMD := swagger.K44sClusterCreateDto{
		Version: version,
	}
	details, resp, err := client.ClustersApi.ClustersNamePost(*auth, createClusterCMD, name)
	log.Printf("[DEBUG] Resource Kluster. CREATE. Returned cluster name: %s", details.Name)
	if err != nil {
		return fmt.Errorf("Resource Kubernetes Cluster. CREATE. Error occured while posting cluster: %s", err)
	}

	log.Printf("[DEBUG] Resource Kubernetes Cluster. CREATE. POST response is %s", resp.Status)
	d.SetId(details.Name)
	return resourceKubernetesClusterRead(d, m)
}

func resourceKubernetesClusterRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource Kubernetes Cluster. READ. Initializing")
	client := m.(*ClientConfig).oktaOKSClient()
	auth := m.(*ClientConfig).getOKSAuth()

	name := d.Id()

	log.Printf("[INFO] Resource Kubernetes Cluster. READ. Getting clusters remote state")
	details, resp, err := client.ClustersApi.ClustersNameGet(*auth, name)
	log.Printf("[DEBUG] Resource Kubernetes Cluster. READ. Version %s", details.Version)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return fmt.Errorf("Resource Kubernetes Cluster. READ. Cluster by the name %s was not found", name)
		}
		return fmt.Errorf("Resource Kubernetes Cluster. READ. Error occured while retrieving remote Kubernetes cluster state: %s", err)
	}

	if err := d.Set("version", details.Version); err != nil {
		return fmt.Errorf("Resource Kubernetes Cluster. READ. Error occured while retrieving remote cluster version")
	}

	if err := d.Set("is_running", details.Running); err != nil {
		return fmt.Errorf("Resource Kubernetes Cluster. READ. Error occured while retrieving clusters running option")
	}
	return nil
}

func resourceKubernetesClusterDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource Kubernetes Cluster. DELETE. Initializing")
	client := m.(*ClientConfig).oktaOKSClient()
	auth := m.(*ClientConfig).getOKSAuth()

	name := d.Id()
	_, resp, err := client.ClustersApi.ClustersNameDelete(*auth, name)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Resource Kubernetes Cluster. DELETE. Cluster was not found by the name %s", name)
		}
		return fmt.Errorf("Resource Kubernetes Cluster. DELETE. Error occured while deleting cluster.  Please delete cluster nodes first")
	}

	d.SetId("")

	return nil
}
