package oktawave

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/oktawave-code/odk"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the group.",
			},
			"affinity_rule_type_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     DICT_AFFINITY_TYPE_NO_SEPARATION,
				Description: "This parameter defines how instances in group are allocated amongst hosts. Minimum affinity allocates instances on single host lowering latency. Maximum affinity allocates instances on multiple hosts minimizing failures risk. Value from dictionary #160",
			},
			"assignment": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Instance id.",
						},
						"ip_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Only one IP of the instance can be assigned to group.",
						},
					},
				},
				Description: "Assigns instance to group.",
			},
			"is_load_balancer": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "When this parameter is set to true, connected load balancer transfers requests to instances in this group.",
			},
			"instances_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of instances in this group.",
			},
			"schedulers_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of schedulers in this group.",
			},
			"autoscaling_type_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Value from dictionary #55",
			},
			"last_change_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when last change was applied.",
			},
			"creation_user_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Id of user who created this resource.",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},
		Description: "Groups are used for configuring instances separation level, load balancing and autoscaling.",
	}
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "creating group")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	createCommand := odk.CreateGroupCommand{
		Name:             d.Get("name").(string),
		AffinityRuleType: int32(d.Get("affinity_rule_type_id").(int)),
	}

	tflog.Debug(ctx, "calling ODK OCIGroupsApi.GroupsCreate")
	group, _, err := client.OCIGroupsApi.GroupsCreate(*auth, createCommand)
	if err != nil {
		return diag.Errorf("ODK Error in OCIGroupsApi.GroupsCreate. %s", err)
	}

	assignments, assignmentsSet := d.GetOk("assignment")
	if assignmentsSet {
		assignmentCommand := odk.ChangeContainerAssignmentsCommand{
			Assignments: createAssignments(assignments),
		}

		tflog.Debug(ctx, "calling ODK OCIGroupsApi.GroupsChangeAssignmentsInGroup")
		_, _, err = client.OCIGroupsApi.GroupsChangeAssignmentsInGroup(*auth, group.Id, assignmentCommand)
		if err != nil {
			return diag.Errorf("Group was created, but instance assignments failed. %s", err)
		}
	}

	d.SetId(strconv.Itoa(int(group.Id)))
	return resourceGroupRead(ctx, d, m)
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "reading group")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	groupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Invalid group id: %v %s", d.Id(), err)
	}

	tflog.Debug(ctx, "calling ODK OCIGroupsApi.GroupsGetGroup")
	group, _, err := client.OCIGroupsApi.GroupsGetGroup(*auth, int32(groupId), nil)
	if err != nil {
		return diag.Errorf("ODK Error in OCIGroupsApi.GroupsGetGroup. %s", err)
	}

	return loadGroupData(ctx, client, auth, d, m, group)
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "updating group")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	groupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Invalid group id: %v %s", d.Id(), err)
	}

	d.Partial(true)

	updateCommand := odk.CreateGroupCommand{}
	updateNeeded := false
	if d.HasChange("name") {
		updateCommand.Name = d.Get("name").(string)
		updateNeeded = true
	}
	if d.HasChange("affinity_rule_type_id") {
		updateCommand.AffinityRuleType = int32(d.Get("affinity_rule_type_id").(int))
		updateNeeded = true
	}
	if updateNeeded {
		tflog.Debug(ctx, "calling ODK OCIGroupsApi.GroupsUpdate")
		_, _, err := client.OCIGroupsApi.GroupsUpdate(*auth, int32(groupId), updateCommand)
		if err != nil {
			return diag.Errorf("ODK Error in OCIGroupsApi.GroupsUpdate. %s", err)
		}
	}

	if d.HasChange("assignment") {
		assignments := d.Get("assignment")
		assignmentCommand := odk.ChangeContainerAssignmentsCommand{
			Assignments: createAssignments(assignments),
		}
		tflog.Debug(ctx, "calling ODK OCIGroupsApi.GroupsChangeAssignmentsInGroup")
		_, _, err = client.OCIGroupsApi.GroupsChangeAssignmentsInGroup(*auth, int32(groupId), assignmentCommand)
		if err != nil {
			return diag.Errorf("ODK Error in OCIGroupsApi.GroupsChangeAssignmentsInGroup. %s", err)
		}
	}

	d.Partial(false)
	return resourceGroupRead(ctx, d, m)
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "deleting group")

	client := m.(*ClientConfig).odkClient
	auth := m.(*ClientConfig).odkAuth

	groupId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Invalid group id: %v %s", d.Id(), err)
	}

	tflog.Debug(ctx, "calling ODK OCIGroupsApi.GroupsDelete")
	_, _, err = client.OCIGroupsApi.GroupsDelete(*auth, int32(groupId))
	if err != nil && err.Error() != "EOF" { // "EOF" condition is a patch for ODK 1.4 bug: it reports error when API returns empty body
		return diag.Errorf("ODK Error in OCIGroupsApi.GroupsDelete. %s", err)
	}

	d.SetId("")
	return nil
}

func loadGroupData(ctx context.Context, client odk.APIClient, auth *context.Context, d *schema.ResourceData, m interface{}, group odk.Group) diag.Diagnostics {
	// Load assignments
	assignments, err := getGroupAssignments(ctx, client, auth, int32(group.Id))
	if err != nil {
		return diag.Errorf("%s", err)
	}

	// Store everything
	tflog.Debug(ctx, "Parsing returned data")
	if d.Set("name", group.Name) != nil {
		return diag.Errorf("Can't retrieve group name")
	}
	if d.Set("affinity_rule_type_id", int(group.AffinityRuleType.Id)) != nil {
		return diag.Errorf("Can't retrieve affinity rule type id")
	}
	if d.Set("assignment", assignments) != nil {
		return diag.Errorf("Can't retrieve group assignments")
	}
	if d.Set("is_load_balancer", group.IsLoadBalancer) != nil {
		return diag.Errorf("Can't retrieve load balancer state")
	}
	if d.Set("instances_count", group.InstancesCount) != nil {
		return diag.Errorf("Can't retrieve instances count")
	}
	if d.Set("schedulers_count", group.SchedulersCount) != nil {
		return diag.Errorf("Can't retrieve schedulers count")
	}
	if d.Set("autoscaling_type_id", int(group.AutoscalingType.Id)) != nil {
		return diag.Errorf("Can't retrieve autoscaling type id")
	}
	if d.Set("last_change_date", group.LastChangeDate.String()) != nil {
		return diag.Errorf("Can't retrieve last change date")
	}
	if d.Set("creation_user_id", int(group.CreationUser.Id)) != nil {
		return diag.Errorf("Can't retrieve creation user id")
	}
	return nil
}

func createAssignments(assignments interface{}) []odk.ContainerAssignmentCommand {
	var list []odk.ContainerAssignmentCommand
	for _, assignment := range assignments.(*schema.Set).List() {
		list = append(list, odk.ContainerAssignmentCommand{
			InstanceId: int32(assignment.(map[string]interface{})["instance_id"].(int)),
			IpId:       int32(assignment.(map[string]interface{})["ip_id"].(int)),
		})
	}
	return list
}

func getGroupAssignments(ctx context.Context, client odk.APIClient, auth *context.Context, groupId int32) ([]map[string]interface{}, error) {
	tflog.Debug(ctx, "calling ODK OCIGroupsApi.GroupsGetAssignmentsInGroup")
	params := map[string]interface{}{
		"pageSize": int32(math.MaxInt16),
	}
	assignmentsGroupCollection, _, err := client.OCIGroupsApi.GroupsGetAssignmentsInGroup(*auth, groupId, params)
	if err != nil {
		return nil, fmt.Errorf("ODK Error in OCIGroupsApi.GroupsGetAssignmentsInGroup. %s", err)
	}

	var assignments []map[string]interface{}
	for _, assignment := range assignmentsGroupCollection.Items {
		assignments = append(assignments, map[string]interface{}{
			"ip_id":       assignment.IpId,
			"instance_id": assignment.InstanceId,
		})
	}
	return assignments, nil
}
