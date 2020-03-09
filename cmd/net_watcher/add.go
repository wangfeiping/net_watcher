package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/wangfeiping/log"
	"github.com/wangfeiping/net_watcher/commands"
)

var addHandler = func() (cancel context.CancelFunc, err error) {
	url := viper.GetString(commands.FlagURL)
	log.Debug("New service: ", url)

	urls, err := addService(url)
	if err != nil {
		return
	}

	log.Debugf("Config add url: %d, %s", len(urls), url)
	viper.Set(commands.FlagService, urls)
	c := viper.GetString(commands.FlagConfig)
	log.Debug("Config file: ", c)
	if _, err = os.Stat(c); err == nil {
		var newFile string
		for i := 1; !os.IsNotExist(err); i++ {
			newFile, err = newPath(c, i)
		}
		os.Rename(c, newFile)
	}
	v := viper.New()
	v.SetConfigFile(viper.GetString(commands.FlagConfig))
	v.Set(commands.FlagService,
		viper.GetStringSlice(commands.FlagService))
	err = v.WriteConfig()
	if err != nil {
		log.Errorf("Write config file error: %v", err)
	}
	return
}

func addService(url string) (urls []string, err error) {
	if err = viper.UnmarshalKey(commands.FlagService, &urls); err != nil {
		log.Errorf("Unmarshal config error: %v", err)
		return
	}
	for i, u := range urls {
		log.Debugf("Config urls: %d, %s", i, u)
		if strings.EqualFold(u, url) {
			err = fmt.Errorf("Service exist: %s", url)
			log.Warn(err)
			return
		}
	}
	urls = append(urls, url)
	return
}

func newPath(file string, i int) (newFile string, err error) {
	newFile = fmt.Sprintf("%s.%d", file, i)
	log.Debug("Backup config: ", newFile)
	_, err = os.Stat(newFile)
	return newFile, err
}
