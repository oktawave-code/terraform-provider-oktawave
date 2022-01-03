package main

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func resourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceLoadBalancerCreate,
		Read:   resourceLoadBalancerRead,
		Update: resourceLoadBalancerUpdate,
		Delete: resourceLoadBalancerDelete,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"ssl_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"service_type_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  43,
			},
			"port_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"target_port_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"session_persistence_type_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  46,
			},
			"load_balancer_algorithm_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  612,
			},
			"ip_version_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  115,
			},
			"health_check_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"common_persistence_for_http_and_https_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},
	}
}

func resourceLoadBalancerCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource Load Balancer. CREATE. Initializing")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx

	log.Printf("[DEBUG] Resource Load Balancer. CREATE. Retrieving attributes from config file")
	isSsl := d.Get("ssl_enabled").(bool)
	serviseTypeId := d.Get("service_type_id").(int)
	port := d.Get("port_number").(int)
	targetPort := d.Get("target_port_number").(int)
	sessionPersistId := d.Get("session_persistence_type_id").(int)
	loadBalancerAlgorithmId := d.Get("load_balancer_algorithm_id").(int)
	ipVersId := d.Get("ip_version_id").(int)
	isHealthChek := d.Get("health_check_enabled").(bool)
	isCmmnPersistHttpHttps := d.Get("common_persistence_for_http_and_https_enabled").(bool)
	groupId := d.Get("group_id").(int)

	log.Printf("[DEBUG] Resource Load Balancer. CREATE. Trying to post new load balancer")
	createCmd := odk.SetLoadBalancerCommand{
		SslEnabled:                              isSsl,
		ServiceType:                             int32(serviseTypeId),
		PortNumber:                              int32(port),
		TargetPortNumber:                        int32(targetPort),
		SessionPersistenceType:                  int32(sessionPersistId),
		LoadBalancerAlgorithm:                   int32(loadBalancerAlgorithmId),
		IpVersion:                               int32(ipVersId),
		HealthCheckEnabled:                      isHealthChek,
		CommonPersistenceForHttpAndHttpsEnabled: isCmmnPersistHttpHttps,
	}
	loadBalancer, resp, err := client.OCIGroupsApi.LoadBalancersCreate(*auth, int32(groupId), createCmd)

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Resource Load Balancer. Create. Group was not found by id %s", strconv.Itoa(groupId))
		}
		return fmt.Errorf("Resource Load Balancer. Create. Error occured while creating load balancer: %s", err)

	}

	d.SetId(strconv.Itoa(int(loadBalancer.GroupId) + 1))
	log.Printf("[DEBUG]Resource Load Balancer. Create. Load balancer was created successufully")

	return resourceLoadBalancerRead(d, m)
}

func resourceLoadBalancerRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource Load Balancer. READ. Initializing")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx
	groupId := d.Get("group_id").(int)

	log.Printf("[DEBUG] Resourece Load Balancer. READ. Trying to get load balancer")
	loadBalancer, resp, err := client.OCIGroupsApi.LoadBalancersGetLoadBalancer(*auth, int32(groupId), nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return fmt.Errorf("Load Balancer. READ. Error: load balancer for gorup by id %s was not found", strconv.Itoa(groupId))
		}
		return fmt.Errorf("Load Balancer. READ. Error occured while getting load balancer: %s", err)
	}

	log.Printf("[DEBUG] Resource Load Balancer. READ. Synchonizing local and remote state")
	if err := d.Set("group_id", int(loadBalancer.GroupId)); err != nil {
		return fmt.Errorf("Resource Load Balancer. READ. Error occured while retrieving load balancer group id. %s", err)
	}

	if err := d.Set("ssl_enabled", loadBalancer.SslEnabled); err != nil {
		return fmt.Errorf("Resource Load Balancer. READ. Error occured while retrieving load balancer ssl enabled option. %s", err)
	}

	if err := d.Set("service_type_id", int(loadBalancer.ServiceType.Id)); err != nil {
		return fmt.Errorf("Resource Load Balancer. READ. Error occured while retrieving load balancer service type id. %s", err)
	}

	if err := d.Set("port_number", int(loadBalancer.PortNumber)); err != nil {
		return fmt.Errorf("Resource Load Balancer. READ. Error occured while retrieving load balancer port number. %s", err)
	}

	if err := d.Set("target_port_number", int(loadBalancer.TargetPortNumber)); err != nil {
		return fmt.Errorf("Resource Load Balancer. READ. Error occured while retrieving load balancer target port number. %s", err)
	}

	if err := d.Set("session_persistence_type_id", int(loadBalancer.SessionPersistenceType.Id)); err != nil {
		return fmt.Errorf("Resource Load Balancer. READ. Error occured while retrieving load balancer session persistence type id. %s", err)
	}

	if err := d.Set("load_balancer_algorithm_id", int(loadBalancer.Algorithm.Id)); err != nil {
		return fmt.Errorf("Resource Load Balancer. READ. Error occured while retrieving load balancer algorithm id. %s", err)
	}

	if err := d.Set("ip_version_id", int(loadBalancer.IpVersion.Id)); err != nil {
		return fmt.Errorf("Resource Load Balancer. READ. Error occured while retrieving load balancer IP version id. %s", err)
	}

	if err := d.Set("health_check_enabled", loadBalancer.HealthCheckEnabled); err != nil {
		return fmt.Errorf("Resource Load Balancer. READ. Error occured while retrieving load balancer health check enable option. %s", err)
	}

	if err := d.Set("common_persistence_for_http_and_https_enabled", loadBalancer.CommonPersistenceForHttpAndHttpsEnabled); err != nil {
		return fmt.Errorf("Resource Load Balancer. READ. Error occured while retrieving load balancer common persistence http and https option. %s", err)
	}

	return nil
}

func resourceLoadBalancerUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource Load Balancer. UPDATE. Initializing")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx
	groupId, _ := d.GetChange("group_id")
	groupId_int32 := int32(groupId.(int))

	log.Printf("[DEBUG] Resource Load Balancer. UPDATE. Checking attributes update state..")
	updateLBCmd := odk.SetLoadBalancerCommand{}

	oldSsl, newSsl := d.GetChange("ssl_enabled")
	if d.HasChange("ssl_enabled") {
		updateLBCmd.SslEnabled = newSsl.(bool)
	} else {
		updateLBCmd.SslEnabled = oldSsl.(bool)
	}

	oldServiceType, newServiceType := d.GetChange("service_type_id")
	if d.HasChange("service_type_id") {
		updateLBCmd.ServiceType = int32(newServiceType.(int))
	} else {
		updateLBCmd.ServiceType = int32(oldServiceType.(int))
	}

	oldPort, newPort := d.GetChange("port_number")
	if d.HasChange("port_number") {
		updateLBCmd.PortNumber = int32(newPort.(int))
	} else {
		updateLBCmd.PortNumber = int32(oldPort.(int))
	}

	oldTargetPort, newTargetPort := d.GetChange("target_port_number")
	if d.HasChange("target_port_number") {
		updateLBCmd.TargetPortNumber = int32(newTargetPort.(int))
	} else {
		updateLBCmd.TargetPortNumber = int32(oldTargetPort.(int))
	}

	oldSessionPersType, newSessionPersType := d.GetChange("session_persistence_type_id")
	if d.HasChange("session_persistence_type_id") {
		updateLBCmd.SessionPersistenceType = int32(newSessionPersType.(int))
	} else {
		updateLBCmd.SessionPersistenceType = int32(oldSessionPersType.(int))
	}

	oldLBAlgorithm, newLBAlgorithm := d.GetChange("load_balancer_algorithm_id")
	if d.HasChange("load_balancer_algorithm_id") {
		updateLBCmd.LoadBalancerAlgorithm = int32(newLBAlgorithm.(int))
	} else {
		updateLBCmd.LoadBalancerAlgorithm = int32(oldLBAlgorithm.(int))
	}

	oldIPVersion, newIPVersion := d.GetChange("ip_version_id")
	if d.HasChange("ip_version_id") {
		updateLBCmd.IpVersion = int32(newIPVersion.(int))
	} else {
		updateLBCmd.IpVersion = int32(oldIPVersion.(int))
	}

	oldHealthCheck, newHealthCheck := d.GetChange("health_check_enabled")
	if d.HasChange("health_check_enabled") {
		updateLBCmd.HealthCheckEnabled = newHealthCheck.(bool)
	} else {
		updateLBCmd.HealthCheckEnabled = oldHealthCheck.(bool)
	}

	oldCmmnPers, newCmmnPers := d.GetChange("common_persistence_for_http_and_https_enabled")
	if d.HasChange("common_persistence_for_http_and_https_enabled") {
		updateLBCmd.CommonPersistenceForHttpAndHttpsEnabled = newCmmnPers.(bool)
	} else {
		updateLBCmd.CommonPersistenceForHttpAndHttpsEnabled = oldCmmnPers.(bool)
	}

	log.Printf("[DEBUG] Resource Load Balancer. UPDATE. Trying to post changes")
	_, resp, err := client.OCIGroupsApi.LoadBalancersUpdate(*auth, groupId_int32, updateLBCmd)

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Resource Load Balancer. UPDATE. Load balancer for group  with id %s was not found", strconv.Itoa(int(groupId_int32)))
		}
		return fmt.Errorf("Resource Load Balancer. UPDATE. Error occured while updating load balancer. %s", err)
	}
	return resourceLoadBalancerRead(d, m)
}

func resourceLoadBalancerDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource Load Balancer.DELETE. Initializing")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx
	groupId, _ := d.GetChange("group_id")
	groupId_int32 := int32(groupId.(int))

	log.Printf("[DEBUG] Resource Load Balancer. DELETE. Trying to delete LB")

	_, resp, err := client.OCIGroupsApi.LoadBalancersDelete(*auth, groupId_int32)
	if err != nil && !strings.Contains(err.Error(), "EOF") {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Resource Load Balancer. DELETE. Load Balancer for group with id %s was not found", strconv.Itoa(int(groupId_int32)))
		}
		return fmt.Errorf("Resource Load Balancer. DELETE. Error occured while deleting Load Balancer. %s", err)
	}
	d.SetId("")
	return nil
}
