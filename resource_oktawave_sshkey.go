package main

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
	"log"
	"net/http"
	"strconv"
	"time"
)

func resourceSshKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceSshKeyCreate,
		Read:   resourceSshKeyRead,
		Delete: resourceSshKeyDelete,
		Schema: map[string]*schema.Schema{
			"ssh_key_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ssh_key_value": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"owner_user_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},
	}
}

func resourceSshKeyCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("Resource SSH Key. CREATE. Initializing")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx

	log.Printf("Resource SSH Key. CREATE. Retrieving attributes from config file")
	sshKeyName := d.Get("ssh_key_name").(string)
	sshKeyValue := d.Get("ssh_key_value").(string)

	log.Printf("Resource SSH key. CREATE. Trying to create ticket to post ssh key")
	createSshKeyCommand := odk.CreateSshKeyCommand{
		SshKeyName: sshKeyName,
		SshKey:     sshKeyValue,
	}

	sshKey, _, err := client.AccountApi.AccountPostSshKey(*auth, createSshKeyCommand)
	if err != nil {
		return fmt.Errorf("Resource Ssh key. CREATE. Cannot post new ssh key. %v", err)
	}
	d.SetId(strconv.Itoa(int(sshKey.Id)))

	log.Printf("Resource Ssh key. CREATE. Created")
	return resourceSshKeyRead(d, m)
}

func resourceSshKeyRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("Resource SSH Key. Read. Initializing")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx
	id, err := strconv.Atoi(d.Id())

	if err != nil {
		return fmt.Errorf("Resource Ssh key. READ. Invalid ssh key id: %v", err)
	}

	log.Printf("Resource Ssh key. READ. Retrieving ssh key form API")
	sshKey, resp, err := client.AccountApi.AccountGetSshKey(*auth, int32(id), nil)

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return fmt.Errorf("Resource Ssh key. READ. Ssh key by id %s was not found", d.Id())
		}
		return fmt.Errorf("Resource Ssh key. READ. Error occured while retrievenig ssh key. %s", err)
	}

	log.Printf("[INFO] Resource Ssh key. READ. Trying to retrieve Ssh key name")
	if err := d.Set("ssh_key_name", sshKey.Name); err != nil {
		return fmt.Errorf("Resource Ssh key. READ. Error: can't retrieve Ssh key name")
	}

	log.Printf("[INFO] Resource Ssh key. READ. Trying to retrieve Ssh key value")
	if err := d.Set("ssh_key_value", sshKey.Value); err != nil {
		return fmt.Errorf("Resource Ssh key. READ. Error: can't retrieve Ssh key value")
	}

	log.Printf("[INFO] Resource Ssh key. READ. Trying to retrieve Ssh key owner user id")
	if err := d.Set("owner_user_id", sshKey.OwnerUser.Id); err != nil {
		return fmt.Errorf("Resource Ssh key. READ. Error: can't retrieve Ssh key owner user id")
	}

	log.Printf("[INFO] Resource Ssh key. READ. Trying to retrieve Ssh key creationg date")
	if err := d.Set("creation_date", sshKey.CreationDate.String()); err != nil {
		return fmt.Errorf("Resource Ssh key. READ. Error: can't retrieve Ssh key creation date")
	}

	return nil
}

func resourceSshKeyDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource ssh key. DELETE. Initializing")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx
	sshKeyId, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Resource Ssh key. DELETE. Invalid ssh key id was set: %s", d.Id())
	}

	log.Printf("[INFO] Resource Ssh key. DELETE. Trying to delete ssh key")
	resp, err := client.AccountApi.AccountDeleteSshKey(*auth, int32(sshKeyId))
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Resource Ssh key. DELETE. Ssh key by id %s was not found", d.Id())
		}
		return fmt.Errorf("Resource Ssh key. DELETE. Error occured while deleting: %s", err)
	}
	return nil
}
