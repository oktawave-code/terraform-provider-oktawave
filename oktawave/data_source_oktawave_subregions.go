package oktawave

import (
	"fmt"
	"math"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

func getSubregionDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"is_active": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	}

}

func dataSourceSubregions() *schema.Resource {
	name := "items"
	dataSourceSchema := makeDataSourceSchema(name, getSubregionDataSourceSchema)
	dataSourceReadFunction := makeDataSourceRead(name, dataSourceSchema, getSubregionsList, mapRawSubregionToDataSourceModel)
	return &schema.Resource{
		ReadContext: dataSourceReadFunction,
		Schema:      dataSourceSchema,
	}
}

func getSubregionsList(config *ClientConfig) ([]odk.Subregion, error) {
	client := config.odkClient
	auth := config.odkAuth

	params := map[string]interface{}{
		"pageSize": int32(math.MaxInt16),
		"orderBy":  "Id",
	}
	list, _, err := client.SubregionsApi.SubregionsGet(*auth, params)
	if err != nil {
		return nil, fmt.Errorf("get subregions request failed, caused by: %s", err)
	}

	return list.Items, nil
}

func mapRawSubregionToDataSourceModel(subregion odk.Subregion) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"id":        subregion.Id,
		"name":      subregion.Name,
		"is_active": subregion.IsActive,
	}
	return result, nil
}
