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
	data string
}

type Response struct {
	Quota_used_today      float64 `json:"quota_used_today"`
	Quota_remaining_today float64 `json:"quota_remaining_today"`
	Quota_allocated_today float64 `json:"quota_allocated_today"`
	Plan                  string  `json:"plan"`
}

func NewClient(apiKey string) (*Client, error) {

	log.Debug("Retrieve configuration")
	url := fmt.Sprintf("https://scoopit.sharedcount.com/v1.0/quota?apikey=%s", apiKey)
	log.Debug(url)
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		} else {
			log.Debug(string(contents))
			return &Client{
				data: string(contents),
			}, nil
		}
	}
	return nil, nil
}

func (client *Client) Metrics() Response {
	result := Response{}
	json.Unmarshal([]byte(client.data), &result)
	return result
}
