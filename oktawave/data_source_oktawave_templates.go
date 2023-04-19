package oktawave

import (
	"fmt"
	"math"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

func dataSourceTemplates() *schema.Resource {
	name := "items"
	dataSourceSchema := makeDataSourceSchema(name, getTemplateDataSourceSchema)
	dataSourceReadFunction := makeDataSourceRead(name, dataSourceSchema, getTemplatesList, mapRawTemplateToDataSourceModel)
	return &schema.Resource{
		ReadContext: dataSourceReadFunction,
		Schema:      dataSourceSchema,
	}
}

func getTemplatesList(config *ClientConfig) ([]odk.Template, error) {
	client := config.odkClient
	auth := config.odkAuth

	params := map[string]interface{}{
		"pageSize": int32(math.MaxInt16),
		// "orderBy":  "Id",
	}
	list, _, err := client.OCITemplatesApi.TemplatesGet(*auth, params)
	if err != nil {
		return nil, fmt.Errorf("get templates request failed, caused by: %s", err)
	}

	return list.Items, nil
}

func mapRawTemplateToDataSourceModel(template odk.Template) (map[string]interface{}, error) {
	var disks []map[string]interface{}
	for _, disk := range template.Disks {
		disks = append(disks, map[string]interface{}{
			"id":             disk.Id,
			"name":           disk.Name,
			"space_capacity": disk.SpaceCapacity,
			"tier_id":        disk.Tier.Id,
			"tier_label":     disk.Tier.Label,
			"creation_date":  disk.CreationDate.String(),
			"controller":     disk.Controller,
			"slot":           disk.Slot,
			"is_system_disk": disk.IsSystemDisk,
		})
	}

	var software []map[string]interface{}
	for _, s := range template.Software {
		software = append(software, map[string]interface{}{
			"id":   s.Id,
			"name": s.Name,
		})
	}

	result := map[string]interface{}{
		"id":                              template.Id,
		"name":                            template.Name,
		"description":                     template.Description,
		"version":                         template.Version,
		"creation_date":                   template.CreationDate.String(),
		"last_change_date":                template.LastChangeDate.String(),
		"creation_user_id":                template.CreationUser.Id,
		"default_instance_type_id":        template.DefaultInstanceType.Id,
		"default_instance_type_label":     template.DefaultInstanceType.Label,
		"minimum_instance_type_id":        template.MinimumInstanceType.Id,
		"minimum_instance_type_label":     template.MinimumInstanceType.Label,
		"ethernet_controllers_number":     template.EthernetControllersNumber,
		"ethernet_controllers_type_id":    template.EthernetControllersType.Id,
		"ethernet_controllers_type_label": template.EthernetControllersType.Label,
		"system_category_id":              template.SystemCategory.Id,
		"system_category_label":           template.SystemCategory.Label,
		"template_type_id":                template.TemplateType.Id,
		"template_type_label":             template.TemplateType.Label,
		"disks":                           disks,
		"software":                        software,
	}

	if template.OwnerAccount != nil {
		result["owner_account_id"] = template.OwnerAccount.Id
	}

	if template.PublicationStatus != nil {
		result["publication_status_id"] = template.PublicationStatus.Id
		result["publication_status_label"] = template.PublicationStatus.Label
	}
	return result, nil
}
