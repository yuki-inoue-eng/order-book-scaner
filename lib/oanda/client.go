package oanda

import (
	"fmt"
	"github.com/yuki-inoue-eng/trade-force/backend/lib"
	"io/ioutil"
	"net/http"
	"time"
)

const authorizationPrefix = "Bearer "

// Client implements operations trade of oanda through OANDA API.
type Client struct {
	client          *http.Client
	endpoint        string
	requiredHeaders http.Header
}

// NewClient constructs OANDA API client objects.
func NewClient(apiKey string, environment string) Client {
	requiredHeaders := http.Header{}
	requiredHeaders.Add("Authorization", authorizationPrefix+apiKey)
	requiredHeaders.Add("Content-Type", "application/json")
	var endpoint string
	if environment == "Trade" {
		endpoint = "https://api-fxtrade.oanda.com"
	} else {
		endpoint = "https://api-fxpractice.oanda.com"
	}
	return Client{
		client:          &http.Client{},
		endpoint:        endpoint,
		requiredHeaders: requiredHeaders,
	}
}

func (c *Client) fetchOrderBook(instrument Instrument, dateTime *time.Time) ([]byte, error) {
	url := c.endpoint + "/v3/instruments/" + string(instrument) + "/orderBook"
	if dateTime != nil {
		url = c.endpoint + "/v3/instruments/" + string(instrument) + "/orderBook?time=" + dateTime.UTC().Format(time.RFC3339Nano)
	}
	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %v", err)
	}
	c.requiredHeaders.Add("Accept-Datetime-Format", "RFC3339")
	req.Header = c.requiredHeaders
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch response: %v", err)
	}
	defer lib.SafeClose(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTTP %s: failed to read response body: %v", resp.Status, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %s: %s", resp.Status, body)
	}
	return body, nil
}

func (c *Client) fetchPositionBook(instrument Instrument, dateTime *time.Time) ([]byte, error) {
	url := c.endpoint + "/v3/instruments/" + string(instrument) + "/positionBook"
	if dateTime != nil {
		url = c.endpoint + "/v3/instruments/" + string(instrument) + "/positionBook?time=" + dateTime.UTC().Format(time.RFC3339Nano)
	}
	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %v", err)
	}
	c.requiredHeaders.Add("Accept-Datetime-Format", "RFC3339")
	req.Header = c.requiredHeaders
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch response: %v", err)
	}
	defer lib.SafeClose(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTTP %s: failed to read response body: %v", resp.Status, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %s: %s", resp.Status, body)
	}
	return body, nil
}