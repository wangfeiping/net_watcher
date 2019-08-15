package main

import (
	"context"
	"sync"
	"time"

	"github.com/wangfeiping/net_watcher/commands"
	"github.com/wangfeiping/net_watcher/util"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/wangfeiping/net_watcher/log"
)

var starter = func() (cancel context.CancelFunc, err error) {
	log.Info("Start watch...")

	t := time.NewTicker(time.Duration(
		viper.GetInt64(commands.FlegDuration)) * time.Second)
	var wg sync.WaitGroup
	cancel = func() {
		t.Stop()
		log.Info("...")
		wg.Wait()
		log.Info("Stop watch")
	}
	reloadServices()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		reloadServices()
		log.Info("Config file changed:", e.Name)
	})

	doJob()

	go func() {
		for {
			select {
			case <-t.C:
				{
					wg.Add(1)
					doJob()
					wg.Done()
				}
			}
		}
	}()
	return
}

func doJob() {
	urls := getServices()
	log.Debugf("Do watch: %d", len(urls))

	for _, u := range urls {
		ok, _ := util.HTTPCall(u)
		log.Debugf("%t, %s", ok, u)
	}
}

var mux sync.RWMutex
var urls []string

func reloadServices() {
	mux.Lock()
	defer mux.Unlock()

	if err := viper.UnmarshalKey("service", &urls); err != nil {
		log.Errorf("Reload config error: %v", err)
		return
	}
}

func getServices() []string {
	mux.RLock()
	defer mux.RUnlock()

	return urls
}
