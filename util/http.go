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
)

// HTTPCall request http(s) service
func HTTPCall(url string) (status int, cost int64) {
	cost = time.Now().UnixNano()
	var body string
	status, body = doHTTPCall(url)
	if status > 0 {
		cost = time.Now().UnixNano() - cost
		cost = cost / 1000000
		log.Infof("Success, status: %d, cost: %d, body: %s", status, cost, body)
	} else {
		cost = 0
	}
	return
}

func doHTTPCall(URL string) (status int, body string) {
	resp, err := http.Get(URL)
	if err != nil {
		log.Error("Failed, request error: ", err.Error())
		if strings.Index(err.Error(),
			"x509: certificate signed by unknown authority") < 0 {
			return
		}
		resp, err = insecureHTTPCall(URL)
		if err != nil {
			log.Error("Failed, insecure request error: ", err.Error())
			return
		}
	}
	defer resp.Body.Close()
	status = resp.StatusCode
	buf := bytes.NewBuffer(nil)
	_, err = io.CopyN(buf, resp.Body, 100)
	// response, err = ioutil.ReadAll(resp.Body)
	if err != nil && err != io.EOF {
		log.Error("Failed, read response error: ", err.Error())
		return
	}
	s := string(buf.Bytes())
	ss := strings.Split(s, "\n")
	var buffer bytes.Buffer
	for _, s = range ss {
		buffer.WriteString(s)
	}
	body = strings.ReplaceAll(buffer.String(), "\r", "")
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Failed, status: %d, body: %s", resp.StatusCode, body)
		log.Error(err.Error())
		return
	}
	return
}

// insecureHTTPCall request http(s) service with InsecureSkipVerify
func insecureHTTPCall(URL string) (resp *http.Response, err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return client.Get(URL)
}
