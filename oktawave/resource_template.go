package oktawave

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

func resourceTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTemplateCreate,
		ReadContext:   resourceTemplateRead,
		UpdateContext: resourceTemplateUpdate,
		DeleteContext: resourceTemplateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": { // write only
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"system_category_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Value from dictionary #70",
			},
			"windows_type_id": { // write only
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Value from dictionary #84 (used only for windows category)",
			},
			"default_type_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Value from dictionary #12",
			},
			"minimum_type_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Value from dictionary #12",
			},
			"support_password": { // write only
				Type:     schema.TypeString,
				Optional: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_change_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_user_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ethernet_controllers_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ethernet_controllers_type_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Value from dictionary #167",
			},
			"publication_status_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Value from dictionary #140",
			},
			"disks": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Optional: true,
					Type:     schema.TypeString,
				},
				Computed:    true,
				Description: "Map id => name",
			},
			"software": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Optional: true,
					Type:     schema.TypeString,
				},
				Computed:    true,
				Description: "Map id => name",
			},
			"template_type_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Value from dictionary #52",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},
	}
}

func resourceTemplateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "creating template")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	descriptions := make([]odk.TemplateDescription, 0)
	descriptions = append(descriptions, odk.TemplateDescription{
		LanguageId:  DICT_LANGUAGE_PL,
		Description: d.Get("description").(string),
	})
	descriptions = append(descriptions, odk.TemplateDescription{
		LanguageId:  DICT_LANGUAGE_EN,
		Description: d.Get("description").(string),
	})

	createCommand := odk.ConvertInstanceToTemplateCommand{
		TemplateName:             d.Get("name").(string),
		TemplateDescriptions:     descriptions,
		TemplateVersion:          d.Get("version").(string),
		TemplateSystemCategoryId: int32(d.Get("system_category_id").(int)),
		TemplateDefaultTypeId:    int32(d.Get("default_type_id").(int)),
		TemplateMinimumTypeId:    int32(d.Get("minimum_type_id").(int)),
	}
	winId, isSet := d.GetOk("windows_type_id")
	if isSet {
		createCommand.TemplateWindowsTypeId = int32(winId.(int))
	}
	supPass, isSet := d.GetOk("support_password")
	if isSet {
		createCommand.TechSupportPassword = supPass.(string)
	}

	instance_id, ok := d.GetOk("instance_id")
	if !ok {
		return diag.Errorf("Instance id must be specified.")
	}

	tflog.Debug(ctx, "calling ODK OCIApi.InstancesConvertToTemplate")
	ticket, _, err := client.OCIApi.InstancesConvertToTemplate(*auth, int32(instance_id.(int)), createCommand)
	if err != nil {
		return diag.Errorf("ODK Error in OCIApi.InstancesConvertToTemplate. %s", err)
	}

	createTicket, err := waitForTicket(client, auth, ticket)
	if err != nil {
		return diag.Errorf("ODK Error in TicketsApi.TicketsGet. %s", err)
	}
	if createTicket.Status.Id != DICT_TICKET_SUCCEED {
		return diag.Errorf("Unable to create template. Ticket status=%v", createTicket.Status.Id)
	}

	tflog.Info(ctx, fmt.Sprintf("successfully created template for instance id=%v", createTicket.ObjectId))

	template, _, err := client.OCIApi.InstancesGetTemplateByBaseVirtualMachineId(*auth, int32(instance_id.(int)), nil)
	if err != nil {
		return diag.Errorf("Instance with id %d no found. Caused by: %s", instance_id.(int), err)
	}

	d.SetId(strconv.Itoa(int(template.Id)))

	return resourceTemplateRead(ctx, d, m)
}

func resourceTemplateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "reading template")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	templateId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Invalid template id: %v %s", d.Id(), err)
	}

	tflog.Debug(ctx, "calling ODK OCITemplatesApi.TemplatesGet_1")
	template, _, err := client.OCITemplatesApi.TemplatesGet_1(*auth, int32(templateId), nil)
	if err != nil {
		return diag.Errorf("template with id %d not found. Api response: %s", templateId, err)
	}

	return loadTemplateData(ctx, d, m, template)
}

func resourceTemplateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "updating template")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	templateId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Invalid template id: %v %s", d.Id(), err)
	}

	updateCommand := odk.UpdateTemplateCommand{}
	updateNeeded := false
	if d.HasChange("name") {
		updateCommand.Name = d.Get("name").(string)
		updateNeeded = true
	}
	if d.HasChange("description") {
		descriptions := make([]odk.TemplateDescription, 0)
		descriptions = append(descriptions, odk.TemplateDescription{
			LanguageId:  DICT_LANGUAGE_PL,
			Description: d.Get("description").(string),
		})
		descriptions = append(descriptions, odk.TemplateDescription{
			LanguageId:  DICT_LANGUAGE_EN,
			Description: d.Get("description").(string),
		})
		updateCommand.TemplateDescriptions = descriptions
		updateNeeded = true
	}
	if updateNeeded {
		tflog.Debug(ctx, "calling ODK OCITemplatesApi.TemplatesPut")
		_, _, err := client.OCITemplatesApi.TemplatesPut(*auth, int32(templateId), updateCommand)
		if err != nil {
			return diag.Errorf("ODK Error in OCITemplatesApi.TemplatesPut. %s", err)
		}
	}

	return resourceTemplateRead(ctx, d, m)
}

func resourceTemplateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "deleting template")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	templateId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Invalid template id: %v %s", d.Id(), err)
	}

	tflog.Debug(ctx, "calling ODK OCITemplatesApi.TemplatesDelete")
	_, _, err = client.OCITemplatesApi.TemplatesDelete(*auth, int32(templateId))
	if err != nil && err.Error() != "EOF" { // "EOF" condition is a patch for ODK 1.4 bug: it reports error when API returns empty body
		return diag.Errorf("ODK Error in OCITemplatesApi.TemplatesDelete. %s", err)
	}

	d.SetId("")
	return nil
}

func loadTemplateData(ctx context.Context, d *schema.ResourceData, m interface{}, template odk.Template) diag.Diagnostics {
	// Prepare maps
	disks := make(map[string]interface{})
	for _, disk := range template.Disks {
		disks[strconv.Itoa(int(disk.Id))] = disk.Name
	}
	software := make(map[string]interface{})
	for _, soft := range template.Software {
		software[strconv.Itoa(int(soft.Id))] = soft.Name
	}

	// Store everything
	tflog.Debug(ctx, "Parsing returned data")
	if d.Set("name", template.Name) != nil {
		return diag.Errorf("Can't retrieve template name")
	}
	if d.Set("description", template.Description) != nil {
		return diag.Errorf("Can't retrieve descriptions")
	}
	if d.Set("version", template.Version) != nil {
		return diag.Errorf("Can't retrieve version")
	}
	if d.Set("system_category_id", int(template.SystemCategory.Id)) != nil {
		return diag.Errorf("Can't retrieve system category id")
	}
	if d.Set("default_type_id", int(template.DefaultInstanceType.Id)) != nil {
		return diag.Errorf("Can't retrieve default type id")
	}
	if d.Set("minimum_type_id", int(template.MinimumInstanceType.Id)) != nil {
		return diag.Errorf("Can't retrieve minimum type id")
	}
	if d.Set("creation_date", template.CreationDate.String()) != nil {
		return diag.Errorf("Can't retrieve creation date")
	}
	if d.Set("last_change_date", template.LastChangeDate.String()) != nil {
		return diag.Errorf("Can't retrieve last change date")
	}
	if d.Set("creation_user_id", int(template.CreationUser.Id)) != nil {
		return diag.Errorf("Can't retrieve creation user id")
	}
	if d.Set("ethernet_controllers_number", int(template.EthernetControllersNumber)) != nil {
		return diag.Errorf("Can't retrieve ethernet controllers number")
	}
	if d.Set("ethernet_controllers_type_id", int(template.EthernetControllersType.Id)) != nil {
		return diag.Errorf("Can't retrieve ethernet controllers type id")
	}
	if d.Set("publication_status_id", int(template.PublicationStatus.Id)) != nil {
		return diag.Errorf("Can't retrieve publication status id")
	}
	if d.Set("disks", disks) != nil {
		return diag.Errorf("Can't retrieve disks")
	}
	if d.Set("software", software) != nil {
		return diag.Errorf("Can't retrieve software")
	}
	if d.Set("template_type_id", int(template.TemplateType.Id)) != nil {
		return diag.Errorf("Can't retrieve template type id")
	}
	return nil
}
