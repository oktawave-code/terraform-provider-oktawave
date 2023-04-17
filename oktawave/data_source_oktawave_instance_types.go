package oktawave

import (
	"fmt"
	"math"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

func getInstanceTypeDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"cpu": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"ram": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"category_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"category_label": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

}

func dataSourceInstanceTypes() *schema.Resource {
	name := "items"
	dataSourceSchema := makeDataSourceSchema(name, getInstanceTypeDataSourceSchema)
	dataSourceReadFunction := makeDataSourceRead(name, dataSourceSchema, getInstanceTypesList, mapRawInstanceTypeToDataSourceModel)
	return &schema.Resource{
		ReadContext: dataSourceReadFunction,
		Schema:      dataSourceSchema,
	}
}

func getInstanceTypesList(config *ClientConfig) ([]odk.InstanceType, error) {
	client := config.odkClient
	auth := config.odkAuth

	params := map[string]interface{}{
		"pageSize": int32(math.MaxInt16),
		// "orderBy":  "Id",
	}
	list, _, err := client.OCIApi.InstancesGetInstancesTypes(*auth, params)
	if err != nil {
		return nil, fmt.Errorf("get instance types request failed, caused by: %s", err)
	}

	return list.Items, nil
}

func mapRawInstanceTypeToDataSourceModel(instanceType odk.InstanceType) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"id":             instanceType.Id,
		"name":           instanceType.Name,
		"cpu":            instanceType.Cpu,
		"ram":            instanceType.Ram,
		"category_id":    instanceType.Category.Id,
		"category_label": instanceType.Category.Label,
	}
	return result, nil
}
