package main

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"

	"github.com/wangfeiping/log"
	"github.com/wangfeiping/net_watcher/commands"
	"github.com/wangfeiping/net_watcher/config"
	"github.com/wangfeiping/net_watcher/exporter"
	"github.com/wangfeiping/net_watcher/util"
)

var starter = func() (cancel context.CancelFunc, err error) {
	log.Info("Start watch...")

	running := true
	t := time.NewTicker(time.Duration(
		viper.GetInt64(commands.FlegDuration)) * time.Second)
	var wg sync.WaitGroup
	cancel = func() {
		running = false
		t.Stop()
		log.Debug("...")
		wg.Wait()
		log.Info("Stop watch")
	}

	config.Load()

	go func() {
		wg.Add(1)
		doJob()
		for running {
			select {
			case <-t.C:
				{
					doJob()
				}
			default:
				{
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
		log.Debug("Done")
		wg.Done()
	}()

	prometheus.MustRegister(exporter.Collector())

	http.Handle("/metrics", promhttp.Handler())
	listen := viper.GetString(commands.FlagListen)
	err = http.ListenAndServe(listen, nil)
	log.Error(err)
	return
}

func doJob() {
	srvs := config.GetServices()
	log.Debugf("Do watch: %d", len(srvs))

	for _, srv := range srvs {
		status, cost := util.HTTPCall(srv.Url)
		exporter.SetStatusCode(srv, status, cost)
	}
}
