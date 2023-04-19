package oktawave

import (
	"fmt"
	"math"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

func dataSourceIps() *schema.Resource {
	name := "items"
	dataSourceSchema := makeDataSourceSchema(name, getIpDataSourceSchema)
	dataSourceReadFunction := makeDataSourceRead(name, dataSourceSchema, getIpsList, mapRawIpToDataSourceModel)
	return &schema.Resource{
		ReadContext: dataSourceReadFunction,
		Schema:      dataSourceSchema,
	}
}

func mapRawIpToDataSourceModel(ip odk.Ip) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"id":                   ip.Id,
		"subregion_id":         ip.Subregion.Id,
		"default_subregion_id": ip.DefaultSubregion.Id,
		"comment":              ip.Comment,
		"address":              ip.Address,
		"address_v6":           ip.AddressV6,
		"mac_address":          ip.MacAddress,
		"rev_dns":              ip.RevDns,
		"rev_dns_v6":           ip.RevDnsV6,
		"type_id":              ip.Type_.Id,
		"type_label":           ip.Type_.Label,
		"gateway":              ip.Gateway,
		"netmask":              ip.NetMask,
		"interface_id":         ip.InterfaceId,
		"dns_prefix":           ip.DnsPrefix,
		"dhcp_branch":          ip.DhcpBranch,
		"mode_id":              ip.Mode.Id,
		"mode_label":           ip.Mode.Label,
		"owner_account_id":     ip.OwnerAccount.Id,
		"creation_user_id":     ip.CreationUser.Id,
	}

	if ip.Instance != nil {
		result["instance_id"] = ip.Instance.Id
	} else {
		result["instance_id"] = nil
	}
	return result, nil
}

func getIpsList(config *ClientConfig) ([]odk.Ip, error) {
	client := config.odkClient
	auth := config.odkAuth

	params := map[string]interface{}{
		"pageSize": int32(math.MaxInt16),
		"orderBy":  "Id",
	}
	ips, _, err := client.OCIInterfacesApi.InstancesGetIps(*auth, params)
	if err != nil {
		return nil, fmt.Errorf("get ips request failed, caused by: %s", err)
	}

	return ips.Items, nil
}
