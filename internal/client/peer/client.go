/******************************************************************************
*
*  Copyright 2023 SAP SE
*
*  Licensed under the Apache License, Version 2.0 (the "License");
*  you may not use this file except in compliance with the License.
*  You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*  Unless required by applicable law or agreed to in writing, software
*  distributed under the License is distributed on an "AS IS" BASIS,
*  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*  See the License for the specific language governing permissions and
*  limitations under the License.
*
******************************************************************************/

package peerclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sapcc/keppel/internal/auth"
	"github.com/sapcc/keppel/internal/keppel"
)

// Client can be used for API access to one of our peers (using our peering
// credentials).
type Client struct {
	peer  keppel.Peer
	token string
}

// New obtains a token for API access to the given peer (using our peering
// credentials), and wraps it into a Client instance.
func New(cfg keppel.Configuration, peer keppel.Peer, scope auth.Scope) (Client, error) {
	c := Client{peer, ""}
	err := c.initToken(cfg, scope)
	if err != nil {
		return Client{}, fmt.Errorf("while trying to obtain a peer token for %s in scope %s: %w",
			peer.HostName, scope.String(), err)
	}
	return c, nil
}

func (c *Client) initToken(cfg keppel.Configuration, scope auth.Scope) error {
	reqURL := c.buildRequestURL(fmt.Sprintf("keppel/v1/auth?service=%[1]s&scope=%[2]s", c.peer.HostName, scope.String()))
	ourUserName := "replication@" + cfg.APIPublicHostname
	authHeader := map[string]string{"Authorization": keppel.BuildBasicAuthHeader(ourUserName, c.peer.OurPassword)}

	respBodyBytes, respStatusCode, _, err := c.doRequest(http.MethodGet, reqURL, http.NoBody, authHeader)
	if err != nil {
		return err
	}
	if respStatusCode != http.StatusOK {
		return fmt.Errorf("expected 200 OK, but got %d: %s", respStatusCode, strings.TrimSpace(string(respBodyBytes)))
	}

	var data struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(respBodyBytes, &data)
	if err != nil {
		return err
	}
	c.token = data.Token
	return nil
}

func (c Client) buildRequestURL(path string) string {
	return fmt.Sprintf("https://%s/%s", c.peer.HostName, path)
}

func (c Client) doRequest(method, url string, body io.Reader, headers map[string]string) (respBodyBytes []byte, respStatusCode int, respHeader http.Header, err error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("during %s %s: %w", method, url, err)
	}
	if c.token != "" { //empty token occurs only during initToken()
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("during %s %s: %w", method, url, err)
	}
	defer resp.Body.Close()
	respBodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("during %s %s: %w", method, url, err)
	}

	return respBodyBytes, resp.StatusCode, resp.Header, nil
}
