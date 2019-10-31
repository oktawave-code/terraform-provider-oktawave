package main

import (
	"context"
	"github.com/oktawave-code/odk"
	swagger "github.com/oktawave-code/oks-sdk"
)

type ClientConfig struct {
	ctx       *context.Context
	oksCtx    *context.Context
	odkClient odk.APIClient
	oksClient swagger.APIClient
}

func (c *ClientConfig) oktaClient() odk.APIClient { return c.odkClient }

func (c *ClientConfig) oktaOKSClient() swagger.APIClient { return c.oksClient }

func (c *ClientConfig) getOKSAuth() *context.Context { return c.oksCtx }
