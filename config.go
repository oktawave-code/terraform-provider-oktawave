package main

import (
	"context"
	"github.com/oktawave-code/odk"
)

type ClientConfig struct {
	ctx       *context.Context
	odkClient odk.APIClient
}

func (c *ClientConfig) oktaClient() odk.APIClient { return c.odkClient }
