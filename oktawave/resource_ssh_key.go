package oktawave

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

func resourceSshKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSshKeyCreate,
		ReadContext:   resourceSshKeyRead,
		DeleteContext: resourceSshKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of ssh key.",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Public ssh key",
			},
			"owner_user_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Id of user who created this resource.",
			},
			"creation_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation date.",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},
		Description: "Ssh keys - used for connections security.",
	}
}

func resourceSshKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "creating ssh key")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	sshKeyName := d.Get("name").(string)
	sshKeyValue := d.Get("value").(string)

	tflog.Debug(ctx, "calling ODK AccountApi.AccountPostSshKey", map[string]interface{}{"name": sshKeyName})
	createSshKeyCommand := odk.CreateSshKeyCommand{
		SshKeyName: sshKeyName,
		SshKey:     sshKeyValue,
	}
	sshKey, _, err := client.AccountApi.AccountPostSshKey(*auth, createSshKeyCommand)
	if err != nil {
		return diag.Errorf("ODK Error in AccountApi.AccountPostSshKey. %s", err)
	}
	d.SetId(strconv.Itoa(int(sshKey.Id)))

	tflog.Debug(ctx, "ssh key created")
	return resourceSshKeyRead(ctx, d, m)
}

func resourceSshKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "reading ssh key")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Invalid ssh key id: %v %s", d.Id(), err)
	}

	tflog.Debug(ctx, "calling ODK AccountApi.AccountGetSshKey", map[string]interface{}{"id": id})
	sshKey, resp, err := client.AccountApi.AccountGetSshKey(*auth, int32(id), nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return diag.Errorf("Ssh key %v not found", id)
		}
		return diag.Errorf("Error while retrieving ssh key %v: %s", id, err)
	}

	tflog.Debug(ctx, "Parsing returned data")
	if d.Set("name", sshKey.Name) != nil {
		return diag.Errorf("Can't retrieve ssh key name")
	}
	if d.Set("value", sshKey.Value) != nil {
		return diag.Errorf("Can't retrieve ssh key value")
	}
	if d.Set("owner_user_id", sshKey.OwnerUser.Id) != nil {
		return diag.Errorf("Can't retrieve ssh key owner user id")
	}
	if d.Set("creation_date", sshKey.CreationDate.String()) != nil {
		return diag.Errorf("Can't retrieve ssh key creation date")
	}

	tflog.Debug(ctx, "ssh key retrieved")
	return *new(diag.Diagnostics)
}

func resourceSshKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "deleting ssh key")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Invalid ssh key id: %v %s", d.Id(), err)
	}

	tflog.Debug(ctx, "calling ODK AccountApi.AccountDeleteSshKey", map[string]interface{}{"id": d.Id()})
	resp, err := client.AccountApi.AccountDeleteSshKey(*auth, int32(id))
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return diag.Errorf("Ssh key %v not found", d.Id())
		}
		return diag.Errorf("Error while retrieving ssh key %v: %s", d.Id(), err)
	}

	d.SetId("")
	tflog.Debug(ctx, "ssh key deleted")
	return *new(diag.Diagnostics)
}
