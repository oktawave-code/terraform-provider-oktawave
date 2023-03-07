package oktawave

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

func getIpDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"subregion_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "ID from subregions resource",
		},
		"default_subregion_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"comment": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"address": {
			Type:     schema.TypeString,
			Required: true,
		},
		"address_v6": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"mac_address": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"instance_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "ID from instances resource",
		},
		"rev_dns": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"rev_dns_v6": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"type_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Value from dictionary #123 (Static/Automatic/...)",
		},
		"type_label": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"gateway": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"netmask": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"interface_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"dns_prefix": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"dhcp_branch": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"mode_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Value from dictionary #301 (Normal/Floating/KAS)",
		},
		"mode_label": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"owner_account_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"creation_user_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func dataSourceIp() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIpAddressRead,
		Schema:      getIpDataSourceSchema(),
	}
}

func dataSourceIpAddressRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "reading ip address")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	address, ok := d.GetOk("address")
	if !ok {
		return diag.Errorf("Address must be specified.")
	}

	vars := map[string]interface{}{}
	ip, _, err := client.FloatingIPsApi.FloatingIpsGetIp(*auth, address.(string), vars)
	if err != nil {
		return diag.Errorf("Ip address %s not found. Caused by: %s", address.(string), err)
	}

	return loadDataSourceIpToSchema(d, ip)
}

func loadDataSourceIpToSchema(d *schema.ResourceData, ip odk.Ip) diag.Diagnostics {
	if err := d.Set("id", ip.Id); err != nil {
		return diag.Errorf("Error setting id: %s", err)
	}
	if err := d.Set("subregion_id", ip.Subregion.Id); err != nil {
		return diag.Errorf("Error setting subregion_id: %s", err)
	}
	if err := d.Set("default_subregion_id", ip.DefaultSubregion.Id); err != nil {
		return diag.Errorf("Error setting default_subregion_id: %s", err)
	}
	if err := d.Set("comment", ip.Comment); err != nil {
		return diag.Errorf("Error setting comment: %s", err)
	}
	if err := d.Set("address", ip.Address); err != nil {
		return diag.Errorf("Error setting address: %s", err)
	}
	if err := d.Set("address_v6", ip.AddressV6); err != nil {
		return diag.Errorf("Error setting address_v6: %s", err)
	}
	if err := d.Set("mac_address", ip.MacAddress); err != nil {
		return diag.Errorf("Error setting mac_address: %s", err)
	}
	if ip.Instance != nil {
		if err := d.Set("instance_id", ip.Instance.Id); err != nil {
			return diag.Errorf("Error setting instance_id: %s", err)
		}
	}
	if err := d.Set("rev_dns", ip.RevDns); err != nil {
		return diag.Errorf("Error setting rev_dns: %s", err)
	}
	if err := d.Set("rev_dns_v6", ip.RevDnsV6); err != nil {
		return diag.Errorf("Error setting rev_dns_v6: %s", err)
	}
	if err := d.Set("type_id", ip.Type_.Id); err != nil {
		return diag.Errorf("Error setting type_id: %s", err)
	}
	if err := d.Set("type_label", ip.Type_.Label); err != nil {
		return diag.Errorf("Error setting type_label: %s", err)
	}
	if err := d.Set("gateway", ip.Gateway); err != nil {
		return diag.Errorf("Error setting gateway: %s", err)
	}
	if err := d.Set("netmask", ip.NetMask); err != nil {
		return diag.Errorf("Error setting netmask: %s", err)
	}
	if err := d.Set("interface_id", ip.InterfaceId); err != nil {
		return diag.Errorf("Error setting interface_id: %s", err)
	}
	if err := d.Set("dns_prefix", ip.DnsPrefix); err != nil {
		return diag.Errorf("Error setting dns_prefix: %s", err)
	}
	if err := d.Set("dhcp_branch", ip.DhcpBranch); err != nil {
		return diag.Errorf("Error setting dhcp_branch: %s", err)
	}
	if err := d.Set("mode_id", ip.Mode.Id); err != nil {
		return diag.Errorf("Error setting mode_id: %s", err)
	}
	if err := d.Set("mode_label", ip.Mode.Label); err != nil {
		return diag.Errorf("Error setting mode_label: %s", err)
	}
	if err := d.Set("owner_account_id", ip.OwnerAccount.Id); err != nil {
		return diag.Errorf("Error setting owner_account_id: %s", err)
	}
	if err := d.Set("creation_user_id", ip.CreationUser.Id); err != nil {
		return diag.Errorf("Error setting creation_user_id: %s", err)
	}

	d.SetId(strconv.Itoa(int(ip.Id)))
	return nil
}
