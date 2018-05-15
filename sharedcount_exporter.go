package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	prom_version "github.com/prometheus/common/version"

	"github.com/vdemay/sharedcount-exporter/sharedCount"
)

const (
	namespace = "sharedcount"
)

var (
	quota_used_today = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "quota_used_today"),
		"Used today",
		nil, nil,
	)
	quota_remaining_today = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "quota_remaining_today"),
		"Remaining today",
		nil, nil,
	)
	quota_allocated_today = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "quota_allocated_today"),
		"Allocated",
		nil, nil,
	)
)

// Exporter collects Speedtest stats from the given server and exports them using
// the prometheus metrics package.
type Exporter struct {
	Client *sharedCount.Client
}

// NewExporter returns an initialized Exporter.
func NewExporter(apikey string, interval time.Duration) (*Exporter, error) {
	log.Infof("Setup sharedcount client with interval %s", interval)
	client, err := sharedCount.NewClient(apikey)
	if err != nil {
		return nil, fmt.Errorf("Can't create the sharedcount client: %s", err)
	}

	log.Debugln("Init exporter")
	return &Exporter{
		Client: client,
	}, nil
}

// Describe describes all the metrics ever exported by the Speedtest exporter.
// It implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- quota_used_today
	ch <- quota_remaining_today
	ch <- quota_allocated_today
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	log.Infof("sharedcount exporter starting")
	if e.Client == nil {
		log.Errorf("sharedcount client not configured.")
		return
	}

	metrics := e.Client.Metrics()
	ch <- prometheus.MustNewConstMetric(quota_used_today, prometheus.GaugeValue, metrics.Quota_used_today)
	ch <- prometheus.MustNewConstMetric(quota_remaining_today, prometheus.GaugeValue, metrics.Quota_remaining_today)
	ch <- prometheus.MustNewConstMetric(quota_allocated_today, prometheus.GaugeValue, metrics.Quota_allocated_today)
	log.Infof("sharedcount exporter finished")
}

func init() {
	prometheus.MustRegister(prom_version.NewCollector("sharecount_exporter"))
}

func main() {
	var (
		showVersion   = flag.Bool("version", false, "Print version information.")
		listenAddress = flag.String("web.listen-address", ":9112", "Address to listen on for web interface and telemetry.")
		metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		apiKey        = flag.String("sharedcount.apikey", "undefined", "aikey from sharedCount")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("sharedcount Prometheus exporter. v%s\n", 1.0)
		os.Exit(0)
	}

	log.Infoln("sharedcount speedtest exporter", prom_version.Info())
	log.Infoln("Build context", prom_version.BuildContext())

	interval := 60 * time.Second
	exporter, err := NewExporter(*apiKey, interval)
	if err != nil {
		log.Errorf("Can't create exporter : %s", err)
		os.Exit(1)
	}
	log.Infoln("Register exporter")
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>SharedCount Exporter</title></head>
             <body>
             <h1>SharedCount Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	log.Infoln("Listening on", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
