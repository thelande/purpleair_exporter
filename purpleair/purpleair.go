package purpleair

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

const apiUrl = "https://api.purpleair.com/v1"

type PurpleAirClient struct {
	ApiKey string
	Logger log.Logger
}

type PurpleAirResponse struct {
	ApiVersion    string `json:"api_version"`
	TimeStamp     int    `json:"time_stamp"`
	DataTimeStamp int    `json:"data_time_stamp"`
}

type SensorResponse struct {
	SensorIndex int     `json:"sensor_index"`
	Name        string  `json:"name"`
	Pm2_5       float64 `json:"pm2.5"`
}

type PurpleAirSensorResponse struct {
	PurpleAirResponse
	Sensor SensorResponse `json:"sensor"`
}

func (p *PurpleAirClient) doRequest(method string, path string, values *url.Values) ([]byte, error) {
	url, err := url.Parse(fmt.Sprintf("%s/%s", apiUrl, path))
	url.RawQuery = values.Encode()

	client := http.Client{}

	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		level.Error(p.Logger).Log("msg", "Failed to create new http request", "err", err)
		return nil, err
	}

	req.Header.Add("X-API-Key", p.ApiKey)
	resp, err := client.Do(req)
	if err != nil {
		level.Error(p.Logger).Log("msg", "Failed to query PurpleAir API", "err", err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		level.Error(p.Logger).Log("msg", "Failed to read response body", "err", err)
		return nil, err
	}
	return body, nil
}

func (p PurpleAirClient) GetPrivateSensor(sensorIndex int, readKey string, fields *[]string) (*PurpleAirSensorResponse, error) {
	path := fmt.Sprintf("%s/%d", "sensors", sensorIndex)
	values := url.Values{}

	if readKey != "" {
		values.Add("read_key", readKey)
	}

	if len(*fields) > 0 {
		values.Add("fields", strings.Join(*fields, ","))
	}

	body, err := p.doRequest("GET", path, &values)
	if err != nil {
		return nil, err
	}

	resp := PurpleAirSensorResponse{}
	if err := json.Unmarshal(body, &resp); err != nil {
		level.Error(p.Logger).Log("msg", "Failed to unmarshal response JSON", "err", err)
		return nil, err
	}

	return &resp, nil
}

func (p PurpleAirClient) GetSensor(sensorIndex int, fields *[]string) (*PurpleAirSensorResponse, error) {
	return p.GetPrivateSensor(sensorIndex, "", fields)
}
