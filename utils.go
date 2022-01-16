package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/oktawave-code/odk"
	swagger "github.com/oktawave-code/oks-sdk"
)

const (
	TICKET_STATUS__SUCCESS int32 = 136
	TICKET_STATUS__ERROR   int32 = 137
)

func evaluateTicket(client odk.APIClient, auth *context.Context, ticket odk.Ticket) (odk.Ticket, error) {
	log.Printf("[DEBUG] Starting ticket polling. Ticket id: [%d], Target object name [%s].", ticket.Id, ticket.ObjectName)
	currentTicket := ticket
	var max_retries = 5

	for currentTicket.EndDate.IsZero() {
		time.Sleep(10 * time.Second)
		returnedTicket, resp, err := client.TicketsApi.TicketsGet_1(*auth, currentTicket.Id, nil)
		log.Printf("[TRACE] Ticket API HTTP response. Ticket Id: [%d], HTTP response status: [%d]", currentTicket.Id, resp.StatusCode)

		if err != nil {
			if max_retries <= 0 {
				log.Printf("[ERROR] Ticket polling failed. Server did not respond too many times. Tickdt Id: [%d], Target object name: [%s].", currentTicket.Id, currentTicket.ObjectName)
				return currentTicket, err
			}
			log.Printf("[DEBUG] Ticket polling - api request failed. Will retry in a moment. Tickdt Id: [%d], Target object name: [%s].", currentTicket.Id, currentTicket.ObjectName)
			max_retries--
		}

		if resp.Body != nil {
			currentTicket = returnedTicket
		}
		log.Printf("[DEBUG] Ticket polling still in progress. Ticket Id: [%d], Target object name: [%s], EndDate: [%s], Progress: [%d]", currentTicket.Id, currentTicket.ObjectName, currentTicket.EndDate, currentTicket.Progress)
	}
	log.Printf("[DEBUG] Ticket polling finished. Ticket Id: [%d], Status id: [%d], Target object name: [%s], EndDate: [%s], Progress: [%d]", currentTicket.Id, currentTicket.Status.Id, currentTicket.ObjectName, currentTicket.EndDate, currentTicket.Progress)
	return currentTicket, nil
}

func retrieve_ids(ids []interface{}) []int32 {
	int32Ids := make([]int32, len(ids))
	for i := 0; i < len(ids); i++ {
		intInstance_id := ids[i].(int)
		int32Instance_id := (int32)(intInstance_id)
		int32Ids[i] = int32Instance_id
	}
	return int32Ids
}

func detachIpById(client odk.APIClient, auth *context.Context, instanceId int32, ipId int32) (odk.Ticket, *http.Response, error) {
	ticket, resp, err := client.OCIInterfacesApi.InstancesPostDetachIpTicket(*auth, instanceId, int32(ipId))
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return ticket, resp, fmt.Errorf("Instance by id %s or "+
				"ip address by id %s were not found to detach", strconv.Itoa(int(instanceId)), strconv.Itoa(int(ipId)))
		}
		return ticket, resp, fmt.Errorf("Error occured while detaching ip %s", err)
	}

	detachTicket, err := evaluateTicket(client, auth, ticket)
	return detachTicket, nil, err
}

func attachIpById(client odk.APIClient, auth *context.Context, instanceId int32, ipId int32) (odk.Ticket, *http.Response, error) {
	localOptions := map[string]interface{}{
		"ipId": int32(ipId),
	}
	ticket, resp, err := client.OCIInterfacesApi.InstancesPostAttachIpTicket(*auth, instanceId, localOptions)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return ticket, resp, fmt.Errorf(" Instance by id %s or "+
				"ip address by id %s were not found to attach", strconv.Itoa(int(instanceId)), strconv.Itoa(int(ipId)))
		}
		return ticket, resp, fmt.Errorf("Resource AddressIp. UPDATE. Error occured while attaching ip %s", err)
	}

	attachTicket, err := evaluateTicket(client, auth, ticket)
	return attachTicket, nil, err
}

func getConnectionsInstancesIds(connections []odk.DiskConnection) []int {
	connectionIds := make([]int, len(connections))
	for i := 0; i < len(connections); i++ {
		connectionIds[i] = int(connections[i].Instance.Id)
	}
	return connectionIds
}

func getConnectionsInstancesIds_int32(connections []odk.DiskConnection) []int32 {
	connectionIds := make([]int32, len(connections))
	for i := 0; i < len(connections); i++ {
		connectionIds[i] = connections[i].Instance.Id
	}
	return connectionIds
}

func retrieveNodeById(nodes []swagger.K44sInstance, nodeId int) (swagger.K44sInstance, error) {
	for _, node := range nodes {
		if nodeId == int(node.Id) {
			return node, nil
		}
	}

	return swagger.K44sInstance{}, fmt.Errorf("Node by id %s was not found", strconv.Itoa(nodeId))
}
