package oktawave

import (
	"fmt"
	"math"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

func dataSourceInstances() *schema.Resource {
	name := "instances"
	dataSourceSchema := makeDataSourceSchema(name, getInstanceDataSourceSchema)
	dataSourceReadFunction := makeDataSourceRead(name, dataSourceSchema, getInstancesList, mapRawInstancesToDataSourceModel)
	return &schema.Resource{
		ReadContext: dataSourceReadFunction,
		Schema:      dataSourceSchema,
	}
}

func getInstancesList(config *ClientConfig) ([]odk.Instance, error) {
	client := config.odkClient
	auth := config.odkAuth

	params := map[string]interface{}{
		"pageSize": int32(math.MaxInt16),
		"orderBy":  "Id",
	}
	list, _, err := client.OCIApi.InstancesGet(*auth, params)
	if err != nil {
		return nil, fmt.Errorf("get instances request failed, caused by: %s", err)
	}

	return list.Items, nil
}

func mapRawInstancesToDataSourceModel(instance odk.Instance) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"id":                         instance.Id,
		"name":                       instance.Name,
		"subregion_id":               instance.Subregion.Id,
		"template_id":                instance.Template.Id,
		"type_id":                    instance.Type_.Id,
		"type_label":                 instance.Type_.Label,
		"is_freemium":                instance.IsFreemium,
		"creation_date":              instance.CreationDate.String(),
		"creation_user_id":           instance.CreationUser.Id,
		"is_locked":                  instance.IsLocked,
		"locking_date":               instance.LockingDate.String(),
		"ip_address":                 instance.IpAddress,
		"private_ip_address":         instance.PrivateIpAddress,
		"dns_address":                instance.DnsAddress,
		"status_id":                  instance.Status.Id,
		"status_label":               instance.Status.Label,
		"system_category_id":         instance.SystemCategory.Id,
		"system_category_label":      instance.SystemCategory.Label,
		"autoscaling_type_id":        instance.AutoscalingType.Id,
		"autoscaling_type_label":     instance.AutoscalingType.Label,
		"vmware_tools_status_id":     instance.VmWareToolsStatus.Id,
		"vmware_tools_status_label":  instance.VmWareToolsStatus.Label,
		"monit_status_id":            instance.MonitStatus.Id,
		"monit_status_label":         instance.MonitStatus.Label,
		"template_type_id":           instance.Template.Id,
		"template_type_label":        instance.TemplateType.Label,
		"payment_type_id":            instance.PaymentType.Id,
		"payment_type_label":         instance.PaymentType.Label,
		"scsi_controller_type_id":    instance.ScsiControllerType.Id,
		"scsi_controller_type_label": instance.ScsiControllerType.Label,
		"total_disks_capacity":       instance.TotalDisksCapacity,
		"cpu_number":                 instance.CpuNumber,
		"ram_mb":                     instance.RamMb,
	}

	if instance.HealthCheck != nil {
		result["health_check_id"] = instance.HealthCheck.Id
	}

	if instance.SupportType != nil {
		result["support_type_id"] = instance.SupportType.Id
		result["support_type_name"] = instance.SupportType.Name
	}

	return result, nil
}
