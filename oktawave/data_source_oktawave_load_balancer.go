package oktawave

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

func getLoadBalancerDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"group_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"group_id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"service_type_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Value from dictionary #15",
		},
		"service_type_label": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"port_number": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"target_port_number": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"ssl_target_port_number": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"session_persistence_type_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Value from dictionary #16",
		},
		"session_persistence_type_label": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"algorithm_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Value from dictionary #77",
		},
		"algorithm_label": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ip_version_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Value from dictionary #36 (v4 / v6 / both)",
		},
		"ip_version_label": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"health_check_enabled": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"ssl_enabled": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Only for \"HTTP\" load balancer service type",
		},
		"common_persistence_enabled": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Only for \"HTTP\" load balancer service type",
		},
		"proxy_protocol_version_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Value from dictionary #302",
		},
		"proxy_protocol_version_label": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ip_address": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ip_address_v6": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLoadBalancerRead,
		Schema:      getLoadBalancerDataSourceSchema(),
	}
}

func dataSourceLoadBalancerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "reading load balancer")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	id, ok := d.GetOk("group_id")
	if !ok {
		return diag.Errorf("Id must be specified.")
	}

	vars := map[string]interface{}{}
	lb, _, err := client.OCIGroupsApi.LoadBalancersGetLoadBalancer(*auth, int32(id.(int)), vars)
	if err != nil {
		return diag.Errorf("Load Balancer for group with id %d not found. Caused by: %s", id.(int), err)
	}

	return loadDataSourceLoadBalancerToSchema(d, lb)
}

func loadDataSourceLoadBalancerToSchema(d *schema.ResourceData, lb odk.LoadBalancer) diag.Diagnostics {
	if err := d.Set("group_name", lb.GroupName); err != nil {
		return diag.Errorf("Error setting group_name: %s", err)
	}

	if err := d.Set("group_id", lb.GroupId); err != nil {
		return diag.Errorf("Error setting group_id: %s", err)
	}

	if err := d.Set("service_type_id", lb.ServiceType.Id); err != nil {
		return diag.Errorf("Error setting service_type_id: %s", err)
	}

	if err := d.Set("service_type_label", lb.ServiceType.Label); err != nil {
		return diag.Errorf("Error setting service_type_label: %s", err)
	}

	if err := d.Set("port_number", lb.PortNumber); err != nil {
		return diag.Errorf("Error setting port_number: %s", err)
	}

	if err := d.Set("target_port_number", lb.TargetPortNumber); err != nil {
		return diag.Errorf("Error setting target_port_number: %s", err)
	}

	if err := d.Set("ssl_target_port_number", lb.SslTargetPortNumber); err != nil {
		return diag.Errorf("Error setting ssl_target_port_number: %s", err)
	}

	if err := d.Set("session_persistence_type_id", lb.SessionPersistenceType.Id); err != nil {
		return diag.Errorf("Error setting session_persistence_type_id: %s", err)
	}

	if err := d.Set("session_persistence_type_label", lb.SessionPersistenceType.Label); err != nil {
		return diag.Errorf("Error setting session_persistence_type_label: %s", err)
	}

	if err := d.Set("algorithm_id", lb.Algorithm.Id); err != nil {
		return diag.Errorf("Error setting algorithm_id: %s", err)
	}

	if err := d.Set("algorithm_label", lb.Algorithm.Label); err != nil {
		return diag.Errorf("Error setting algorithm_label: %s", err)
	}

	if err := d.Set("ip_version_id", lb.IpVersion.Id); err != nil {
		return diag.Errorf("Error setting ip_version_id: %s", err)
	}

	if err := d.Set("ip_version_label", lb.IpVersion.Label); err != nil {
		return diag.Errorf("Error setting ip_version_label: %s", err)
	}

	if err := d.Set("health_check_enabled", lb.HealthCheckEnabled); err != nil {
		return diag.Errorf("Error setting health_check_enabled: %s", err)
	}

	if err := d.Set("ssl_enabled", lb.SslEnabled); err != nil {
		return diag.Errorf("Error setting ssl_enabled: %s", err)
	}

	if err := d.Set("common_persistence_enabled", lb.CommonPersistenceForHttpAndHttpsEnabled); err != nil {
		return diag.Errorf("Error setting common_persistence_enabled: %s", err)
	}

	if err := d.Set("proxy_protocol_version_id", lb.ProxyProtocolVersion.Id); err != nil {
		return diag.Errorf("Error setting proxy_protocol_version_id: %s", err)
	}

	if err := d.Set("proxy_protocol_version_label", lb.ProxyProtocolVersion.Label); err != nil {
		return diag.Errorf("Error setting proxy_protocol_version_label: %s", err)
	}

	if err := d.Set("ip_address", lb.IpAddress); err != nil {
		return diag.Errorf("Error setting ip_address: %s", err)
	}

	if err := d.Set("ip_address_v6", lb.IpV6Address); err != nil {
		return diag.Errorf("Error setting ip_address_v6: %s", err)
	}

	d.SetId(strconv.Itoa(int(lb.GroupId)))
	return nil
}
