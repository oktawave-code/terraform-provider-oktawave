package oktawave

import (
	"context"
	"math"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

func getOpnDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"creation_user_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"creation_date": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"last_change_date": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"private_ips": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"interface_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"mac_address": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"address": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"address_v6": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"instance_id": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"creation_date": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func dataSourceOpn() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOpnRead,
		Schema:      getOpnDataSourceSchema(),
	}
}

func dataSourceOpnRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "reading opn")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	id, ok := d.GetOk("id")
	if !ok {
		return diag.Errorf("Id must be specified.")
	}

	vars := map[string]interface{}{
		"pageSize": int32(math.MaxInt16),
	}
	opn, _, err := client.NetworkingApi.OpnsGet_1(*auth, int32(id.(int)), vars)
	if err != nil {
		return diag.Errorf("Opn with id %d not found. Caused by: %s", id.(int), err)
	}

	return loadDataSourceOpnToSchema(d, opn)
}

func loadDataSourceOpnToSchema(d *schema.ResourceData, opn odk.Opn) diag.Diagnostics {
	if err := d.Set("id", opn.Id); err != nil {
		return diag.Errorf("Error setting id: %s", err)
	}

	if err := d.Set("name", opn.Name); err != nil {
		return diag.Errorf("Error setting name: %s", err)
	}

	if err := d.Set("creation_user_id", opn.CreationUser.Id); err != nil {
		return diag.Errorf("Error setting creation_user_id: %s", err)
	}

	if err := d.Set("creation_date", opn.CreationDate.String()); err != nil {
		return diag.Errorf("Error setting creation_date: %s", err)
	}

	if err := d.Set("last_change_date", opn.LastChangeDate.String()); err != nil {
		return diag.Errorf("Error setting last_change_date: %s", err)
	}

	privateIps := make([]map[string]interface{}, len(opn.PrivateIps))
	for idx, ip := range opn.PrivateIps {
		privateIps[idx] = privateIpToMap(ip)
	}

	if err := d.Set("private_ips", privateIps); err != nil {
		return diag.Errorf("Error setting private_ips: %s", err)
	}

	d.SetId(strconv.Itoa(int(opn.Id)))
	return nil
}

func privateIpToMap(privateIp odk.PrivateIp) map[string]interface{} {
	return map[string]interface{}{
		"interface_id":  privateIp.InterfaceId,
		"mac_address":   privateIp.MacAddress,
		"address":       privateIp.Address,
		"address_v6":    privateIp.AddressV6,
		"instance_id":   privateIp.Instance.Id,
		"creation_date": privateIp.CreationDate.String(),
	}
}
