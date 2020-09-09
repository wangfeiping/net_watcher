package util

import (
	"testing"
)

func Test_Capture(t *testing.T) {
	resp := `{	"info":"Failed to connect to /10.64.72.70:8080 executing POST http://service-ad-post/post/advert/get",	"message":"系统异常",	"result":null,	"status":"50000003"}`
	_, err := Capture(resp, "url")
	if err == nil {
		t.Error("Test_Capture should return an err")
	}
	t.Log("ok")
}
