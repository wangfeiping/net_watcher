package main

import (
	"context"

	"github.com/spf13/viper"
	"github.com/wangfeiping/net_watcher/commands"
	"github.com/wangfeiping/net_watcher/util"
)

var callHandler = func() (context.CancelFunc, error) {
	util.HTTPCall(viper.GetString(commands.FlagURL))
	return nil, nil
}
