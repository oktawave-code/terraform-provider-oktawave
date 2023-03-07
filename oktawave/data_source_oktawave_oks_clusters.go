package oktawave

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/oktawave-code/oks-sdk"
)

func dataSourceOksClusters() *schema.Resource {
	name := "oks_clusters"
	dataSourceSchema := makeDataSourceSchema(name, getOksClusterDataSourceSchema)
	dataSourceReadFunction := makeDataSourceRead(name, dataSourceSchema, getOksClustersList, mapRawOksClusterToDataSourceModel)
	return &schema.Resource{
		ReadContext: dataSourceReadFunction,
		Schema:      dataSourceSchema,
	}
}

func getOksClustersList(config *ClientConfig) ([]swagger.K44SClusterDetailsDto, error) {
	client := config.oksClient
	auth := config.oksAuth

	list, _, err := client.ClustersApi.ClustersGet(*auth)
	if err != nil {
		return nil, fmt.Errorf("get oks clusters request failed, caused by: %s", err)
	}

	return list, nil
}

func mapRawOksClusterToDataSourceModel(cluster swagger.K44SClusterDetailsDto) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"name":          cluster.Name,
		"version":       cluster.Version,
		"creation_date": cluster.CreationDate.String(),
		"is_running":    cluster.Running,
	}
	return result, nil
}
