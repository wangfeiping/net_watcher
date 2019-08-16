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

	doJob()

	go func() {
		wg.Add(1)
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