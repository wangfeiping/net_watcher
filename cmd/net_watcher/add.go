package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/wangfeiping/log"
	"github.com/wangfeiping/net_watcher/commands"
	"github.com/wangfeiping/net_watcher/config"
)

var addHandler = func() (cancel context.CancelFunc, err error) {
	url := viper.GetString(commands.FlagURL)
	alias := viper.GetString(commands.FlagAlias)
	body := viper.GetString(commands.FlagBody)
	log.Debugf("New service: %s - %s", alias, url)

	srvs, err := addService(url, alias, body)
	if err != nil {
		return
	}
	viper.Set(commands.FlagService, srvs)
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
	v.Set(commands.FlagService, srvs)
	err = v.WriteConfig()
	if err != nil {
		log.Errorf("Failed: write config file error: %v", err)
	} else {
		log.Infoz("Success: config add service",
			zap.Field{Key: "url", String: url, Type: zapcore.StringType})
	}
	return
}

func addService(url, alias, body string) (srvs []*config.Service, err error) {
	if err = viper.UnmarshalKey(commands.FlagService, &srvs); err != nil {
		log.Errorf("Unmarshal config error: %v", err)
		return
	}
	// for i, u := range srvs {
	// 	log.Debugf("Config urls: %d, %s", i, u)
	// 	if strings.EqualFold(u, url) {
	// 		err = fmt.Errorf("Service exist: %s", url)
	// 		log.Warn(err)
	// 		return
	// 	}
	// }
	service := &config.Service{
		Alias: alias, Url: url, Body: body}
	srvs = append(srvs, service)
	return
}

func newPath(file string, i int) (newFile string, err error) {
	newFile = fmt.Sprintf("%s.%d", file, i)
	log.Debug("Backup config: ", newFile)
	_, err = os.Stat(newFile)
	return newFile, err
}
