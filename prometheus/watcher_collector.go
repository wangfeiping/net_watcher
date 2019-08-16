package prometheus

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var collector *watcherCollector

func init() {
	collector = &watcherCollector{
		serviceStatusDesc: prometheus.NewDesc(
			"network_service_status",
			"Status code of network service response ",
			[]string{"url"}, nil),
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

	for url, status := range c.mapper {
		ch <- prometheus.MustNewConstMetric(
			c.serviceStatusDesc,
			prometheus.GaugeValue,
			float64(status), url)
	}
}

func (c *watcherCollector) setStatusCode(url string, status int) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.mapper[url] = status
}

// SetStatusCode set status code to the collector mapper
func SetStatusCode(url string, status int) {
	collector.setStatusCode(url, status)
}
