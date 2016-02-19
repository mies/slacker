package slacker

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// DefaultAPIURL is the default URL + Path for slack API requests
const DefaultAPIURL = "https://slack.com/api"

// APIClient contains simple logic for starting the RTM Messaging API for Slack
type APIClient struct {
	client   *http.Client
	token    string
	SlackURL string
}

// NewAPIClient returns a new APIClient with a token set.
func NewAPIClient(token string, url string) *APIClient {
	if url == "" {
		url = DefaultAPIURL
	}

	return &APIClient{
		client:   http.DefaultClient,
		token:    token,
		SlackURL: url,
	}
}

// RunMethod runs an RPC method and returns the response body as a byte slice
func (c *APIClient) RunMethod(name string, params ...string) ([]byte, error) {
	resp, err := c.slackMethod(name, params...)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (c *APIClient) slackMethod(method string, params ...string) (*http.Response, error) {
	query := fmt.Sprintf("%s/%s?token=%s", c.SlackURL, method, c.token)
	for param := range params {
		query = fmt.Sprintf("%s/?=%s", param)
	}
	req, err := http.NewRequest("GET", query, nil)

	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

func (c *APIClient) slackMethodAndParse(method string, dest interface{}) error {
	resp, err := c.slackMethod(method)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return ParseResponse(resp.Body, dest)
}
