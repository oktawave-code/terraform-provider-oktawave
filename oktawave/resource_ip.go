package oktawave

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

func resourceIpAddress() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIpAddressCreate,
		ReadContext:   resourceIpAddressRead,
		UpdateContext: resourceIpAddressUpdate,
		DeleteContext: resourceIpAddressDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"subregion_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID from subregions resource",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Comment to this ip. Helps to quickly identify IP purpose.",
			},
			"rev_dns": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Reverse DNS for v4 IP.",
			},
			"rev_dns_v6": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Reverse DNS for v6 IP.",
			},
			"type_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Static or automatic.",
			},
			"address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IPv4 address",
			},
			"address_v6": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IPv6 address",
			},
			"gateway": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default gateway.",
			},
			"netmask": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP netmask",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of instance this IP is connected to.",
			},
			"mac_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Mac address",
			},
			"interface_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Interface id",
			},
			"dns_prefix": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNS prefix",
			},
			"dhcp_branch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DHCP branch",
			},
			"mode_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Value from dictionary #301 (Normal/Floating/KAS)",
			},
			"creation_user_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Id of user who created this ip.",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},
		Description: "IP addresses allows communication between devices over internet.",
	}
}

func resourceIpAddressCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "creating ip address")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	subregion_id := d.Get("subregion_id").(int)
	bookCommand := odk.BookIpCommand{ // only subregion can be defined here
		SubregionId: int32(subregion_id),
	}

	tflog.Debug(ctx, "calling ODK FloatingIPsApi.FloatingIpsBookNewIp")
	ip, _, err := client.OCIInterfacesApi.InstancesBookNewIp(*auth, bookCommand)
	if err != nil {
		return diag.Errorf("ODK Error in FloatingIPsApi.FloatingIpsBookNewIp. %s", err)
	}

	updateCommand := odk.UpdateIpCommand{}
	updateNeeded := false
	comment, isSet := d.GetOk("comment")
	if isSet {
		updateCommand.Comment = comment.(string)
		updateNeeded = true
	}
	revDns, isSet := d.GetOk("rev_dns")
	if isSet {
		updateCommand.RevDns = revDns.(string)
		updateNeeded = true
	}
	revDnsV6, isSet := d.GetOk("rev_dns_v6")
	if isSet {
		updateCommand.RevDnsV6 = revDnsV6.(string)
		updateNeeded = true
	}
	if updateNeeded {
		tflog.Debug(ctx, "calling ODK FloatingIPsApi.FloatingIpsUpdateIp")
		_, _, err := client.OCIInterfacesApi.InstancesUpdateIp(*auth, ip.Id, updateCommand)
		if err != nil && err.Error() != "EOF" { // "EOF" condition is a patch for ODK 1.4 bug: it reports error when API returns empty body
			return diag.Errorf("ODK Error in FloatingIPsApi.FloatingIpsUpdateIp. %s", err)
		}
	}

	tflog.Info(ctx, fmt.Sprintf("successfully created IP address. id=%v", ip.Address))
	d.SetId(strconv.Itoa(int(ip.Id)))

	return resourceIpAddressRead(ctx, d, m)
}

func resourceIpAddressRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "reading ip address")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Invalid Ip id: %v %s", d.Id(), err)
	}

	tflog.Debug(ctx, "calling ODK FloatingIPsApi.FloatingIpsGetIp")
	ip, resp, err := client.OCIInterfacesApi.InstancesGetInstanceIp(*auth, (int32)(id), nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
		}
		return diag.Errorf("ODK Error in FloatingIPsApi.FloatingIpsGetIp. %s", err)
	}

	return loadIpAddressData(ctx, d, m, ip)
}

func resourceIpAddressUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "updating ip address")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Invalid Ip id: %v %s", d.Id(), err)
	}

	updateCommand := odk.UpdateIpCommand{}
	updateNeeded := false
	if d.HasChange("comment") {
		updateCommand.Comment = d.Get("comment").(string)
		updateNeeded = true
	}
	if d.HasChange("rev_dns") {
		updateCommand.RevDns = d.Get("rev_dns").(string)
		updateNeeded = true
	}
	if d.HasChange("rev_dns_v6") {
		updateCommand.RevDnsV6 = d.Get("rev_dns_v6").(string)
		updateNeeded = true
	}
	if updateNeeded {
		tflog.Debug(ctx, "calling ODK FloatingIPsApi.FloatingIpsUpdateIp")
		_, _, err := client.OCIInterfacesApi.InstancesUpdateIp(*auth, (int32)(id), updateCommand)
		if err != nil && err.Error() != "EOF" { // "EOF" condition is a patch for ODK 1.4 bug: it reports error when API returns empty body
			return diag.Errorf("ODK Error in FloatingIPsApi.FloatingIpsUpdateIp. %s", err)
		}
	}

	if d.HasChange("subregion_id") {
		tflog.Debug(ctx, "calling ODK FloatingIPsApi.FloatingIpsChangeIpSubregionTicket")
		ticket, _, err := client.OCIInterfacesApi.InstancesChangeIpSubregionTicket(*auth, (int32)(id), map[string]interface{}{
			"subregionId": int32(d.Get("subregion_id").(int)),
		})
		if err != nil {
			return diag.Errorf("ODK Error in FloatingIPsApi.FloatingIpsChangeIpSubregionTicket. %s", err)
		}

		ticket, err = waitForTicket(client, auth, ticket)
		if err != nil {
			return diag.Errorf("ODK Error in TicketsApi.TicketsGet. %s", err)
		}
		if ticket.Status.Id != DICT_TICKET_SUCCEED {
			return diag.Errorf("Unable to change ip subregion. Ticket status=%v", ticket.Status.Id)
		}
	}

	return resourceIpAddressRead(ctx, d, m)
}

func resourceIpAddressDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "deleting ip address")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Invalid Ip id: %v %s", d.Id(), err)
	}

	ip, resp, err := client.OCIInterfacesApi.InstancesGetInstanceIp(*auth, (int32)(id), nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Failed to get IP data IP id: %d, Cause: %s", id, err)
	}

	if ip.Instance != nil {
		ticket, _, err := detachIpById(client, auth, ip.Instance.Id, (int32)(id))
		if err != nil {
			return diag.Errorf("Can't detach IP. %s", err)
		}
		if ticket.Status.Id != DICT_TICKET_SUCCEED {
			return diag.Errorf("Can't detach IP. Ticket status=%v", ticket.Status.Id)
		}
	}

	tflog.Debug(ctx, "calling ODK FloatingIPsApi.FloatingIpsDeleteIp")
	_, _, err = client.OCIInterfacesApi.InstancesDeleteIp(*auth, (int32)(id))
	if err != nil && err.Error() != "EOF" { // "EOF" condition is a patch for ODK 1.4 bug: it reports error when API returns empty body
		return diag.Errorf("ODK Error in FloatingIPsApi.FloatingIpsDeleteIp. %s", err)
	}

	d.SetId("")
	return nil
}

func loadIpAddressData(ctx context.Context, d *schema.ResourceData, m interface{}, ip odk.Ip) diag.Diagnostics {
	// Store everything
	tflog.Debug(ctx, "Parsing returned data")

	if d.Set("subregion_id", int(ip.Subregion.Id)) != nil {
		return diag.Errorf("Can't retrieve subregion id")
	}
	if d.Set("comment", ip.Comment) != nil {
		return diag.Errorf("Can't retrieve comment")
	}
	if d.Set("rev_dns", ip.RevDns) != nil {
		return diag.Errorf("Can't retrieve rev dns")
	}
	if d.Set("rev_dns_v6", ip.RevDnsV6) != nil {
		return diag.Errorf("Can't retrieve rev dns v6")
	}
	if ip.Type_ != nil {
		if d.Set("type_id", int(ip.Type_.Id)) != nil {
			return diag.Errorf("Can't retrieve type id")
		}
	}
	if d.Set("address", ip.Address) != nil {
		return diag.Errorf("Can't retrieve address")
	}
	if d.Set("address_v6", ip.AddressV6) != nil {
		return diag.Errorf("Can't retrieve address v6")
	}
	if d.Set("gateway", ip.Gateway) != nil {
		return diag.Errorf("Can't retrieve gateway")
	}
	if d.Set("netmask", ip.NetMask) != nil {
		return diag.Errorf("Can't retrieve netmask")
	}
	if ip.Instance != nil {
		if d.Set("instance_id", strconv.Itoa(int(ip.Instance.Id))) != nil {
			return diag.Errorf("Can't retrieve instance id")
		}
	} else {
		if d.Set("instance_id", "") != nil {
			return diag.Errorf("Can't retrieve instance id")
		}
	}
	if d.Set("mac_address", ip.MacAddress) != nil {
		return diag.Errorf("Can't retrieve mac address")
	}
	if d.Set("interface_id", ip.InterfaceId) != nil {
		return diag.Errorf("Can't retrieve interface id")
	}
	if d.Set("dns_prefix", ip.DnsPrefix) != nil {
		return diag.Errorf("Can't retrieve dns prefix")
	}
	if d.Set("dhcp_branch", ip.DhcpBranch) != nil {
		return diag.Errorf("Can't retrieve dhcp branch")
	}
	if ip.Mode != nil {
		if d.Set("mode_id", int(ip.Mode.Id)) != nil {
			return diag.Errorf("Can't retrieve mode id")
		}
	}
	if ip.CreationUser != nil {
		if d.Set("creation_user_id", ip.CreationUser.Id) != nil {
			return diag.Errorf("Can't retrieve creation user id")
		}
	}
	return nil
}
