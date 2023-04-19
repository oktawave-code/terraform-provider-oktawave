package oktawave

import (
	"fmt"
	"math"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

func dataSourceSshKeys() *schema.Resource {
	name := "items"
	dataSourceSchema := makeDataSourceSchema(name, getSshKeyDataSourceSchema)
	dataSourceReadFunction := makeDataSourceRead(name, dataSourceSchema, getSshKeysList, mapRawSshKeyToDataSourceModel)
	return &schema.Resource{
		ReadContext: dataSourceReadFunction,
		Schema:      dataSourceSchema,
	}
}

func getSshKeysList(config *ClientConfig) ([]odk.SshKey, error) {
	client := config.odkClient
	auth := config.odkAuth

	params := map[string]interface{}{
		"pageSize": int32(math.MaxInt16),
		"orderBy":  "Id",
	}
	list, _, err := client.AccountApi.AccountGetSshKeys(*auth, params)
	if err != nil {
		return nil, fmt.Errorf("get ssh keys request failed, caused by: %s", err)
	}

	return list.Items, nil
}

func mapRawSshKeyToDataSourceModel(key odk.SshKey) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"id":            key.Id,
		"name":          key.Name,
		"owner_user_id": key.OwnerUser.Id,
		"creation_date": key.CreationDate.String(),
	}
	return result, nil
}
