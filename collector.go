package main

import (
	"errors"
	"fmt"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/thelande/purpleair_exporter/purpleair"
)

const namespace = "purpleair"

var (
	upDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"'1' if the API was successfully scraped, '0' otherwise.",
		[]string{"sensor_index"},
		nil,
	)
	infoDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "info"),
		"Information about the sensor.",
		[]string{"sensor_index", "name"},
		nil,
	)
	sensorValueDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "sensor", "count"),
		"The value of the sensor.",
		[]string{"sensor_index", "field"},
		nil,
	)
)

type PurpleAirExporter struct {
	SensorIndices []int
	Fields        []string
	Client        *purpleair.PurpleAirClient
	Logger        *log.Logger
}

func NewPurpleAirExporter(
	indices []int,
	fields []string,
	client *purpleair.PurpleAirClient,
) (*PurpleAirExporter, error) {
	if len(indices) == 0 {
		return nil, errors.New("no indices provided")
	}

	return &PurpleAirExporter{SensorIndices: indices, Fields: fields, Client: client, Logger: &logger}, nil
}

func (c PurpleAirExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- upDesc
	ch <- infoDesc
	ch <- sensorValueDesc
}

func (c PurpleAirExporter) Collect(ch chan<- prometheus.Metric) {
	for _, index := range c.SensorIndices {
		resp, err := c.Client.GetSensor(index, &c.Fields)
		sensorIndex := fmt.Sprintf("%d", index)
		if err != nil {
			ch <- prometheus.MustNewConstMetric(
				upDesc, prometheus.GaugeValue, 0, sensorIndex,
			)
		} else {
			ch <- prometheus.MustNewConstMetric(
				infoDesc, prometheus.GaugeValue, 1, sensorIndex, resp.Sensor.Name,
			)
			ch <- prometheus.MustNewConstMetric(
				sensorValueDesc, prometheus.GaugeValue, resp.Sensor.Pm2_5, sensorIndex, "pm2.5",
			)
			ch <- prometheus.MustNewConstMetric(
				upDesc, prometheus.GaugeValue, 1, sensorIndex,
			)
		}
	}
}
