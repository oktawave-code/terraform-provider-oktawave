package oktawave

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

func resourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoadBalancerCreate,
		ReadContext:   resourceLoadBalancerRead,
		UpdateContext: resourceLoadBalancerUpdate,
		DeleteContext: resourceLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the group associated with this load balancer. Incoming traffic will be load balanced over instances in this group.",
			},
			"service_type_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     DICT_SERVICE_TYPE_HTTP,
				Description: "Type of load balancer service. Value from dictionary #15",
			},
			"port_number": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "LB external port number can be configured for Port ervice type",
			},
			"target_port_number": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Traffic will be sent to instances on this port.",
			},
			"ssl_target_port_number": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Ssl target port number.",
			},
			"session_persistence_type_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     DICT_SESSION_HANDLING_SOURCE_IP,
				Description: "Type of session persistence. Value from dictionary #16",
			},
			"algorithm_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     DICT_LB_ALGORITHM_ROUND_ROBIN,
				Description: "Load balancer algorithm. Value from dictionary #77",
			},
			"ip_version_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     DICT_IP_VERSION_IPV4,
				Description: "Type of IP. Value from dictionary #36 (v4 / v6 / both)",
			},
			"health_check_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Tells if healthcheck is enabled for this lb.",
			},
			"ssl_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Only for \"HTTP\" load balancer service type",
			},
			"common_persistence_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Only for \"HTTP\" load balancer service type",
			},
			"proxy_protocol_version_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     DICT_PROXY_PROTOCOL_NONE,
				Description: "Value from dictionary #302",
			},
			"group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of associated group.",
			},
			"address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IPv4 address of this load balancer",
			},
			"address_v6": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IPv6 address of this load balancer",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},
		Description: "Load balancers help in managing traffic load and increase reliability of services. Service based on physical hardware produced by Citrix.",
	}
}

func resourceLoadBalancerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "creating load balancer")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	createCmd := odk.SetLoadBalancerCommand{
		SslEnabled:                              d.Get("ssl_enabled").(bool),
		ServiceType:                             int32(d.Get("service_type_id").(int)),
		SessionPersistenceType:                  int32(d.Get("session_persistence_type_id").(int)),
		LoadBalancerAlgorithm:                   int32(d.Get("algorithm_id").(int)),
		IpVersion:                               int32(d.Get("ip_version_id").(int)),
		HealthCheckEnabled:                      d.Get("health_check_enabled").(bool),
		CommonPersistenceForHttpAndHttpsEnabled: d.Get("common_persistence_enabled").(bool),
		ProxyProtocolVersion:                    int32(d.Get("proxy_protocol_version_id").(int)),
	}
	port, isSet := d.GetOk("port_number")
	if isSet {
		createCmd.PortNumber = int32(port.(int))
	}
	port, isSet = d.GetOk("target_port_number")
	if isSet {
		createCmd.TargetPortNumber = int32(port.(int))
	}
	port, isSet = d.GetOk("ssl_target_port_number")
	if isSet {
		createCmd.SslTargetPortNumber = int32(port.(int))
	}

	tflog.Debug(ctx, "calling ODK OCIGroupsApi.LoadBalancersCreate")
	loadBalancer, _, err := client.OCIGroupsApi.LoadBalancersCreate(*auth, int32(d.Get("group_id").(int)), createCmd)
	if err != nil {
		return diag.Errorf("ODK Error in OCIGroupsApi.LoadBalancersCreate. %s", err)
	}

	tflog.Info(ctx, fmt.Sprintf("successfully created load balancer. id=%v", loadBalancer.GroupId))
	d.SetId(strconv.Itoa(int(loadBalancer.GroupId)))

	return resourceLoadBalancerRead(ctx, d, m)
}

func resourceLoadBalancerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "reading load balancer")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	groupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Invalid load balancer id: %v %s", d.Id(), err)
	}

	tflog.Debug(ctx, "calling ODK OCIGroupsApi.LoadBalancersGetLoadBalancer")
	loadBalancer, _, err := client.OCIGroupsApi.LoadBalancersGetLoadBalancer(*auth, int32(groupId), nil)
	if err != nil {
		return diag.Errorf("ODK Error in OCIGroupsApi.LoadBalancersGetLoadBalancer. %s", err)
	}

	return loadLoadBalancerData(ctx, d, m, loadBalancer)
}

func resourceLoadBalancerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "updating load balancer")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	groupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Invalid load balancer id: %v %s", d.Id(), err)
	}

	d.Partial(true)

	updateCmd := odk.SetLoadBalancerCommand{
		SslEnabled:                              d.Get("ssl_enabled").(bool),
		ServiceType:                             int32(d.Get("service_type_id").(int)),
		SessionPersistenceType:                  int32(d.Get("session_persistence_type_id").(int)),
		LoadBalancerAlgorithm:                   int32(d.Get("algorithm_id").(int)),
		IpVersion:                               int32(d.Get("ip_version_id").(int)),
		HealthCheckEnabled:                      d.Get("health_check_enabled").(bool),
		CommonPersistenceForHttpAndHttpsEnabled: d.Get("common_persistence_enabled").(bool),
		ProxyProtocolVersion:                    int32(d.Get("proxy_protocol_version_id").(int)),
	}
	port, isSet := d.GetOk("port_number")
	if isSet {
		updateCmd.PortNumber = int32(port.(int))
	}
	port, isSet = d.GetOk("target_port_number")
	if isSet {
		updateCmd.TargetPortNumber = int32(port.(int))
	}
	port, isSet = d.GetOk("ssl_target_port_number")
	if isSet {
		updateCmd.SslTargetPortNumber = int32(port.(int))
	}

	tflog.Debug(ctx, "calling ODK OCIGroupsApi.LoadBalancersUpdate")
	_, _, err = client.OCIGroupsApi.LoadBalancersUpdate(*auth, int32(groupId), updateCmd)
	if err != nil {
		return diag.Errorf("ODK Error in OCIGroupsApi.LoadBalancersUpdate. %s", err)
	}

	d.Partial(false)
	return resourceLoadBalancerRead(ctx, d, m)
}

func resourceLoadBalancerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "deleting load balancer")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	groupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Invalid load balancer id: %v %s", d.Id(), err)
	}

	tflog.Debug(ctx, "calling ODK OCIGroupsApi.LoadBalancersDelete")
	_, _, err = client.OCIGroupsApi.LoadBalancersDelete(*auth, int32(groupId))
	if err != nil && err.Error() != "EOF" { // "EOF" condition is a patch for ODK 1.4 bug: it reports error when API returns empty body
		return diag.Errorf("ODK Error in OCIGroupsApi.LoadBalancersDelete. %s", err)
	}

	d.SetId("")
	return nil
}

func loadLoadBalancerData(ctx context.Context, d *schema.ResourceData, m interface{}, loadBalancer odk.LoadBalancer) diag.Diagnostics {
	// Store everything
	tflog.Debug(ctx, "Parsing returned data")

	if d.Set("group_id", int(loadBalancer.GroupId)) != nil {
		return diag.Errorf("Can't retrieve group id")
	}
	if d.Set("service_type_id", int(loadBalancer.ServiceType.Id)) != nil {
		return diag.Errorf("Can't retrieve service type id")
	}
	if d.Set("port_number", int(loadBalancer.PortNumber)) != nil {
		return diag.Errorf("Can't retrieve port number")
	}
	if d.Set("target_port_number", int(loadBalancer.TargetPortNumber)) != nil {
		return diag.Errorf("Can't retrieve target port number")
	}
	if d.Set("ssl_target_port_number", int(loadBalancer.SslTargetPortNumber)) != nil {
		return diag.Errorf("Can't retrieve ssl target port number")
	}
	if d.Set("session_persistence_type_id", int(loadBalancer.SessionPersistenceType.Id)) != nil {
		return diag.Errorf("Can't retrieve session persistence type id")
	}
	if d.Set("algorithm_id", int(loadBalancer.Algorithm.Id)) != nil {
		return diag.Errorf("Can't retrieve algorithm id")
	}
	if d.Set("ip_version_id", int(loadBalancer.IpVersion.Id)) != nil {
		return diag.Errorf("Can't retrieve ip version id")
	}
	if d.Set("health_check_enabled", loadBalancer.HealthCheckEnabled) != nil {
		return diag.Errorf("Can't retrieve health check state")
	}
	if d.Set("ssl_enabled", loadBalancer.SslEnabled) != nil {
		return diag.Errorf("Can't retrieve ssl state")
	}
	if d.Set("common_persistence_enabled", loadBalancer.CommonPersistenceForHttpAndHttpsEnabled) != nil {
		return diag.Errorf("Can't retrieve common persistence state")
	}
	if d.Set("proxy_protocol_version_id", int(loadBalancer.ProxyProtocolVersion.Id)) != nil {
		return diag.Errorf("Can't retrieve proxy protocol version id")
	}
	if d.Set("group_name", loadBalancer.GroupName) != nil {
		return diag.Errorf("Can't retrieve group name")
	}
	if d.Set("address", loadBalancer.IpAddress) != nil {
		return diag.Errorf("Can't retrieve ip address")
	}
	if d.Set("address_v6", loadBalancer.IpV6Address) != nil {
		return diag.Errorf("Can't retrieve ip v6 address")
	}
	return nil
}
