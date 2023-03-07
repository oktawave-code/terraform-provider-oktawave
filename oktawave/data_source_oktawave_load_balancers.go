package oktawave

import (
	"fmt"
	"math"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

func dataSourceLoadBalancers() *schema.Resource {
	name := "load_balancers"
	dataSourceSchema := makeDataSourceSchema(name, getLoadBalancerDataSourceSchema)
	dataSourceReadFunction := makeDataSourceRead(name, dataSourceSchema, getLoadBalancersList, mapRawLoadBalancerToDataSourceModel)
	return &schema.Resource{
		ReadContext: dataSourceReadFunction,
		Schema:      dataSourceSchema,
	}
}

func getLoadBalancersList(config *ClientConfig) ([]odk.LoadBalancer, error) {
	client := config.odkClient
	auth := config.odkAuth

	params := map[string]interface{}{
		"pageSize": int32(math.MaxInt16),
		"orderBy":  "GroupId",
	}
	list, _, err := client.OCIGroupsApi.GroupsGetLoadBalancers(*auth, params)
	if err != nil {
		return nil, fmt.Errorf("get load balancers request failed, caused by: %s", err)
	}

	return list.Items, nil
}

func mapRawLoadBalancerToDataSourceModel(lb odk.LoadBalancer) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"group_name":                     lb.GroupName,
		"group_id":                       lb.GroupId,
		"service_type_id":                lb.ServiceType.Id,
		"service_type_label":             lb.ServiceType.Label,
		"port_number":                    lb.PortNumber,
		"target_port_number":             lb.TargetPortNumber,
		"ssl_target_port_number":         lb.SslTargetPortNumber,
		"session_persistence_type_id":    lb.SessionPersistenceType.Id,
		"session_persistence_type_label": lb.SessionPersistenceType.Label,
		"algorithm_id":                   lb.Algorithm.Id,
		"algorithm_label":                lb.Algorithm.Label,
		"ip_version_id":                  lb.IpVersion.Id,
		"ip_version_label":               lb.IpVersion.Label,
		"health_check_enabled":           lb.HealthCheckEnabled,
		"ssl_enabled":                    lb.SslEnabled,
		"common_persistence_enabled":     lb.CommonPersistenceForHttpAndHttpsEnabled,
		"ip_address":                     lb.IpAddress,
		"ip_address_v6":                  lb.IpV6Address,
	}

	if lb.ProxyProtocolVersion != nil {
		result["proxy_protocol_version_id"] = lb.ProxyProtocolVersion.Id
		result["proxy_protocol_version_label"] = lb.ProxyProtocolVersion.Label
	}
	return result, nil
}
