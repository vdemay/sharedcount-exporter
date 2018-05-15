package sharedCount

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/common/log"
	"io/ioutil"
	"net/http"
)

// Client defines the Speedtest client
type Client struct {
	apiKey string
}

type Response struct {
	Quota_used_today      float64 `json:"quota_used_today"`
	Quota_remaining_today float64 `json:"quota_remaining_today"`
	Quota_allocated_today float64 `json:"quota_allocated_today"`
	Plan                  string  `json:"plan"`
}

func NewClient(apiKey string) (*Client, error) {
	log.Infof("Retrieve configuration")
	// TODO check API Key
	return &Client{
		apiKey: apiKey,
	}, nil
}

func (client *Client) Metrics() (Response, error) {
	result := Response{}

	url := fmt.Sprintf("https://scoopit.sharedcount.com/v1.0/quota?apikey=%s", client.apiKey)
	response, err := http.Get(url)

	if err != nil {
		return result, err
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return result, err
		} else {
			json.Unmarshal([]byte(contents), &result)
			return result, nil
		}
	}
	return result, nil
}
