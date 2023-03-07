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
	oks "github.com/oktawave-code/oks-sdk"
)

func resourceOksNode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOksNodeCreate,
		ReadContext:   resourceOksNodeRead,
		UpdateContext: resourceOksNodeUpdate,
		DeleteContext: resourceOksNodeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster id.",
			},
			"subregion_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "ID from subregions resource",
			},
			"type_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Value from dictionary #12",
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Value from dictionary #27",
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"total_disks_capacity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cpu_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ram_mb": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},
	}
}

func resourceOksNodeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "creating oks node")

	client := m.(*ClientConfig).oksClient
	auth := m.(*ClientConfig).oksAuth
	odkClient := m.(*ClientConfig).odkClient
	odkAuth := m.(*ClientConfig).odkAuth

	createCmd := oks.K44SNodesSpecification{
		Nodes: []oks.Node{
			{
				Subregion: float64(d.Get("subregion_id").(int)),
				Type_:     float64(d.Get("type_id").(int)),
			},
		},
	}
	tflog.Debug(ctx, "calling OKS ClustersApi.ClustersInstancesNamePost")
	operations, _, err := client.ClustersApi.ClustersInstancesNamePost(*auth, createCmd, d.Get("cluster_id").(string))
	if err != nil {
		return diag.Errorf("OKS Error in ClustersApi.ClustersInstancesNamePost. %s", err)
	}
	if operations[0].Error_ != "" {
		return diag.Errorf("OKS Error in ClustersApi.ClustersInstancesNamePost. %s", operations[0].Error_)
	}

	ticket := odk.Ticket{EndDate: time.Time{}, Progress: 0, Id: int64(operations[0].Ticket.Id)}
	createTicket, err := waitForTicket(odkClient, odkAuth, ticket)
	if err != nil {
		return diag.Errorf("ODK Error in TicketsApi.TicketsGet. %s", err)
	}
	if createTicket.Status.Id != DICT_TICKET_SUCCEED {
		return diag.Errorf("Unable to create node. Ticket status=%v", createTicket.Status.Id)
	}

	tflog.Info(ctx, fmt.Sprintf("successfully created OKS node. id=%v", createTicket.ObjectId))
	d.SetId(strconv.Itoa(int(createTicket.ObjectId)))

	return resourceOksNodeRead(ctx, d, m)
}

func resourceOksNodeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return diag.Errorf("Oks node update not implemented.")
}

func resourceOksNodeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "reading oks node")

	client := m.(*ClientConfig).oksClient
	auth := m.(*ClientConfig).oksAuth

	instanceId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Invalid node id: %v %s", d.Id(), err)
	}

	clusterId, ok := d.GetOk("cluster_id")
	if !ok {
		return diag.Errorf("Cluster id must be specified.")
	}

	tflog.Debug(ctx, "calling OKS ClustersApi.ClustersInstancesNameGet")
	nodes, _, err := client.ClustersApi.ClustersInstancesNameGet(*auth, clusterId.(string))
	if err != nil {
		return diag.Errorf("OKS Error in ClustersApi.ClustersInstancesNameGet. %s", err)
	}

	for _, node := range nodes {
		if int(node.Id) == instanceId {
			return loadOksNodeData(ctx, d, m, clusterId.(string), node)
		}
	}
	d.SetId("")
	return diag.Errorf("Node not found")
}

func resourceOksNodeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "deleting oks node")

	client := m.(*ClientConfig).oksClient
	auth := m.(*ClientConfig).oksAuth
	odkClient := m.(*ClientConfig).odkClient
	odkAuth := m.(*ClientConfig).odkAuth

	instanceId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Invalid node id: %v %s", d.Id(), err)
	}

	clusterId, ok := d.GetOk("cluster_id")
	if !ok {
		return diag.Errorf("Cluster id must be specified.")
	}

	nodes := oks.K44SNodesList{
		InstancesIds: []float64{float64(instanceId)},
	}
	tflog.Debug(ctx, "calling OKS ClustersApi.ClustersInstancesNameDelete")
	operations, _, err := client.ClustersApi.ClustersInstancesNameDelete(*auth, nodes, clusterId.(string))
	if err != nil {
		return diag.Errorf("ODK Error in ClustersApi.ClustersInstancesNameDelete. %s", err)
	}

	ticket := odk.Ticket{EndDate: time.Time{}, Progress: 0, Id: int64(operations[0].Ticket.Id)}
	deleteTicket, err := waitForTicket(odkClient, odkAuth, ticket)
	if err != nil {
		return diag.Errorf("ODK Error in TicketsApi.TicketsGet. %s", err)
	}
	if deleteTicket.Status.Id != DICT_TICKET_SUCCEED {
		return diag.Errorf("Unable to delete node. Ticket status=%v", deleteTicket.Status.Id)
	}

	d.SetId("")
	return nil
}

func loadOksNodeData(ctx context.Context, d *schema.ResourceData, m interface{}, clusterId string, node oks.K44sInstance) diag.Diagnostics {
	// Store everything
	tflog.Debug(ctx, "Parsing returned data")
	if d.Set("cluster_id", clusterId) != nil {
		return diag.Errorf("Can't retrieve cluster id")
	}
	if d.Set("name", node.Name) != nil {
		return diag.Errorf("Can't retrieve name")
	}
	if d.Set("creation_date", node.CreationDate.String()) != nil {
		return diag.Errorf("Can't retrieve creation date")
	}
	if d.Set("subregion_id", int(node.Subregion.Id)) != nil {
		return diag.Errorf("Can't retrieve subregion id")
	}
	if d.Set("type_id", node.Type_.Id) != nil {
		return diag.Errorf("Can't retrieve type id")
	}
	if d.Set("status_id", node.Status.Id) != nil {
		return diag.Errorf("Can't retrieve status")
	}
	if d.Set("ip_address", node.IpAddress) != nil {
		return diag.Errorf("Can't retrieve ip address")
	}
	if d.Set("total_disks_capacity", node.TotalDisksCapacity) != nil {
		return diag.Errorf("Can't retrieve total disks capacity")
	}
	if d.Set("cpu_number", node.CpuNumber) != nil {
		return diag.Errorf("Can't retrieve cpu number")
	}
	if d.Set("ram_mb", node.RamMb) != nil {
		return diag.Errorf("Can't retrieve ram size")
	}
	return nil
}
