package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"net/http"
)

func dataSourceOci() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceOciRead,
		Schema:	map[string]*schema.Schema {
			"instance_name": {
				Type:					schema.TypeString,
				Description:	"instance_name in https://api.oktawave.com/beta/docs/index#!/OCI/Instances_Get",
				Required:			true,
			},
			"subregion_id": {
				Type:					schema.TypeInt,
				Description:	"Subregion ID https://api.oktawave.com/beta/docs/index#!/Subregions/Subregions_Get",
				Optional:			true,
			},
			"template_id": {
				Type:					schema.TypeInt,
				Description:	"Template ID https://api.oktawave.com/beta/docs/index#!/OCI_Templates/Templates_Get",
				Optional:			true,
			},
			"type_id": {
				Type:					schema.TypeInt,
				Description:	"Instance type https://api.oktawave.com/beta/docs/index#!/OCI/Instances_GetInstancesTypes",
				Optional:			true,
			},
		},
	}
}

func dataSourceOciRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx
	instances, resp, err := client.OCIApi.InstancesGet(*auth, map[string]interface{}{
		"pageSize": int32(0),
	})
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Data Source OCI. READ. No resources found")
		}
		return fmt.Errorf("Data Source OCI. READ. Error retrieving instances: %s", err)
	}
	for _, instance := range instances.Items {
		if (instance.Name == d.Get("instance_name") &&
		((int32)(d.Get("subregion_id").(int)) == instance.Subregion.Id || d.Get("subregion_id") == 0) &&
		((int32)(d.Get("template_id").(int)) == instance.Template.Id || d.Get("template_id") == 0) &&
		((int32)(d.Get("type_id").(int)) == instance.Type_.Id || d.Get("type_id") == 0)) {
			d.SetId(fmt.Sprint(instance.Id))
			return nil
		}
	}
	return fmt.Errorf("Data Source OCI. READ. Instance not found: %s", d.Get("instance_name"))
}
