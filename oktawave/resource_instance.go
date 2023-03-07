package oktawave

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

func resourceInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceCreate,
		ReadContext:   resourceInstanceRead,
		UpdateContext: resourceInstanceUpdate,
		DeleteContext: resourceInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Required
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of instance.",
			},
			"subregion_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID from subregions resource.",
			},
			"system_disk_class_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Defines disk performance class. Value from dictionary #17",
			},
			"template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Defines which image will be used for instance initialization. User can use standard image with one of popular operating systems or choose its own template. ID from templates resource",
			},
			"type_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Defines vCPU and RAM for this instance. Value from dictionary #12",
			},
			// Optional
			"authorization_method_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     DICT_LOGIN_TYPE_USER_AND_PASS,
				Description: "Two authorization methods are available - login/password or ssh-keys. Value from dictionary #159",
			},
			"opn_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "List of OPNs this instance is in.",
			},
			"system_disk_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "Disk size in GB. At least 5 GB. Disk capacity can be only scaled up.",
			},
			"without_public_ip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Allows to create instance without public IP. In this case this instance must be in at least one OPN.",
			},
			"init_script": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Must be base64 encoded. This script will be invoked during instance initialization.",
			},
			"ssh_keys_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "List of ssh keys injected to this instance during initialization.",
			},
			"public_ips": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "List of public IPs attached to this instance.",
			},
			// Computed
			"system_disk_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Id of instance system disk.",
			},
			"opn_mac": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Computed: true,
					Type:     schema.TypeString,
				},
				Computed:    true,
				Description: "MAC address of OPN.",
			},
			"creation_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when instance was created.",
			},
			"creation_user_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Id of user who created this resource.",
			},
			"is_locked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Tells if instance is locked.",
			},
			"locking_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tells when instance was locked.",
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Main ip address of this instance.",
			},
			"mac_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "MAC address of this instance.",
			},
			"private_ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private ip address",
			},
			"dns_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Auto generated dns address for this instance.",
			},
			"status_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Tells if instance is running or shut down, etc. Value from dictionary #27",
			},
			"system_category_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Value from dictionary #70",
			},
			"autoscaling_type_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Value from dictionary #56",
			},
			"vmware_tools_status_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Value from dictionary #155",
			},
			"monit_status_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Value from dictionary #92",
			},
			"template_type_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Value from dictionary #52",
			},
			"payment_type_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Value from dictionary #11",
			},
			"scsi_controller_type_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Value from dictionary #182",
			},
			"total_disks_capacity": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Capacity sum of all connected disks.",
			},
			"cpu_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Instance vCPU number.",
			},
			"ram_mb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Instance RAM size.",
			},
			"health_check_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Id of connected healthcheck.",
			},
			"support_type_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Support type id.",
			},
			"disks_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Ids of connected disks.",
			},
			"converted_to_template_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "Id of the template this instance was converted to. When instance is converted to template it ceases to exist and this attribute is set. After this, instance state will not be synchronized to prevent instance recreation. Instance definition may be safely removed from definition and state.",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},
		Description: "Oktawave Cloud Instance(OCI) is a virtual machine, base building block of cloud computing world.",
	}
}

func resourceInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "creating instance")

	err := validateInstanceResource(d)
	if err != nil {
		return diag.FromErr(err)
	}

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	authorizationMethod := (int32)(d.Get("authorization_method_id").(int))
	var initialIpAddress int32 = 0
	var ipAddressToAttach []int32
	if ipIdSet, ipIsSet := d.GetOk("public_ips"); ipIsSet {
		publicIps := castToInt32(ipIdSet.(*schema.Set).List())
		initialIpAddress = publicIps[0]
		ipAddressToAttach = publicIps[1:]
	}

	createCommand := odk.CreateInstanceCommand{
		InstanceName:          d.Get("name").(string),
		SubregionId:           (int32)(d.Get("subregion_id").(int)),
		DiskClass:             (int32)(d.Get("system_disk_class_id").(int)),
		TemplateId:            (int32)(d.Get("template_id").(int)),
		TypeId:                (int32)(d.Get("type_id").(int)),
		AuthorizationMethodId: authorizationMethod,
		OpnsIds:               castToInt32(d.Get("opn_ids").(*schema.Set).List()),
		DiskSize:              (int32)(d.Get("system_disk_size").(int)),
		InstancesCount:        1,
		WithoutPublicIp:       (bool)(d.Get("without_public_ip").(bool)),
		IPAddressId:           initialIpAddress,
		InitScript:            d.Get("init_script").(string),
	}

	if authorizationMethod == DICT_LOGIN_TYPE_SSH_KEYS {
		sshKeyIds := castToInt32(d.Get("ssh_keys_ids").(*schema.Set).List())
		if len(sshKeyIds) == 0 {
			return diag.Errorf("Empty ssh keys list used with authorization method == ssh keys")
		}
		tflog.Debug(ctx, "calling ODK AccountApi.AccountGetSshKeys")
		params := map[string]interface{}{
			"pageSize": int32(math.MaxInt16),
		}
		sshKeys, _, err := client.AccountApi.AccountGetSshKeys(*auth, params)
		if err != nil {
			return diag.Errorf("ODK Error in AccountApi.AccountGetSshKeys. %s", err)
		}
		if err := checkSshKeysList(sshKeyIds, sshKeys.Items); err != nil {
			return diag.Errorf("SSH keys problem. %s", err)
		}
		createCommand.SshKeysIds = sshKeyIds
	}

	tflog.Debug(ctx, "calling ODK OCIApi.InstancesPost", map[string]interface{}{"authorizationMethod": authorizationMethod, "name": d.Get("name").(string)})
	ticket, _, err := client.OCIApi.InstancesPost(*auth, createCommand)
	if err != nil {
		return diag.Errorf("ODK Error in OCIApi.InstancesPost. %s", err)
	}

	createTicket, err := waitForTicket(client, auth, ticket)
	if err != nil {
		return diag.Errorf("ODK Error in TicketsApi.TicketsGet. %s", err)
	}
	if createTicket.Status.Id != DICT_TICKET_SUCCEED {
		return diag.Errorf("Unable to create instance. Ticket status=%v", createTicket.Status.Id)
	}

	tflog.Info(ctx, fmt.Sprintf("successfully created OCI. id=%v", createTicket.ObjectId))
	d.SetId(strconv.Itoa(int(createTicket.ObjectId)))

	if len(ipAddressToAttach) > 0 {
		err := attachInstanceToIps(client, auth, ipAddressToAttach, createTicket.ObjectId)
		if err != nil {
			return diag.Errorf("Attaching IPs failed. %s", err)
		}
	}

	if disksIdSet, disksIsSet := d.GetOk("disks_ids"); disksIsSet {
		disksIds := castToInt32(disksIdSet.(*schema.Set).List())
		for _, diskId := range disksIds {
			err := attachDiskToInstance(client, auth, int32(diskId), createTicket.ObjectId)
			if err != nil {
				return diag.Errorf("Attaching disks failed. %s", err)
			}
		}
	}

	return resourceInstanceRead(ctx, d, m)
}

func resourceInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "reading instance")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	instanceId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Invalid OCI id: %v %s", d.Id(), err)
	}

	// Check if instance was converted to template. When instance is converted to template
	// it is no longer accessible on instances list. This is workaround to prevent instance recreation.
	// In this situation if template is found for instance, any data will not be further synchronized.
	// get template id
	isConverted, templateId, err := isConvertedToTemplate(client, auth, d, int32(instanceId))
	if err != nil {
		return diag.FromErr(err)
	}

	if isConverted {
		// get template id and set it, then print info and return
		if d.Set("converted_to_template_id", templateId) != nil {
			return diag.Errorf("Can't retrieve templateId")
		}
		tflog.Warn(ctx, "Instance was converted to template, definition may be safely removed from definition and state as it no longer exists.")
		return nil
	}
	// End of template workaround

	tflog.Debug(ctx, "calling ODK OCIApi.InstancesGet", map[string]interface{}{"id": instanceId})
	instance, resp, err := client.OCIApi.InstancesGet_2(*auth, (int32)(instanceId), nil)
	if err != nil {
		if resp != nil && (resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusForbidden) { // api returns 403 on missing instance
			d.SetId("")
		}
		return diag.Errorf("Error while retrieving OCI %v (possibly \"not found\" error): %s", instanceId, err)
	}

	return loadInstanceData(ctx, d, m, instance)
}

func resourceInstanceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "updating instance")

	err := validateInstanceResource(d)
	if err != nil {
		return diag.FromErr(err)
	}

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	instanceId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Invalid OCI id: %v %s", d.Id(), err)
	}

	// Converted to template workaround. Instance no longer exists
	isConverted, _, err := isConvertedToTemplate(client, auth, d, int32(instanceId))
	if err != nil {
		return diag.FromErr(err)
	}

	if isConverted {
		tflog.Warn(ctx, "Instance was converted to template, definition may be safely removed from definition and state as it no longer exists. No changes will be applied.")
		return nil
	}
	// End of workaround

	d.Partial(true)

	if d.HasChange("name") {
		tflog.Info(ctx, "instance name change detected")
		newName := d.Get("name").(string)
		tflog.Debug(ctx, "calling ODK OCIApi.InstancesChangeName")
		updateTicket, resp, err := client.OCIApi.InstancesChangeName(*auth, (int32)(instanceId), newName)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return diag.Errorf("OCI %v not found", instanceId)
			}
			return diag.Errorf("Error while updating OCI %v: %s", instanceId, err)
		}

		updateTicket, err = waitForTicket(client, auth, updateTicket)
		if err != nil {
			return diag.Errorf("ODK Error in TicketsApi.TicketsGet. %s", err)
		}
		if updateTicket.Status.Id != DICT_TICKET_SUCCEED {
			return diag.Errorf("Unable to update instance. Ticket status=%v", updateTicket.Status.Id)
		}
	}

	if d.HasChange("subregion_id") {
		tflog.Info(ctx, "subregion id change detected")
		//	newSubregionId := d.Get("subregion_id").(int)
		//	subregionUpdateCommand := odk.ChangeInstanceSubregionCommand{(int32)(newSubregionId)}
		//	updateTicket, _, err:=client.OCIApi.InstancesChangeSubregion(*auth, (int32)(instanceId), subregionUpdateCommand)
		//	... TODO (eventually)
		return diag.Errorf("Subregion changes are not supported for now")
	}

	if d.HasChange("type_id") {
		tflog.Info(ctx, "type id change detected")
		newTypeId := d.Get("type_id").(int)
		tflog.Debug(ctx, "calling ODK OCIApi.InstancesChangeType")
		updateTicket, _, err := client.OCIApi.InstancesChangeType_1(*auth, (int32)(instanceId), (int32)(newTypeId))
		if err != nil {
			return diag.Errorf("Error while updating OCI %v: %s", instanceId, err)
		}

		updateTicket, err = waitForTicket(client, auth, updateTicket)
		if err != nil {
			return diag.Errorf("ODK Error in TicketsApi.TicketsGet. %s", err)
		}
		if updateTicket.Status.Id != DICT_TICKET_SUCCEED {
			return diag.Errorf("Unable to update instance. Ticket status=%v", updateTicket.Status.Id)
		}
	}

	if d.HasChange("system_disk_size") || d.HasChange("system_disk_class_id") {
		tflog.Info(ctx, "system disk size or class change detected")
		systemDiskId := d.Get("system_disk_id").(int)
		tflog.Debug(ctx, "calling ODK OVSApi.DisksGet")
		disk, resp, err := client.OVSApi.DisksGet(*auth, int32(systemDiskId), nil)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return diag.Errorf("Disk %v not found", systemDiskId)
			}
			return diag.Errorf("ODK Error in OVSApi.DisksGet. %s", err)
		}

		updCmd := odk.UpdateDiskCommand{
			DiskName:        disk.Name,
			SpaceCapacity:   disk.SpaceCapacity,
			TierId:          disk.Tier.Id,
			SubregionId:     disk.Subregion.Id,
			InstanceIdsList: castIntToInt32(getConnectionInstanceIds(disk.Connections)),
		}
		if d.HasChange("system_disk_size") {
			_, newDiskSize := d.GetChange("system_disk_size")
			updCmd.SpaceCapacity = int32(newDiskSize.(int))
		}

		if d.HasChange("system_disk_class_id") {
			_, newDiskClass := d.GetChange("system_disk_class_id")
			updCmd.TierId = int32(newDiskClass.(int))
		}

		tflog.Debug(ctx, "calling ODK OVSApi.DisksPut")
		ticket, resp, err := client.OVSApi.DisksPut(*auth, disk.Id, updCmd)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return diag.Errorf("Disk %v not found", disk.Id)
			}
			return diag.Errorf("ODK Error in OVSApi.DisksPut. %s", err)
		}
		respTicket, err := waitForTicket(client, auth, ticket)
		if err != nil {
			return diag.Errorf("ODK Error in TicketsApi.TicketsGet. %s", err)
		}
		if respTicket.Status.Id != DICT_TICKET_SUCCEED {
			return diag.Errorf("Unable to update disk. Ticket status=%v", respTicket.Status.Id)
		}
	}

	var opnsToAttach []int32
	var opnsToDetach []int32
	if d.HasChange("opn_ids") {
		tflog.Info(ctx, "opn ids change detected")
		oldOpnId, newOpnId := d.GetChange("opn_ids")
		oldOpnIds := oldOpnId.(*schema.Set).List()
		newOpnIds := newOpnId.(*schema.Set).List()

		oldOpnIds_int32 := castToInt32(oldOpnIds)
		newOpnIds_int32 := castToInt32(newOpnIds)

		opnsToDetach = calcListAMinusListB(oldOpnIds_int32, newOpnIds_int32)
		opnsToAttach = calcListAMinusListB(newOpnIds_int32, oldOpnIds_int32)
	}

	var ipsToAttach []int32
	var ipsToDetach []int32
	if d.HasChange("public_ips") {
		tflog.Info(ctx, "ip change detected")
		oldIpId, newIpId := d.GetChange("public_ips")
		oldIps := oldIpId.(*schema.Set).List()
		newIps := newIpId.(*schema.Set).List()

		oldIpsList := castToInt32(oldIps)
		newIpsList := castToInt32(newIps)

		ipsToDetach = calcListAMinusListB(oldIpsList, newIpsList)
		ipsToAttach = calcListAMinusListB(newIpsList, oldIpsList)
	}

	if len(opnsToAttach) > 0 {
		if err := attachInstanceToOpns(client, auth, opnsToAttach, int32(instanceId)); err != nil {
			return diag.Errorf("Attaching OPNs failed. %s", err)
		}
	}

	if len(ipsToAttach) > 0 {
		if err := attachInstanceToIps(client, auth, ipsToAttach, int32(instanceId)); err != nil {
			return diag.Errorf("Attaching IPs failed. %s", err)
		}
	}

	if len(opnsToDetach) > 0 {
		if err := detachInstanceFromOpns(client, auth, opnsToDetach, int32(instanceId)); err != nil {
			return diag.Errorf("Detaching OPNs with ids %v failed. Caused by: %s", opnsToDetach, err)
		}
	}

	if len(ipsToDetach) > 0 {
		if err := detachInstanceFromIps(client, auth, ipsToDetach, int32(instanceId)); err != nil {
			return diag.Errorf("Detaching IPs with ids %v failed. Caused by: %s", ipsToDetach, err)
		}
	}

	if d.HasChange("disks_ids") {
		oldDisksId, newDisksId := d.GetChange("disks_ids")
		oldDisks := oldDisksId.(*schema.Set).List()
		newDisks := newDisksId.(*schema.Set).List()

		oldDisksList := castToInt32(oldDisks)
		newDisksList := castToInt32(newDisks)

		disksIdListToDetach := calcListAMinusListB(oldDisksList, newDisksList)
		disksIdListToAttach := calcListAMinusListB(newDisksList, oldDisksList)

		for _, diskId := range disksIdListToAttach {
			err := attachDiskToInstance(client, auth, int32(diskId), int32(instanceId))
			if err != nil {
				return diag.Errorf("Attaching disks failed. %s", err)
			}
		}

		for _, diskId := range disksIdListToDetach {
			err := detachDiskFromInstance(client, auth, int32(diskId), int32(instanceId))
			if err != nil {
				return diag.Errorf("Detaching disks failed. %s", err)
			}
		}
	}

	d.Partial(false)
	return resourceInstanceRead(ctx, d, m)
}

func resourceInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "deleting instance")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	instanceId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Invalid OCI id: %v %s", d.Id(), err)
	}

	// Converted to template workaround. Instance no longer exists
	isConverted, _, err := isConvertedToTemplate(client, auth, d, int32(instanceId))
	if err != nil {
		return diag.FromErr(err)
	}

	if isConverted {
		// instance was converted to template and no longer exists, api should not be called
		d.SetId("")
		return nil
	}
	// End of workaround

	tflog.Debug(ctx, "calling ODK OCIApi.InstancesDelete")
	deleteTicket, _, err := client.OCIApi.InstancesDelete(*auth, (int32)(instanceId), nil)
	if err != nil {
		return diag.Errorf("ODK Error in OCIApi.InstancesDelete. %s", err)
	}

	deleteTicket, err = waitForTicket(client, auth, deleteTicket)
	if err != nil {
		return diag.Errorf("ODK Error in TicketsApi.TicketsGet. %s", err)
	}
	if deleteTicket.Status.Id != DICT_TICKET_SUCCEED {
		return diag.Errorf("Unable to delete instance. Ticket status=%v", deleteTicket.Status.Id)
	}

	d.SetId("")
	return nil
}

func loadInstanceData(ctx context.Context, d *schema.ResourceData, m interface{}, instance odk.Instance) diag.Diagnostics {
	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	// Load system disk data
	tflog.Debug(ctx, "calling ODK OCIApi.InstancesGetDisks", map[string]interface{}{"id": instance.Id})
	params := map[string]interface{}{
		"pageSize": int32(math.MaxInt16),
	}
	list, resp, err := client.OCIApi.InstancesGetDisks(*auth, int32(instance.Id), params)
	var systemDisk odk.Disk
	var customDisks []int32
	for _, disk := range list.Items {
		for _, connection := range disk.Connections {
			if connection.Instance.Id == instance.Id {
				if connection.IsSystemDisk {
					systemDisk = disk
				} else {
					customDisks = append(customDisks, disk.Id)
				}
			}
		}
	}
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return diag.Errorf("System disk for OCI %v not found", instance.Id)
		}
		return diag.Errorf("Error while retrieving system disk for OCI %v: %s", instance.Id, err)
	}

	// Load networking data
	opnMacMap, opnIds, err := getOpnsData(client, *auth, int32(instance.Id))
	if err != nil {
		return diag.Errorf("failed to load OPNs. %s", err)
	}

	tflog.Debug(context.Background(), "calling ODK FloatingIPsApi.FloatingIpsGetIp")
	params2 := map[string]interface{}{
		"instanceId": instance.Id,
		"pageSize":   int32(math.MaxInt16),
	}
	ips, _, err := client.FloatingIPsApi.FloatingIpsGetIps(*auth, params2)
	if err != nil {
		return diag.Errorf("ODK Error in FloatingIPsApi.FloatingIpsGetIp. %s", err)
	}
	publicIps := make([]int32, 0)
	var ipMac *string = nil
	for i, ip := range ips.Items {
		publicIps = append(publicIps, ip.Id)
		if i == 0 {
			// first public ip is used as instance mac address
			ipMac = &ip.MacAddress
		}
	}

	// Load ssh keys
	tflog.Debug(context.Background(), "calling ODK OCIApi.InstancesGetSshKeys")
	params3 := map[string]interface{}{
		"pageSize": int32(math.MaxInt16),
	}
	keys, _, err := client.OCIApi.InstancesGetSshKeys(*auth, int32(instance.Id), params3)
	if err != nil {
		return diag.Errorf("ODK Error in OCIApi.InstancesGetSshKeys. %s", err)
	}
	keyIds := make([]int32, 0)
	for _, key := range keys.Items {
		keyIds = append(keyIds, key.Id)
	}

	// Load init script
	tflog.Debug(context.Background(), "calling ODK OCIApi.InstancesGetInstanceInitScript")
	initScript, resp, err := client.OCIApi.InstancesGetInstanceInitScript(*auth, int32(instance.Id), nil)
	if err != nil {
		if resp.StatusCode != http.StatusNotFound {
			return diag.Errorf("Error while retrieving init script for OCI %v: %s", instance.Id, err)
		}
	}

	// Store everything
	tflog.Debug(ctx, "Parsing returned data")
	if d.Set("name", instance.Name) != nil {
		return diag.Errorf("Can't retrieve instance name")
	}
	if d.Set("subregion_id", int(instance.Subregion.Id)) != nil {
		return diag.Errorf("Can't retrieve subregion id")
	}
	if d.Set("system_disk_class_id", systemDisk.Tier.Id) != nil {
		return diag.Errorf("Can't retrieve system disk tier id")
	}
	if d.Set("template_id", instance.Template.Id) != nil {
		return diag.Errorf("Can't retrieve template id")
	}
	if d.Set("type_id", instance.Type_.Id) != nil {
		return diag.Errorf("Can't retrieve type id")
	}
	// Not obtainable: authorization_method_id
	if d.Set("opn_ids", opnIds) != nil {
		return diag.Errorf("Can't retrieve opn connections")
	}
	if d.Set("system_disk_size", int(systemDisk.SpaceCapacity)) != nil {
		return diag.Errorf("Can't retrieve system disk size")
	}
	// Not obtainable: without_public_ip
	if d.Set("init_script", initScript) != nil {
		return diag.Errorf("Can't retrieve init script")
	}
	if d.Set("ssh_keys_ids", keyIds) != nil {
		return diag.Errorf("Can't retrieve ssh key ids")
	}
	if d.Set("public_ips", publicIps) != nil {
		return diag.Errorf("Can't retrieve ip addresses")
	}
	if d.Set("system_disk_id", int(systemDisk.Id)) != nil {
		return diag.Errorf("Can't retrieve system disk id")
	}
	if d.Set("opn_mac", opnMacMap) != nil {
		return diag.Errorf("Can't retrieve opn mac adresses")
	}
	if d.Set("creation_date", instance.CreationDate.String()) != nil {
		return diag.Errorf("Can't retrieve creation date")
	}
	if d.Set("creation_user_id", instance.CreationUser.Id) != nil {
		return diag.Errorf("Can't retrieve creation user id")
	}
	if d.Set("is_locked", instance.IsLocked) != nil {
		return diag.Errorf("Can't retrieve locked option")
	}
	if d.Set("locking_date", instance.LockingDate.String()) != nil {
		return diag.Errorf("Can't retrieve locking date")
	}
	if d.Set("ip_address", instance.IpAddress) != nil {
		return diag.Errorf("Can't retrieve ip address")
	}
	if d.Set("mac_address", ipMac) != nil {
		return diag.Errorf("Can't retrieve mac address")
	}
	if d.Set("private_ip_address", instance.PrivateIpAddress) != nil {
		return diag.Errorf("Can't retrieve private ip address")
	}
	if d.Set("dns_address", instance.DnsAddress) != nil {
		return diag.Errorf("Can't retrieve dns ip address")
	}
	if d.Set("status_id", instance.Status.Id) != nil {
		return diag.Errorf("Can't retrieve status")
	}
	if d.Set("system_category_id", instance.SystemCategory.Id) != nil {
		return diag.Errorf("Can't retrieve system category")
	}
	if d.Set("autoscaling_type_id", instance.AutoscalingType.Id) != nil {
		return diag.Errorf("Can't retrieve autoscaling type")
	}
	if d.Set("vmware_tools_status_id", instance.VmWareToolsStatus.Id) != nil {
		return diag.Errorf("Can't retrieve vmware tools status")
	}
	if d.Set("monit_status_id", instance.MonitStatus.Id) != nil {
		return diag.Errorf("Can't retrieve monit status")
	}
	if d.Set("template_type_id", instance.TemplateType.Id) != nil {
		return diag.Errorf("Can't retrieve template type")
	}
	if d.Set("payment_type_id", instance.PaymentType.Id) != nil {
		return diag.Errorf("Can't retrieve payment type")
	}
	if d.Set("scsi_controller_type_id", instance.ScsiControllerType.Id) != nil {
		return diag.Errorf("Can't retrieve scsi controller type")
	}
	if d.Set("total_disks_capacity", instance.TotalDisksCapacity) != nil {
		return diag.Errorf("Can't retrieve total disks capacity")
	}
	if d.Set("cpu_number", instance.CpuNumber) != nil {
		return diag.Errorf("Can't retrieve cpu number")
	}
	if d.Set("ram_mb", instance.RamMb) != nil {
		return diag.Errorf("Can't retrieve ram size")
	}
	if instance.HealthCheck != nil {
		if d.Set("health_check_id", instance.HealthCheck.Id) != nil {
			return diag.Errorf("Can't retrieve health check")
		}
	}
	if instance.SupportType != nil {
		if d.Set("support_type_id", instance.SupportType.Id) != nil {
			return diag.Errorf("Can't retrieve support type")
		}
	}
	if d.Set("disks_ids", customDisks) != nil {
		return diag.Errorf("Can't retrieve attached disks list")
	}
	return nil
}

func checkSshKeysList(sshKeyIds []int32, acceptedSshKeys []odk.SshKey) error {
	// Build accepted keys set
	check := make(map[int32]struct{})
	for _, key := range acceptedSshKeys {
		check[key.Id] = struct{}{}
	}
	// Verify every key
	for _, id := range sshKeyIds {
		if _, ok := check[id]; !ok {
			return fmt.Errorf("SSH key %v is not present on the accepted ssh keys list", id)
		}
	}
	// No problems found
	return nil
}

func attachDiskToInstance(client odk.APIClient, auth *context.Context, diskId int32, instanceId int32) error {
	ticket, resp, err := client.OVSApi.DisksAttachToInstance(*auth, diskId, instanceId)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("instance id %v or disk id %v was not found", instanceId, diskId)
		}
		return fmt.Errorf("ODK Error in OVSApi.DisksAttachToInstance. %v", err)
	}
	ticket, err = waitForTicket(client, auth, ticket)
	if err != nil {
		return fmt.Errorf("can't attach disk. %s", err)
	}
	if ticket.Status.Id != DICT_TICKET_SUCCEED {
		return fmt.Errorf("can't attach disk. Ticket status=%v", ticket.Status.Id)
	}
	return nil
}

func detachDiskFromInstance(client odk.APIClient, auth *context.Context, diskId int32, instanceId int32) error {
	ticket, resp, err := client.OVSApi.DisksDetachFromInstance(*auth, diskId, instanceId)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("instance id %v or disk id %v was not found", instanceId, diskId)
		}
		return fmt.Errorf("ODK Error in OVSApi.DisksDetachFromInstance. %v", err)
	}
	ticket, err = waitForTicket(client, auth, ticket)
	if err != nil {
		return fmt.Errorf("can't detach disk. %s", err)
	}
	if ticket.Status.Id != DICT_TICKET_SUCCEED {
		return fmt.Errorf("can't detach disk. Ticket status=%v", ticket.Status.Id)
	}
	return nil
}

func attachInstanceToIps(client odk.APIClient, auth *context.Context, ips []int32, instanceId int32) error {
	for _, ip := range ips {
		ticket, _, err := attachIpById(client, auth, instanceId, ip)
		if err != nil {
			return fmt.Errorf("can't attach IP with id: %d. Caused by %s", ip, err)
		}
		if ticket.Status.Id != DICT_TICKET_SUCCEED {
			return fmt.Errorf("can't attach IP with id %d. Ticket id: %d, Ticket status=%v", ip, ticket.Id, ticket.Status.Id)
		}
	}
	return nil
}

func detachInstanceFromIps(client odk.APIClient, auth *context.Context, ips []int32, instanceId int32) error {
	for _, ip := range ips {
		ticket, _, err := detachIpById(client, auth, instanceId, ip)
		if err != nil {
			return fmt.Errorf("can't detach IP with id: %d. Caused by %s", ip, err)
		}
		if ticket.Status.Id != DICT_TICKET_SUCCEED {
			return fmt.Errorf("can't detach IP with id %d. Ticket id: %d, Ticket status=%v", ip, ticket.Id, ticket.Status.Id)
		}
	}
	return nil
}

func getOpnsData(client odk.APIClient, auth context.Context, instanceId int32) (map[string]string, []int, error) {
	tflog.Debug(context.Background(), "calling ODK OCIInterfacesApi.InstancesGetOpns")
	params := map[string]interface{}{
		"pageSize": int32(math.MaxInt16),
	}
	opns, response, err := client.OCIInterfacesApi.InstancesGetOpns(auth, instanceId, params)
	if err != nil {
		if response != nil && response.StatusCode != http.StatusNotFound {
			return nil, nil, fmt.Errorf("instance id %v was not found", instanceId)
		}
		return nil, nil, fmt.Errorf("ODK Error in OCIInterfacesApi.InstancesGetOpns. %s", err)
	}
	opnMacMap := make(map[string]string)
	opnIds := make([]int, len(opns.Items))
	for opnIx, opn := range opns.Items {
		opnIds[opnIx] = int(opn.Id)
		for i, ip := range opn.PrivateIps {
			if ip.Instance.Id == instanceId {
				key := strconv.Itoa(int(opn.Id))
				opnMacMap[key] = opn.PrivateIps[i].MacAddress
				break
			}
		}
	}
	return opnMacMap, opnIds, nil
}

func attachInstanceToOpns(client odk.APIClient, auth *context.Context, opnIds []int32, instanceId int32) error {
	for _, opnId := range opnIds {
		attachOpnCmd := odk.AttachInstanceToOpnCommand{
			OpnId: opnId,
		}
		tflog.Debug(context.Background(), "calling ODK OCIInterfacesApi.InstancesAttachOpn")
		ticket, resp, err := client.OCIInterfacesApi.InstancesAttachOpn(*auth, instanceId, attachOpnCmd)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("instance id %v or opn id %v was not found", instanceId, opnId)
			}
			return fmt.Errorf("ODK Error in OCIInterfacesApi.InstancesAttachOpn. %s", err)
		}
		resultTicket, err := waitForTicket(client, auth, ticket)
		if err != nil {
			return err
		}
		if resultTicket.Status.Id != DICT_TICKET_SUCCEED {
			return fmt.Errorf("unable to attach instance to opn. Ticket status=%v", resultTicket.Status.Id)
		}
	}
	return nil
}

func detachInstanceFromOpns(client odk.APIClient, auth *context.Context, opnIds []int32, instanceId int32) error {
	for _, opnId := range opnIds {
		detachOpnCmd := odk.DetachInstanceFromOpnCommand{
			OpnId: opnId,
		}
		tflog.Debug(context.Background(), "calling ODK OCIInterfacesApi.InstancesDetachFromOpn")
		ticket, resp, err := client.OCIInterfacesApi.InstancesDetachFromOpn(*auth, instanceId, detachOpnCmd)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("instance id %v or opn id %v was not found", instanceId, opnId)
			}
			return fmt.Errorf("ODK Error in OCIInterfacesApi.InstancesDetachFromOpn. %s", err)
		}
		resultTicket, err := waitForTicket(client, auth, ticket)
		if err != nil {
			return err
		}
		if resultTicket.Status.Id != DICT_TICKET_SUCCEED {
			return fmt.Errorf("can't detach opn with id %d from instance %d. Ticket id: %d, Ticket status=%v", opnId, instanceId, ticket.Id, ticket.Status.Id)
		}
	}
	return nil
}

func validateInstanceResource(d *schema.ResourceData) error {
	publicIps := d.Get("public_ips").(*schema.Set).List()
	privateIps := d.Get("opn_ids").(*schema.Set).List()

	if len(publicIps) == 0 && len(privateIps) == 0 {
		return fmt.Errorf("instance must have at least one entry in public_ips or opn_ids specified")
	}

	return nil
}

func getInstanceTemplate(client odk.APIClient, auth *context.Context, instanceId int32) (template *odk.Template, err error) {
	// if template id is not set, call api for template id
	tmpl, resp, err := client.OCIApi.InstancesGetTemplateByBaseVirtualMachineId(*auth, instanceId, nil)
	if err != nil {
		// if api returns error 404 - no template found, ok
		if resp.StatusCode == 404 {
			return nil, nil
		}
		// if api returns error other than 404 - print error and fail
		if resp.StatusCode != 404 {
			return nil, fmt.Errorf("failed to check template data for instance %d. Caused by: %s", instanceId, err)
		}
	}
	return &tmpl, nil
}

func isConvertedToTemplate(client odk.APIClient, auth *context.Context, d *schema.ResourceData, instanceId int32) (isConverted bool, templateId int32, err error) {
	tmplId, isTemplateIdSet := d.GetOk("converted_to_template_id")
	if isTemplateIdSet {
		return true, int32(tmplId.(int)), nil
	}

	template, err := getInstanceTemplate(client, auth, int32(instanceId))
	if err != nil {
		return false, 0, err
	}

	if template == nil {
		return false, 0, nil
	}

	return true, template.Id, nil
}
