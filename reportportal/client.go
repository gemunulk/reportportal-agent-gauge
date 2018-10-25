//Package reportportal provides a go api for reportportal
package reportportal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// A Client stores the client informations
// and implement all the api functions
// to communicate with testrail
type Client struct {
	url        string
	password   string
	httpClient *http.Client
}

// NewClient returns a new client
// with the given credential
// for the given testrail domain
func NewClient(url, password string) (c *Client) {
	c = &Client{}
	c.password = "bearer " + password

	c.url = url
	if !strings.HasSuffix(c.url, "/") {
		c.url += "/"
	}
	c.url += "api/v1/"

	c.httpClient = &http.Client{}

	return
}

// sendRequest sends a request of type "method"
// to the url "client.url+uri" and with optional data "data"
// Returns an error if any and the optional data "v"
func (c *Client) sendRequest(method, uri string, data, v interface{}) error {
	//fmt.Println(c.url + uri)

	var body io.Reader
	if data != nil {

		jsonReq, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("marshaling data: %s", err)
		}

		body = bytes.NewBuffer(jsonReq)

		//body = bytes.NewBuffer(data)
		//fmt.Println(body)

	}

	req, err := http.NewRequest(method, c.url+uri, body)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", c.password)
	//req.SetBasicAuth(c.username, c.password)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	jsonCnt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading: %s", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("response: status: %q, body: %s", resp.Status, jsonCnt)
	}

	if v != nil {
		err = json.Unmarshal(jsonCnt, v)
		if err != nil {
			return fmt.Errorf("unmarshaling response: %s", err)
		}
	}

	return nil
}
