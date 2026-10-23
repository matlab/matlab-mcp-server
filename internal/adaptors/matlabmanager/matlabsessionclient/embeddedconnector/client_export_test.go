// Copyright 2025-2026 The MathWorks, Inc.

package embeddedconnector

import (
	httpclient "github.com/matlab/matlab-mcp-server/internal/adaptors/http/client"
)

type MATLABException = matlabException

type MATLABExceptionsError = matlabExceptionsError

func (c *Client) SetHttpClient(httpClient httpclient.HttpClient) {
	c.httpClient = httpClient
}

func (c *Client) SetHost(host string) {
	c.host = host
}

func (c *Client) SetPort(port string) {
	c.port = port
}

func (c *Client) SetBasePath(basePath string) {
	c.basePath = basePath
}

func (c *Client) MessageServiceEndpoint(channel string) (string, error) {
	return c.messageServiceEndpoint(channel)
}
