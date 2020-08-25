package util

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
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
		if err != nil {
			log.Error("Failed, insecure call error: ", err.Error())
			return
		}
		defer resp.Body.Close()
		if srv.Service != nil {
			response, err = read(resp)
			if err != nil {
				log.Error("Failed, read response error: ", err.Error())
				return
			}
			log.Debugf("status: %d, resp: %s", resp.StatusCode, response)
			url, err := capture(response, "url")
			if err != nil {
				log.Error("Failed, capture data error: ", err.Error())
				return
			}
			log.Debugf("Capture url: %s", url)
			service := &config.Service{
				Url: config.Check(srv.Service.Url, url)}
			log.Debugf("Continue call url: %s", service.Url)
			resp, err = insecureCall(service)
		}
	} else {
		resp, err = http.Get(srv.Url)
	}
	if err != nil {
		log.Error("Failed, request error: ", err.Error())
		return
	}
	if resp == nil {
		log.Error("Failed, resp is nil")
		return
	}
	if resp.Body == nil {
		log.Error("Failed, resp.Body is nil")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		response, err = read(resp)
		if err != nil {
			log.Error("Failed, read response error: ", err.Error())
			return
		}
		err = fmt.Errorf("Failed, status: %d, resp: %s", resp.StatusCode, response)
		log.Error(err.Error())
		return
	}
	buf := bytes.NewBuffer(nil)
	_, err = io.CopyN(buf, resp.Body, 100)
	if err != nil && err != io.EOF {
		log.Error("Failed, read response error: ", err.Error())
		return
	}
	response = string(buf.Bytes())
	response = strings.ReplaceAll(response, "\n", "")
	response = strings.ReplaceAll(response, "\r", "")
	status = resp.StatusCode
	return
}

// insecureCall request http(s) service with InsecureSkipVerify
func insecureCall(srv *config.Service) (resp *http.Response, err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	if srv.Body == "" {
		return client.Get(srv.Url)
	}
	req, err := http.NewRequest("GET", srv.Url, bytes.NewBuffer([]byte(srv.Body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return client.Do(req)
}

func read(resp *http.Response) (response string, err error) {
	var bytes []byte
	bytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	response = strings.ReplaceAll(string(bytes), "\n", "")
	response = strings.ReplaceAll(response, "\r", "")
	return
}

func capture(response, name string) (string, error) {
	r := regexp.MustCompile(`"url":"(?P<url>.*?)"`)
	groups := r.FindStringSubmatch(response)
	// fmt.Printf("%#v\n", groups)
	// fmt.Printf("%#v\n", r.SubexpNames())
	// fmt.Printf("%s\n", groups[1])
	for i, n := range r.SubexpNames() {
		if strings.EqualFold(n, name) {
			return groups[i], nil
		}
	}
	return "", fmt.Errorf("name not found: %s", name)
}
