package config

import (
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/wangfeiping/log"
	"github.com/wangfeiping/net_watcher/commands"
)

type Service struct {
	Alias   string   `json:"alias,omitempty" yaml:"alias,omitempty"`
	Url     string   `json:"url" yaml:"url"`
	Method  string   `json:"method,omitempty" yaml:"method,omitempty"`
	Body    string   `json:"body,omitempty" yaml:"body,omitempty"`
	Regex   string   `json:"regex,omitempty" yaml:"regex,omitempty"`
	Service *Service `json:"service,omitempty" yaml:"service,omitempty"`
}

var mux sync.RWMutex
var srvs []*Service

func ReloadServices() {
	mux.Lock()
	defer mux.Unlock()

	if err := viper.UnmarshalKey(commands.FlagService, &srvs); err != nil {
		log.Errorf("Load config error: %v", err)
		return
	}
}

func GetServices() []*Service {
	mux.RLock()
	defer mux.RUnlock()

	return srvs
}

func Load() {
	ReloadServices()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		ReloadServices()
		log.Info("Config file changed:", e.Name)
	})
}

func Check(data string, val ...string) string {
	return val[0] + data[6:]
}
