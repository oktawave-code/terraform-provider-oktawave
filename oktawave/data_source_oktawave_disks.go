package oktawave

import (
	"fmt"
	"math"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

func dataSourceDisks() *schema.Resource {
	name := "disks"
	dataSourceSchema := makeDataSourceSchema(name, getDiskDataSourceSchema)
	dataSourceReadFunction := makeDataSourceRead(name, dataSourceSchema, getDisksList, mapRawDiskToDataSourceModel)
	return &schema.Resource{
		ReadContext: dataSourceReadFunction,
		Schema:      dataSourceSchema,
	}
}

func getDisksList(config *ClientConfig) ([]odk.Disk, error) {
	client := config.odkClient
	auth := config.odkAuth

	params := map[string]interface{}{
		"pageSize": int32(math.MaxInt16),
		"orderBy":  "Id",
	}
	list, _, err := client.OVSApi.DisksGetDisks(*auth, params)
	if err != nil {
		return nil, fmt.Errorf("get disks request failed, caused by: %s", err)
	}

	return list.Items, nil
}

func mapRawDiskToDataSourceModel(disk odk.Disk) (map[string]interface{}, error) {
	connections := make([]map[string]interface{}, len(disk.Connections))
	for idx, connection := range disk.Connections {
		connections[idx] = connectionToMap(connection)
	}

	result := map[string]interface{}{
		"id":               disk.Id,
		"name":             disk.Name,
		"tier_id":          disk.Tier.Id,
		"tier_label":       disk.Tier.Label,
		"subregion_id":     disk.Subregion.Id,
		"capacity":         disk.SpaceCapacity,
		"creation_user_id": disk.CreationUser.Id,
		"creation_date":    disk.CreationDate.String(),
		"is_shared":        disk.IsShared,
		"is_locked":        disk.IsLocked,
		"locking_date":     disk.LockingDate.String(),
		"is_freemium":      disk.IsFreemium,
		"connections":      connections,
	}

	if disk.SharedDiskType != nil {
		result["shared_disk_type_id"] = disk.SharedDiskType.Id
		result["shared_disk_type_label"] = disk.SharedDiskType.Label
	}

	return result, nil
}
