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
		log.Infof("Success, status: %d, cost: %d, method: %s, resp: %s", status, cost, srv.Method, resp)
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
			log.Debugf("status: %d, method: %s, resp: %s", resp.StatusCode, srv.Method, response)
			url, err := Capture(response, "url")
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
		resp, err = httpCall(srv)
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
	// TODO
	r, err := regexp.Compile(srv.Regex)
	if err != nil {
		log.Error("regexp error: ", err.Error())
		return
	}
	if strings.EqualFold(srv.Method, "POST") {
		var bytes []byte
		bytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error("Failed, read response error: ", err.Error())
			return
		}
		response = string(bytes)
		if !r.MatchString(string(bytes)) {
			log.Warn("regex match failed: ", response)
			return
		}
		status = resp.StatusCode
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

	req, err := newRequest(srv)

	return client.Do(req)
}

func httpCall(srv *config.Service) (*http.Response, error) {
	req, err := newRequest(srv)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	return client.Do(req)
}

func newRequest(srv *config.Service) (*http.Request, error) {
	method := "GET"
	if srv.Method != "" {
		method = srv.Method
	}

	body := bytes.NewBuffer([]byte(srv.Body))
	req, err := http.NewRequest(method, srv.Url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
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

func Capture(response, name string) (string, error) {
	r := regexp.MustCompile(`"url":"(?P<url>.*?)"`)
	groups := r.FindStringSubmatch(response)
	// fmt.Printf("%#v\n", groups)
	// fmt.Printf("%#v\n", r.SubexpNames())
	// fmt.Printf("%s\n", groups[1])
	l := len(groups)
	for i, n := range r.SubexpNames() {
		if l <= i {
			break
		} else if strings.EqualFold(n, name) {
			return groups[i], nil
		}
	}
	return "", fmt.Errorf("name not found: %s", name)
}
