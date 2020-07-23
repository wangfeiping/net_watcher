package util

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"

	// "io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/wangfeiping/log"
	"github.com/wangfeiping/net_watcher/config"
)

// Call request http(s) service
func Call(srv *config.Service) (status int, cost int64, resp string) {
	cost = time.Now().UnixNano()
	status, resp = doCall(srv)
	if status > 0 {
		cost = time.Now().UnixNano() - cost
		cost = cost / 1000000
		log.Infof("Success, status: %d, cost: %d, resp: %s", status, cost, resp)
	} else {
		cost = 0
	}
	return
}

func doCall(srv *config.Service) (status int, response string) {
	var err error
	var resp *http.Response
	if strings.HasPrefix(srv.Url, "https://") {
		resp, err = insecureCall(srv)
	} else {
		resp, err = http.Get(srv.Url)
	}
	if err != nil {
		log.Error("Failed, request error: ", err.Error())
		return
	}
	defer resp.Body.Close()
	buf := bytes.NewBuffer(nil)
	_, err = io.CopyN(buf, resp.Body, 100)
	if err != nil && err != io.EOF {
		log.Error("Failed, read response error: ", err.Error())
		return
	}
	response = string(buf.Bytes())
	response = strings.ReplaceAll(response, "\n", "")
	response = strings.ReplaceAll(response, "\r", "")
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Failed, status: %d, resp: %s", resp.StatusCode, response)
		log.Error(err.Error())
		return
	}
	status = resp.StatusCode
	return
}

// insecureCall request http(s) service with InsecureSkipVerify
func insecureCall(srv *config.Service) (resp *http.Response, err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return client.Get(srv.Url)
}
