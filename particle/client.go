package particle

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	AccessToken string

	client *http.Client
}

type DeviceInfoResponse struct {
	ID           string `json:"id"`
	SerialNumber string `json:"serial_number"`
	Name         string `json:"name"`
	LastApp      string `json:"last_app"`
	Connected    bool   `json:"connected"`
	Notes        string `json:"notes"`
}

type ErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"error_description"`
}

func (c *Client) GetDeviceInfo(deviceID string) (*DeviceInfoResponse, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.particle.io/v1/devices/%s", deviceID),
		nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("access_token", c.AccessToken)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	di := &DeviceInfoResponse{}
	if err := json.NewDecoder(resp.Body).Decode(di); err != nil {
		return nil, err
	}

	return di, nil
}
