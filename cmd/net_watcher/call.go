package main

import (
	"context"

	"github.com/spf13/viper"

	"github.com/wangfeiping/net_watcher/commands"
	"github.com/wangfeiping/net_watcher/config"
	"github.com/wangfeiping/net_watcher/util"
)

var callHandler = func() (context.CancelFunc, error) {
	url := viper.GetString(commands.FlagURL)
	alias := viper.GetString(commands.FlagAlias)
	body := viper.GetString(commands.FlagBody)

	service := &config.Service{
		Alias: alias, Url: url, Body: body}

	util.Call(service)
	return nil, nil
}
