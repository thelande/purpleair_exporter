package main

import (
	"net/http"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/thelande/purpleair_exporter/purpleair"

	kingpin "github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

var (
	configFile = kingpin.Flag(
		"config.file",
		"Path to configuration file.",
	).Default("config.yaml").String()
	webConfig = webflag.AddFlags(kingpin.CommandLine, ":9811")
	logger    log.Logger
)

func main() {
	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Version(version.Print("purpleair_exporter"))
	kingpin.Parse()
	logger = promlog.New(promlogConfig)
	level.Info(logger).Log("msg", "Starting purpleair_exporter", "version", version.Info())

	config := PurpleAirConfig{}
	if err := config.ReloadConfig(*configFile); err != nil {
		level.Error(logger).Log("msg", "Failed to load configuration file", "err", err)
		os.Exit(1)
	}

	client := purpleair.PurpleAirClient{ApiKey: config.ApiKey, Logger: logger}
	fields := append(config.Fields, "name")
	collector, err := NewPurpleAirExporter(config.SensorIndices, fields, &client)
	if err != nil {
		level.Error(logger).Log("msg", "Failed to create exporter", "err", err)
		os.Exit(1)
	}
	prometheus.MustRegister(collector)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
            <head>
            <title>PurpleAir Exporter</title>
            <style>
            label{
            display:inline-block;
            width:75px;
            }
            form label {
            margin: 10px;
            }
            form input {
            margin: 10px;
            }
            </style>
            </head>
            <body>
            <h1>PurpleAir Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
            </body>
            </html>`))
	})

	srv := &http.Server{}
	if err := web.ListenAndServe(srv, webConfig, logger); err != nil {
		level.Error(logger).Log("msg", "HTTP listener stopped", "error", err)
		os.Exit(1)
	}
}
