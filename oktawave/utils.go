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

func castToInt32(ids []interface{}) []int32 {
	int32Ids := make([]int32, len(ids))
	for i, v := range ids {
		int32Ids[i] = (int32)(v.(int))
	}
	return int32Ids
}

func castToString(ips []interface{}) []string {
	strings := make([]string, len(ips))
	for i, v := range ips {
		strings[i] = (string)(v.(string))
	}
	return strings
}

func castIntToInt32(ids []int) []int32 {
	int32Ids := make([]int32, len(ids))
	for i, v := range ids {
		int32Ids[i] = (int32)(v)
	}
	return int32Ids
}

func calcListAMinusListB(listA []int32, listB []int32) []int32 {
	diff := make([]int32, 0)
	for _, aId := range listA {
		found := false
		for _, bId := range listB {
			if aId == bId {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, aId)
		}
	}
	return diff
}

func calcListAMinusListB_string(listA []string, listB []string) []string {
	diff := make([]string, 0)
	for _, aId := range listA {
		found := false
		for _, bId := range listB {
			if aId == bId {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, aId)
		}
	}
	return diff
}

func waitForTicket(client odk.APIClient, auth *context.Context, ticket odk.Ticket) (odk.Ticket, error) {
	tflog.Info(context.Background(), fmt.Sprintf("Waiting for ticket %v", ticket.Id))
	var maxRetries = 5
	var err error
	first := true
	for ticket.EndDate.IsZero() {
		if !first {
			tflog.Debug(context.Background(), fmt.Sprintf("Still waiting (ticket=%v; progress=%v)", ticket.Id, ticket.Progress))
			time.Sleep(10 * time.Second)
		}
		first = false
		tflog.Debug(context.Background(), "calling ODK TicketsApi.TicketsGet")
		ticket, _, err = client.TicketsApi.TicketsGet_1(*auth, ticket.Id, nil)
		if err != nil {
			tflog.Warn(context.Background(), fmt.Sprintf("ODK Error in TicketsApi.TicketsGet. %v", err))
			if maxRetries <= 0 {
				return ticket, err
			}
			maxRetries--
		}
	}
	return ticket, nil
}

func detachIpById(client odk.APIClient, auth *context.Context, instanceId int32, ipId int32) (odk.Ticket, *http.Response, error) {
	tflog.Debug(context.Background(), "calling ODK OCIInterfacesApi.InstancesPostDetachIpTicket")
	ticket, resp, err := client.OCIInterfacesApi.InstancesPostDetachIpTicket(*auth, instanceId, int32(ipId))
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return ticket, resp, fmt.Errorf("instance id %v or IP id %v was not found", instanceId, ipId)
		}
		return ticket, resp, fmt.Errorf("ODK Error in OCIInterfacesApi.InstancesPostDetachIpTicket. %v", err)
	}
	ticket, err = waitForTicket(client, auth, ticket)
	return ticket, nil, err
}

func attachIpById(client odk.APIClient, auth *context.Context, instanceId int32, ipId int32) (odk.Ticket, *http.Response, error) {
	localOptions := map[string]interface{}{
		"ipId": int32(ipId),
	}
	tflog.Debug(context.Background(), "calling ODK OCIInterfacesApi.InstancesPostAttachIpTicket")
	ticket, resp, err := client.OCIInterfacesApi.InstancesPostAttachIpTicket(*auth, instanceId, localOptions)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return ticket, resp, fmt.Errorf("instance id %v or IP id %v was not found", instanceId, ipId)
		}
		return ticket, resp, fmt.Errorf("ODK Error in OCIInterfacesApi.InstancesPostAttachIpTicket. %v", err)
	}
	ticket, err = waitForTicket(client, auth, ticket)
	return ticket, nil, err
}

// func detachIp(client odk.APIClient, auth *context.Context, instanceId int32, ip string) (odk.Ticket, *http.Response, error) {
// 	tflog.Debug(context.Background(), "calling ODK FloatingIPsApi.FloatingIpsPostDetachIpTicket")
// 	ticket, resp, err := client.FloatingIPsApi.FloatingIpsPostDetachIpTicket(*auth, ip, instanceId)
// 	if err != nil {
// 		if resp != nil && resp.StatusCode == http.StatusNotFound {
// 			return ticket, resp, fmt.Errorf("instance id %v or IP %v was not found", instanceId, ip)
// 		}
// 		return ticket, resp, fmt.Errorf("ODK Error in FloatingIPsApi.FloatingIpsPostDetachIpTicket. %v", err)
// 	}
// 	ticket, err = waitForTicket(client, auth, ticket)
// 	return ticket, nil, err
// }

// func attachIp(client odk.APIClient, auth *context.Context, instanceId int32, ip string) (odk.Ticket, *http.Response, error) {
// 	localOptions := map[string]interface{}{
// 		"ipV4": ip,
// 	}
// 	tflog.Debug(context.Background(), "calling ODK FloatingIPsApi.FloatingIpsPostAttachIpTicket")
// 	ticket, resp, err := client.FloatingIPsApi.FloatingIpsPostAttachIpTicket(*auth, instanceId, localOptions)
// 	if err != nil {
// 		if resp != nil && resp.StatusCode == http.StatusNotFound {
// 			return ticket, resp, fmt.Errorf("instance id %v or IP %v was not found", instanceId, ip)
// 		}
// 		return ticket, resp, fmt.Errorf("ODK Error in FloatingIPsApi.FloatingIpsPostAttachIpTicket. %v", err)
// 	}
// 	ticket, err = waitForTicket(client, auth, ticket)
// 	return ticket, nil, err
// }

func getConnectionInstanceIds(connections []odk.DiskConnection) []int {
	connectionIds := make([]int, len(connections))
	for i, conn := range connections {
		connectionIds[i] = int(conn.Instance.Id)
	}
	return connectionIds
}

// ------------

func makeDataSourceSchema(name string, getSchema func() map[string]*schema.Schema) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter": {
			Type: schema.TypeSet,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Required: true,
						// todo
						// ValidateDiagFunc:func(v any, p cty.Path) diag.Diagnostics {}
					},
					"values": {
						Type:     schema.TypeList,
						Required: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
			Optional: true,
		},
		name: {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: getSchema(),
			},
		},
	}
}

func makeDataSourceRead[T any](
	dataSourceName string,
	dataSourceSchema map[string]*schema.Schema,
	getData func(config *ClientConfig) ([]T, error),
	mapRawDataToDataSourceModel func(rawElem T) (map[string]interface{}, error),
) func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

		config := m.(*ClientConfig)
		rawDataItems, err := getData(config)
		if err != nil {
			return diag.FromErr(err)
		}

		// map results
		results := make([]map[string]interface{}, len(rawDataItems))
		for i, rawDataItem := range rawDataItems {
			result, err := mapRawDataToDataSourceModel(rawDataItem)
			if err != nil {
				return diag.FromErr(err)
			}
			results[i] = result
		}

		// filter results
		if filters, ok := d.GetOk("filter"); ok {
			for _, filterEntry := range filters.(*schema.Set).List() {

				key := filterEntry.(map[string]interface{})["key"].(string)
				values := filterEntry.(map[string]interface{})["values"].([]interface{})
				dss := dataSourceSchema[dataSourceName].Elem.(*schema.Resource).Schema
				valuesToTest, err := castStringsToValues(values, dss[key])
				if err != nil {
					return diag.FromErr(err)
				}

				filterFn := func(item map[string]interface{}) bool {
					res := false
					for _, v := range valuesToTest {
						res = res || testValue(dss[key], item[key], v)
					}
					return res
				}

				results = filter(results, filterFn)
			}
		}

		if err := d.Set(dataSourceName, results); err != nil {
			return diag.Errorf("Set %s failed: %s", dataSourceName, err)
		}

		d.SetId(time.Now().UTC().String())
		return nil
	}
}

func castStringsToValues(values []interface{}, s *schema.Schema) ([]interface{}, error) {
	results := []interface{}{}
	for _, value := range values {
		v, err := castStringToValue(value.(string), s)
		if err != nil {
			return nil, err
		}
		results = append(results, v)
	}
	return results, nil
}

func castStringToValue(value string, s *schema.Schema) (interface{}, error) {
	switch s.Type {
	case schema.TypeString:
		return value, nil
	case schema.TypeBool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return nil, err
		}
		return boolValue, nil
	case schema.TypeInt:
		intValue, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return nil, err
		}
		return int32(intValue), nil
	case schema.TypeFloat:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		return floatValue, nil
	}
	return nil, fmt.Errorf("unrecognized type for data conversion")
}

func testValue(s *schema.Schema, item interface{}, value interface{}) bool {
	switch s.Type {
	case schema.TypeString:
		return item.(string) == value.(string)
	case schema.TypeBool:
		return item.(bool) == value.(bool)
	case schema.TypeInt:
		return item.(int32) == value.(int32)
	case schema.TypeFloat:
		return math.Abs(item.(float64)-value.(float64)) < 0.001
	case schema.TypeList:
		return false
	case schema.TypeSet:
		return false
	case schema.TypeMap:
		// fixme - return error or at least log warning
		return false
	}
	return false
}

func filter[T any](items []T, fn func(item T) bool) []T {
	filteredItems := []T{}
	for _, value := range items {
		if fn(value) {
			filteredItems = append(filteredItems, value)
		}
	}
	return filteredItems
}
