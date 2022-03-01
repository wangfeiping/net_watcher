package main

import (
	"context"

	"github.com/wangfeiping/net_watcher/util"
)

var callHandler = func() (context.CancelFunc, error) {
	srv := checkService()

	util.Call(srv)
	return nil, nil
}
