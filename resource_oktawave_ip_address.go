package main

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func resourceIpAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceIpAddressCreate,
		Read:   resourceIpAddressRead,
		Update: resourceIpAddressUpdate,
		Delete: resourceIpAddressDelete,
		Schema: map[string]*schema.Schema{
			"subregion_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"restore_rev_dns": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"restore_rev_dns_v6": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"rev_dns": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rev_dns_v6": {
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
			"gateway": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"netmask": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mac_address": {
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
			"type_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"creation_user_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},
	}
}

func resourceIpAddressCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource IpAddress. CREATE. Initializing")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx

	log.Printf("[DEBUG] Resource IpAddress. CREATE. Retrieving attributes from config file")
	subregion_id := d.Get("subregion_id").(int)
	comment, commentSet := d.GetOk("comment")
	//restRevDns := d.Get("restore_rev_dns")
	//restRevDnsV6 := d.Get("restore_rev_dns_v6")

	log.Printf("[DEBUG] Resource IpAddress. CREATE. Successfull retrieved attributes.Trying to post new ip")
	bookCommand := odk.BookIpCommand{
		SubregionId: int32(subregion_id),
	}
	ip, _, err := client.OCIInterfacesApi.InstancesBookNewIp(*auth, bookCommand)
	if err != nil {
		return fmt.Errorf("Resource IpAddress. CREATE. Cannot post new ip addres. %s", err)
	}

	log.Printf("[DEBUG]Resource IpAddress. CREATE. New ip was successfully created")
	d.SetId(strconv.Itoa(int(ip.Id)))

	log.Printf("[DEBUG]Resource IpAddress. CREATE. Trying to set additional parameters")
	if commentSet {
		updIp := odk.UpdateIpCommand{
			SetStatic: true,
			Comment:   comment.(string),
			//RestoreRevDns:   restRevDns.(bool),
			//RestoreRevDnsV6: restRevDnsV6.(bool),
		}
		_, resp, err := client.OCIInterfacesApi.InstancesUpdateIp(*auth, ip.Id, updIp)
		if err != nil && !strings.Contains(err.Error(), "EOF") {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				d.SetId("")
				return fmt.Errorf("Resource IpAddress. CREATE. Ip address was not found by id: %s", d.Id())
			}
			return fmt.Errorf("Resource IpAddress. CREATE. Error occured while setting additinoal parameters. %s", err)
		}
	}

	return resourceIpAddressRead(d, m)
}

func resourceIpAddressRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource IpAddress. READ. Initializing")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Resource IpAddress. READ. Invalid ip address id: %v", err)
	}

	log.Printf("[INFO] Resource IpAddress. READ. Initializing completed. Getting ip from remote server")
	ip, resp, err := client.OCIInterfacesApi.InstancesGetInstanceIp(*auth, int32(id), nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return fmt.Errorf("Resource IpAddress. READ. Ip address by id=%s was not found", strconv.Itoa(id))
		}
		return fmt.Errorf("Resource IpAddress. READ. Error retrieving ip: %s", err)
	}

	log.Printf("[DEBUG] Resource IpAddres. READ. Ip address was found. Synchronize local and remote state..")

	if err := d.Set("subregion_id", int(ip.Subregion.Id)); err != nil {
		return fmt.Errorf("Resource IpAddress. READ, error occured while retrieving ips remote subregion id. %s", err)
	}

	if err := d.Set("comment", ip.Comment); err != nil {
		return fmt.Errorf("Resource IpAddress. READ, error occured while retrieving ips remote comment. %s", err)
	}

	if err := d.Set("rev_dns", ip.RevDns); err != nil {
		return fmt.Errorf("Resource IpAddress. READ, error occured while retrieving ips remote reverse dns. %s", err)
	}

	if err := d.Set("rev_dns_v6", ip.RevDnsV6); err != nil {
		return fmt.Errorf("Resource IpAddress. READ, error occured while retrieving ips remote reverse dns v6. %s", err)
	}

	restoreRevDns := (ip.RevDns != "")
	if err := d.Set("restore_rev_dns", restoreRevDns); err != nil {
		return fmt.Errorf("Resource IpAddress. READ, error occured while retrieving ips remote reverse dns option. %s", err)
	}

	restoreRevDnsV6 := (ip.RevDnsV6 != "")
	if err := d.Set("restore_rev_dns_v6", restoreRevDnsV6); err != nil {
		return fmt.Errorf("Resource IpAddress. READ, error occured while retrieving ips remote reverse dns option V6. %s", err)
	}

	log.Printf("[DEBUG] Resource IpAddr. READ. Type id value is %s", strconv.Itoa(int(ip.Type_.Id)))
	if err := d.Set("type_id", int(ip.Type_.Id)); err != nil {
		return fmt.Errorf("Resource IpAddress. READ, error occured while retrieving ips remote type id. %s", err)
	}

	if err := d.Set("address", ip.Address); err != nil {
		return fmt.Errorf("Resource IpAddress. READ, error occured while retrieving ips remote address %s", err)
	}

	if err := d.Set("address_v6", ip.AddressV6); err != nil {
		return fmt.Errorf("Resource IpAddress. READ, error occured while retrieving ips remote address v6. %s", err)
	}

	if err := d.Set("gateway", ip.Gateway); err != nil {
		return fmt.Errorf("Resource IpAddress. READ, error occured while retrieving ips remote gateway. %s", err)
	}

	if err := d.Set("netmask", ip.NetMask); err != nil {
		return fmt.Errorf("Resource IpAddress. READ, error occured while retrieving ips remote netmask. %s", err)
	}

	if err := d.Set("mac_address", ip.MacAddress); err != nil {
		return fmt.Errorf("Resource IpAddress. READ, error occured while retrieving ips remote mac address. %s", err)
	}

	if err := d.Set("interface_id", ip.InterfaceId); err != nil {
		return fmt.Errorf("Resource IpAddress. READ, error occured while retrieving ips remote interface id. %s", err)
	}

	if err := d.Set("dns_prefix", ip.DnsPrefix); err != nil {
		return fmt.Errorf("Resource IpAddress. READ, error occured while retrieving ips remote dns prefix. %s", err)
	}

	if err := d.Set("creation_user_id", ip.CreationUser.Id); err != nil {
		return fmt.Errorf("Resource IpAddress. READ, error occured while retrieving ips remote creation user id. %s", err)
	}

	if err := d.Set("dhcp_branch", ip.DhcpBranch); err != nil {
		return fmt.Errorf("Resource IpAddress. READ, error occured while retrieving ips remote dhcp branch. %s", err)
	}

	log.Printf("[DEBUG] Resource IpAddress. READ. Remote and local state was synchronized successfully")
	return nil
}

func resourceIpAddressUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource IpAddress. UPDATE. Initializing")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx
	ipId, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Resource IpAddress. UPDATE. Invalid IpAddress id: %v", err)
	}

	updateIpCmd := odk.UpdateIpCommand{
		SetStatic: true,
		Comment:   d.Get("comment").(string),
	}

	if d.HasChange("instance_id") {
		oldInstanceId, newInstanceId := d.GetChange("instance_id")
		oldInstanceId_int32 := int32(oldInstanceId.(int))
		newInstanceId_int32 := int32(newInstanceId.(int))

		if oldInstanceId_int32 != 0 {
			log.Printf("[DEBUG] oldInstanceId = %s. New instance id %s", strconv.Itoa(oldInstanceId.(int)), strconv.Itoa(newInstanceId.(int)))
			ticket, resp, err := detachIpById(client, auth, oldInstanceId_int32, int32(ipId))
			if err != nil {
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					return fmt.Errorf("Resource IP id. UPDATE. %s", err)
				}
				return fmt.Errorf("Resource IpAddress. UPDATE. Error: %s", err)
			}
			switch ticket.Status.Id {
			//error case
			case TICKET_STATUS__ERROR:
				apiTicketStatusId := int(ticket.Status.Id)
				log.Printf("Resource IpAddress. UPDATE. Unable to detach ip. Ticket status: %s", strconv.Itoa(apiTicketStatusId))
			//succeed case
			case TICKET_STATUS__SUCCESS:
				log.Printf("[INFO] Resource IpAddress. UPDATE. Ip detaching. remote state was updated successfully")
			}
		}
		if newInstanceId_int32 != 0 {
			ticket, resp, err := attachIpById(client, auth, newInstanceId_int32, int32(ipId))
			if err != nil {
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					return fmt.Errorf("Resource IP id. UPDATE. %s", err)
				}
				return fmt.Errorf("Resource IpAddress. UPDATE. Error: %s", err)
			}
			switch ticket.Status.Id {
			//error case
			case TICKET_STATUS__ERROR:
				apiTicketStatusId := int(ticket.Status.Id)
				log.Printf("Resource IpAddress. UPDATE. Unable to attach ip. Ticket status: %s", strconv.Itoa(apiTicketStatusId))
			//succeed case
			case TICKET_STATUS__SUCCESS:
				log.Printf("[INFO] Resource IpAddress. UPDATE. Ip attach. remote state was updated successfully")
			}
		}
	}

	if d.HasChange("comment") {
		_, newComment := d.GetChange("comment")
		log.Printf("Found comment %s", newComment.(string))
		updateIpCmd.Comment = newComment.(string)
		log.Printf(updateIpCmd.Comment)
	}
	//if is not set -> can offer to change value. No guarantee to change in fact
	//if d.HasChange("restore_rev_dns") {
	//	_, newRestRevDns := d.GetChange("restore_rev_dns")
	//	updateIpCmd.RestoreRevDns = newRestRevDns.(bool)
	//}
	//
	////if is not set -> can offer to change value. No guarantee to change in fact
	//if d.HasChange("restore_rev_dns_v6") {
	//	_, newRestRev6Dns := d.GetChange("restore_rev_dns_v6")
	//	updateIpCmd.RestoreRevDnsV6 = newRestRev6Dns.(bool)
	//}

	_, resp, err := client.OCIInterfacesApi.InstancesUpdateIp(*auth, int32(ipId), updateIpCmd)

	if err != nil {
		if resp != nil && resp.StatusCode != 200 {
			if resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("Resource IpAddress. UPDATE. Ip address by id %s was not found", d.Id())
			}
		} else {
			if resp == nil {
				return fmt.Errorf("Resource IpAddress. UPDATE. Error occured while putting ip update. %s", err)
			}
		}
	}
	return resourceIpAddressRead(d, m)
}

func resourceIpAddressDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource IpAddress. DELETE. Initializing.")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Resource IpAddress. DELETE. Invalid ip id: %v", err)
	}

	log.Printf("[INFO]Resource IpAddress. DELETE. Trying to create detach ticket")
	instanceId, instanceIsSet := d.GetOk("instance_id")
	if instanceIsSet {
		instanceId_int32 := int32(instanceId.(int))
		ticket, resp, err := client.OCIInterfacesApi.InstancesPostDetachIpTicket(*auth, instanceId_int32, int32(id))
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("Resource IpAddress. DELETE. Instance by id %s or "+
					"ip address by id %s were not found to detach", strconv.Itoa(instanceId.(int)), d.Id())
			}
			return fmt.Errorf("Resource AddressIp. DELETE. Error occured while detaching ip %s", err)
		}

		detachTicket, err := evaluateTicket(client, auth, ticket)
		if err != nil {
			return fmt.Errorf("[INFO] Resource IpAddress. DELETE. Can't get ticket %s", err)
		}
		switch detachTicket.Status.Id {
		//error case
		case TICKET_STATUS__ERROR:
			apiTicketStatusId := int(detachTicket.Status.Id)
			log.Printf("Resource IpAddress. DELETE. Unable to detach ip. Ticket status: %s", strconv.Itoa(apiTicketStatusId))
		//succeed case
		case TICKET_STATUS__SUCCESS:
			log.Printf("[INFO] Resource IpAddress. DELETE. Ip detaching. remote state was updated successfully")
		}
	}
	_, resp, err := client.OCIInterfacesApi.InstancesDeleteIp(*auth, int32(id))
	if err != nil && !strings.Contains(err.Error(), "EOF") {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Resource IpAddress. DELETE. Ip was not found by id: %s", d.Id())
		}
		return fmt.Errorf("Resource IpAddress. DELETE. Error occured while deleting ip. %s", err)
	}
	d.SetId("")
	return nil
}
