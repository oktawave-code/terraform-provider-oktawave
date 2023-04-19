package oktawave

import (
	"fmt"
	"math"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

func dataSourceGroups() *schema.Resource {
	name := "items"
	dataSourceSchema := makeDataSourceSchema(name, getGroupDataSourceSchema)
	dataSourceReadFunction := makeDataSourceRead(name, dataSourceSchema, getGroupsList, mapRawGroupToDataSourceModel)
	return &schema.Resource{
		ReadContext: dataSourceReadFunction,
		Schema:      dataSourceSchema,
	}
}

func getGroupsList(config *ClientConfig) ([]odk.Group, error) {
	client := config.odkClient
	auth := config.odkAuth

	params := map[string]interface{}{
		"pageSize": int32(math.MaxInt16),
		"orderBy":  "Id",
	}
	list, _, err := client.OCIGroupsApi.GroupsGetGroups(*auth, params)
	if err != nil {
		return nil, fmt.Errorf("get groups request failed, caused by: %s", err)
	}

	return list.Items, nil
}

func mapRawGroupToDataSourceModel(group odk.Group) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"id":                       group.Id,
		"name":                     group.Name,
		"affinity_rule_type_id":    group.AffinityRuleType.Id,
		"affinity_rule_type_label": group.AffinityRuleType.Label,
		"is_load_balancer":         group.IsLoadBalancer,
		"instances_count":          group.InstancesCount,
		"schedulers_count":         group.SchedulersCount,
		"autoscaling_type_id":      group.AutoscalingType.Id,
		"autoscaling_type_label":   group.AutoscalingType.Label,
		"last_change_date":         group.LastChangeDate.String(),
		"creation_user_id":         group.CreationUser.Id,
	}
	return result, nil
}
