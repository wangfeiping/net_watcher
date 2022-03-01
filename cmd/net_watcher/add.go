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
	srv := checkService()

	srvs, err := addService(srv)
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
			zap.Field{Key: "url", String: srv.Url, Type: zapcore.StringType})
	}
	return
}

func addService(srv *config.Service) (srvs []*config.Service, err error) {
	if err = viper.UnmarshalKey(commands.FlagService, &srvs); err != nil {
		log.Errorf("Unmarshal config error: %v", err)
		return
	}
	srvs = append(srvs, srv)
	return
}

func checkService() *config.Service {
	url := viper.GetString(commands.FlagURL)
	alias := viper.GetString(commands.FlagAlias)
	body := viper.GetString(commands.FlagBody)
	method := viper.GetString(commands.FlagMethod)
	regex := viper.GetString(commands.FlagRegex)

	srv := &config.Service{
		Alias: alias, Url: url, Method: method, Body: body, Regex: regex}
	log.Debugf("checking service: %s - %s %s", srv.Alias, srv.Method, srv.Url)
	return srv
}

func newPath(file string, i int) (newFile string, err error) {
	newFile = fmt.Sprintf("%s.%d", file, i)
	log.Debug("Backup config: ", newFile)
	_, err = os.Stat(newFile)
	return newFile, err
}
