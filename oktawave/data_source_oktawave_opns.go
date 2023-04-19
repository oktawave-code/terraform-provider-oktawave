package oktawave

import (
	"fmt"
	"math"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

func dataSourceOpns() *schema.Resource {
	name := "items"
	dataSourceSchema := makeDataSourceSchema(name, getOpnDataSourceSchema)
	dataSourceReadFunction := makeDataSourceRead(name, dataSourceSchema, getOpnsList, mapRawOpnToDataSourceModel)
	return &schema.Resource{
		ReadContext: dataSourceReadFunction,
		Schema:      dataSourceSchema,
	}
}

func getOpnsList(config *ClientConfig) ([]odk.Opn, error) {
	client := config.odkClient
	auth := config.odkAuth

	params := map[string]interface{}{
		"pageSize": int32(math.MaxInt16),
		"orderBy":  "Id",
	}
	list, _, err := client.NetworkingApi.OpnsGet(*auth, params)
	if err != nil {
		return nil, fmt.Errorf("get opns request failed, caused by: %s", err)
	}

	return list.Items, nil
}

func mapRawOpnToDataSourceModel(opn odk.Opn) (map[string]interface{}, error) {
	privateIps := make([]map[string]interface{}, len(opn.PrivateIps))
	for idx, ip := range opn.PrivateIps {
		privateIps[idx] = privateIpToMap(ip)
	}

	result := map[string]interface{}{
		"id":               opn.Id,
		"name":             opn.Name,
		"creation_user_id": opn.CreationUser.Id,
		"creation_date":    opn.CreationDate.String(),
		"last_change_date": opn.LastChangeDate.String(),
		"private_ips":      privateIps,
	}
	return result, nil
}
