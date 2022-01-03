package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
			}, "affinity_rule_type_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1403,
			},
			"group_instance_ip_ids": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Optional: true,
					Type:     schema.TypeInt,
				},
				Optional: true,
			},
			"schedulers_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"last_change_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instances_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},
	}
}

func resourceGroupCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource Groups. CREATE. Initializing")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx

	log.Printf("[DEBUG] Resource Group. CREATE. Retrieving attributes from config file")
	affinityRuleType := int32(d.Get("affinity_rule_type_id").(int))
	//Creating create command
	createGroupCommand := odk.CreateGroupCommand{
		Name:             d.Get("group_name").(string),
		AffinityRuleType: affinityRuleType,
	}
	group, _, err := client.OCIGroupsApi.GroupsCreate(*auth, createGroupCommand)
	if err != nil {
		return fmt.Errorf("Resource Group. CREATE. Cannot post new group. %v", err)
	}
	log.Printf("[DEBUG] Resource Group. CREATE. Group was created. Retrieving groups instance ids and ip ids..")
	instanceIpIds, instancesIpIdsAreSet := d.GetOk("group_instance_ip_ids")
	//Checking whether instance ip ids map is empty and if not -> creating list of assignments for group and try to make assignment
	if instancesIpIdsAreSet {
		log.Printf("[DEBUG] Resource Group. CREATE. Instance ip ids map was found. Initializing...")
		instanceIpIds := instanceIpIds.(map[string]interface{})
		var assignments []odk.ContainerAssignmentCommand
		log.Printf("[DEBUG] Resource Group. CREATE. Instance ip ids map was found. Getting list of assignments...")
		//create list of assignments from instance ip ids map
		assignments, err = createAssignments(client, *auth, instanceIpIds)
		if err != nil {
			return fmt.Errorf("Resource Group. CREATE. %s ", err)
		}
		log.Printf("[DEBUG] Resource Group. CREATE. Trying to post assignments")
		changeContainerAssignments := odk.ChangeContainerAssignmentsCommand{
			Assignments: assignments,
		}
		//Trying to apply assignments
		_, _, err = client.OCIGroupsApi.GroupsChangeAssignmentsInGroup(*auth, group.Id, changeContainerAssignments)
		if err != nil {
			return fmt.Errorf("[DEBUG] Resource Group. CREATE. Group was created, but instance assignments failed. Error: %s", err)
		}
	}
	d.SetId(strconv.Itoa(int(group.Id)))
	return resourceGroupRead(d, m)
}

func resourceGroupRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource Groups. READ. Initializing")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx
	groupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Resource Group. READ. Invalid group id. %s", err)
	}
	log.Printf("[DEBUG] Trying to retrieve i")
	group, resp, err := client.OCIGroupsApi.GroupsGetGroup(*auth, int32(groupId), nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return fmt.Errorf("Resource Group. READ. Group was not found by id: %s", strconv.Itoa(groupId))
		}
		return fmt.Errorf("Resource Group. Read. READ. Error occured while group retrieving. %s", err)
	}
	log.Printf("[INFO] Resource Group. Group was found. Synch local and remote state...")
	if err = d.Set("group_name", group.Name); err != nil {
		return fmt.Errorf("Resource Group. READ. Error occured while retrieving group name. Error: %s", err)
	}

	if err = d.Set("instances_count", group.InstancesCount); err != nil {
		return fmt.Errorf("Resource Group. READ. Error occured while retrieving group instance count. Error: %s", err)
	}

	if err = d.Set("schedulers_count", group.SchedulersCount); err != nil {
		return fmt.Errorf("Resource Group. READ. Error occured while retrieving group schedulers count. Error: %s", err)
	}

	if err = d.Set("affinity_rule_type_id", int(group.AffinityRuleType.Id)); err != nil {
		return fmt.Errorf("Resource Group. READ. Error occured while retrieving group affinity rule type id. Error: %s", err)
	}

	if err = d.Set("last_change_date", group.LastChangeDate.String()); err != nil {
		return fmt.Errorf("Resource Group. READ. Error occured while retrieving group last change date. Error: %s", err)
	}

	if err = d.Set("creation_user_id", strconv.Itoa(int(group.CreationUser.Id))); err != nil {
		return fmt.Errorf("Resource Group. READ. Error occured while retrieving group creation user id. Error: %s", err)
	}

	assignmentsGroupInstIdsIps, err := getGroupAssignmentIds(client, auth, int32(groupId))
	if err != nil {
		return fmt.Errorf("Resource Group. READ. %s", err)
	}
	if err = d.Set("group_instance_ip_ids", assignmentsGroupInstIdsIps); err != nil {
		return fmt.Errorf("Resource Group. READ. Error occured while retrieving group instance ip ids. Error: %s", err)
	}

	return nil
}

func resourceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource Groups. UPDATE. Initializing")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx
	groupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Resource Group. UPDATE. Invalid group id. %s", err)
	}
	log.Printf("[DEBUG] Resource Group. UPDATE. Check attribute update states")
	d.Partial(true)
	if d.HasChange("group_name") || d.HasChange("affinity_rule_type_id") {
		_, newGroupName := d.GetChange("group_name")
		_, newRuleTypeId := d.GetChange("affinity_rule_type_id")
		updateGroupCommand := odk.CreateGroupCommand{
			Name:             newGroupName.(string),
			AffinityRuleType: int32(newRuleTypeId.(int)),
		}

		_, resp, err := client.OCIGroupsApi.GroupsUpdate(*auth, int32(groupId), updateGroupCommand)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("Resource Group. UPDATE. Group was not found by ip %s", d.Id())
			}
			return fmt.Errorf("Resource Group. UPDATE. Error occured while updating group. %s", err)
		}
		d.SetPartial("group_name")
		d.SetPartial("affinity_rule_type_id")
	}

	if d.HasChange("group_instance_ip_ids") {
		var assignments []odk.ContainerAssignmentCommand
		_, newInstanceIpIds := d.GetChange("group_instance_ip_ids")
		//creating list of assignments from new instance ip ids map
		assignments, err = createAssignments(client, *auth, newInstanceIpIds.(map[string]interface{}))
		changeContainerAssignments := odk.ChangeContainerAssignmentsCommand{
			Assignments: assignments,
		}
		_, _, err = client.OCIGroupsApi.GroupsChangeAssignmentsInGroup(*auth, int32(groupId), changeContainerAssignments)
		if err != nil {
			return fmt.Errorf("ERROR Resource Group. UPDATE. instance assignments failed. Error: %s", err)
		}
		d.SetPartial("group_instance_ip_ids")
	}

	return resourceGroupRead(d, m)
}

func resourceGroupDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Resource Groups. DELETE. Initializing")
	client := m.(*ClientConfig).oktaClient()
	auth := m.(*ClientConfig).ctx
	groupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Resource Group. DELETE. Invalid group id. %s", err)
	}

	_, resp, err := client.OCIGroupsApi.GroupsDelete(*auth, int32(groupId))
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Resource Group. DELETE. Group was not found by id: %s", d.Id())
		}
		if strings.Contains(err.Error(), "EOF") {
			for i := 0; i < 500; i++ {
				time.Sleep(5 * time.Second)
				_, resp, _ := client.OCIGroupsApi.GroupsGetGroup(*auth, int32(groupId), nil)
				if resp != nil && resp.StatusCode == 200 {
					log.Printf("[INFO] Resource Group. DELETE. Group was successfully deleted")
					d.SetId("")
					return nil
				}
				log.Printf("[DEBUG] Resource Group. DELETE. HTTP response code: %s", strconv.Itoa(resp.StatusCode))
			}
		}
		return fmt.Errorf("Resource Group. DELETE. Error occured while deleting group: %s", err)
	}
	log.Printf("[INFO] Resource Group. DELETE. Group was successfully deleted")
	d.SetId("")
	return nil
}

//Create list of assignment to attach to group using map of ids and ips
func createAssignments(client odk.APIClient, auth context.Context, instanceIds map[string]interface{}) ([]odk.ContainerAssignmentCommand, error) {
	assignments := make([]odk.ContainerAssignmentCommand, len(instanceIds))
	log.Printf("[DEBUG] length %s", strconv.Itoa(len(instanceIds)))
	log.Printf("[DEBUG] length %s", strconv.Itoa(len(assignments)))
	count := 0
	//iterating over map, converting key and values and creating new ContainerAssignmentCommand
	for instanceId, ipId := range instanceIds {
		instanceId_int, err := strconv.Atoi(instanceId)
		if err != nil {
			return nil, fmt.Errorf("Invalid instance id %s", err)
		}
		instanceId_int32 := int32(instanceId_int)
		ipId_int32 := int32(ipId.(int))
		//Checking, whether ip id value was set. If not - trying to get first ip id from the list of ip ids of instance
		if ipId_int32 == 0 {
			ipId_int32, err = findFirstIpIdForInstance(client, &auth, instanceId_int32)
			if err != nil {
				return nil, err
			}
			if ipId_int32 == 0 {
				return nil, fmt.Errorf("Resource Group. Ip id was not found for instance id %s", strconv.Itoa(int(instanceId_int32)))
			}
		}
		log.Printf("[DEBUG] Assighning new assignment")
		log.Printf("[DEBUG] Count %s", strconv.Itoa(count))
		assignments[count] = odk.ContainerAssignmentCommand{
			InstanceId: instanceId_int32,
			IpId:       ipId_int32,
		}
		count++
	}
	return assignments, nil
}

//read-only method for getting remote list of ids of instances and list of instances ip ids that were assigned to group
func getGroupAssignmentIds(client odk.APIClient, auth *context.Context, groupId int32) (map[string]interface{}, error) {
	assignmentsGroupCollection, resp, err := client.OCIGroupsApi.GroupsGetAssignmentsInGroup(*auth, int32(groupId), nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("Group not found by id %s", strconv.Itoa(int(groupId)))
		}
		return nil, fmt.Errorf("Error occured while retrieving group assignment list: %s", err)
	}

	assignmentsGroupIds := make(map[string]interface{})
	//iterating over list of assignemnts of group and configuring and on the basis trying to configure instance ip ids local map
	for _, assignment := range assignmentsGroupCollection.Items {
		instanceId_string := strconv.Itoa(int(assignment.InstanceId))
		assignmentsGroupIds[instanceId_string] = assignment.IpId
	}
	return assignmentsGroupIds, nil
}

//Getting first ip id on map of ips of given instance
func findFirstIpIdForInstance(client odk.APIClient, auth *context.Context, instanceId int32) (int32, error) {
	collectionIp, resp, err := client.OCIInterfacesApi.InstancesGetInstanceIps(*auth, instanceId, nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return 0, fmt.Errorf("Ip were not found for instance by id %s", strconv.Itoa(int(instanceId)))
		}
		return 0, fmt.Errorf("Error occured while retrieiving instances ip list %s", err)
	}
	if len(collectionIp.Items) > 0 {
		return collectionIp.Items[0].Id, nil
	} else {
		return 0, nil
	}
}
