package main

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/wangfeiping/net_watcher/commands"
	"github.com/wangfeiping/net_watcher/exporter"
	"github.com/wangfeiping/net_watcher/util"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/wangfeiping/net_watcher/log"
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
	reloadServices()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		reloadServices()
		log.Info("Config file changed:", e.Name)
	})

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
	urls := getServices()
	log.Debugf("Do watch: %d", len(urls))

	for _, u := range urls {
		status, cost := util.HTTPCall(u)
		exporter.SetStatusCode(u, status, cost)
		log.Debugf("Call: status %d, cost %d, %s", status, cost, u)
	}
}

var mux sync.RWMutex
var urls []string

func reloadServices() {
	mux.Lock()
	defer mux.Unlock()

	if err := viper.UnmarshalKey(commands.FlagService, &urls); err != nil {
		log.Errorf("Reload config error: %v", err)
		return
	}
}

func getServices() []string {
	mux.RLock()
	defer mux.RUnlock()

	return urls
}
