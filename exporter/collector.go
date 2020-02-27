package exporter

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var collector *watcherCollector

func init() {
	collector = &watcherCollector{
		serviceStatusDesc: prometheus.NewDesc(
			"network_service_status",
			"Status code of network service response ",
			[]string{"code", "url"}, nil),
		mapper: make(map[string]int)}
}

type watcherCollector struct {
	serviceStatusDesc *prometheus.Desc

	mapper map[string]int
	mux    sync.RWMutex
}

// Collector returns a collector
// which exports metrics about status code of network service response
func Collector() prometheus.Collector {
	return collector
}

// Describe returns all descriptions of the collector.
func (c *watcherCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.serviceStatusDesc
}

// Collect returns the current state of all metrics of the collector.
func (c *watcherCollector) Collect(ch chan<- prometheus.Metric) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	for url, code := range c.mapper {
		up := 1
		if code == 0 {
			up = 0
		}
		ch <- prometheus.MustNewConstMetric(
			c.serviceStatusDesc,
			prometheus.GaugeValue,
			float64(up), fmt.Sprintf("%d", code), url)
	}
}

func (c *watcherCollector) setStatusCode(url string, code int) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.mapper[url] = code
}

// SetStatusCode set status code to the collector mapper
func SetStatusCode(url string, code int, cost int64) {
	collector.setStatusCode(url, code)
}
