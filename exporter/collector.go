package exporter

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/wangfeiping/net_watcher/config"
)

var collector *watcherCollector

func init() {
	collector = &watcherCollector{
		serviceStatusDesc: prometheus.NewDesc(
			"network_service_status",
			"Status of network service response ",
			[]string{"code", "url"}, nil),
		mapper: make(map[string]*callRecord)}
}

type callRecord struct {
	status int
	cost   int64
}

type watcherCollector struct {
	serviceStatusDesc *prometheus.Desc

	mapper map[string]*callRecord
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

	for url, record := range c.mapper {
		ch <- prometheus.MustNewConstMetric(
			c.serviceStatusDesc,
			prometheus.GaugeValue,
			float64(record.cost), fmt.Sprintf("%d", record.status), url)
	}
}

func (c *watcherCollector) setStatusCode(url string, code int, cost int64) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.mapper[url] = &callRecord{
		status: code,
		cost:   cost}
}

// SetStatusCode set status code to the collector mapper
func SetStatusCode(srv *config.Service, code int, cost int64) {
	collector.setStatusCode(srv.Url, code, cost)
}
