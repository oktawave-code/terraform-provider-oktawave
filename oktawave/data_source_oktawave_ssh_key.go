package oktawave

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

func getSshKeyDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		// "value": {
		// 	Type:     schema.TypeString,
		// 	Computed: true,
		// },
		"owner_user_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"creation_date": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

}

func dataSourceSshKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSshKeyRead,
		Schema:      getSshKeyDataSourceSchema(),
	}
}

func dataSourceSshKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "reading ssh key")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	id, ok := d.GetOk("id")
	if !ok {
		return diag.Errorf("Id must be specified.")
	}

	vars := map[string]interface{}{}
	key, _, err := client.AccountApi.AccountGetSshKey(*auth, int32(id.(int)), vars)
	if err != nil {
		return diag.Errorf("Ssh key with id %d not found. Caused by: %s", id.(int), err)
	}

	return loadDataSourceSshKeyToSchema(d, key)
}

func loadDataSourceSshKeyToSchema(d *schema.ResourceData, key odk.SshKey) diag.Diagnostics {
	if err := d.Set("id", key.Id); err != nil {
		return diag.Errorf("Error setting id: %s", err)
	}

	if err := d.Set("name", key.Name); err != nil {
		return diag.Errorf("Error setting name: %s", err)
	}

	if err := d.Set("owner_user_id", key.OwnerUser.Id); err != nil {
		return diag.Errorf("Error setting owner_user_id: %s", err)
	}

	if err := d.Set("creation_date", key.CreationDate.String()); err != nil {
		return diag.Errorf("Error setting creation_date: %s", err)
	}

	d.SetId(strconv.Itoa(int(key.Id)))
	return nil
}
