// +build ignore

package main

import (
	"bytes"
	"fmt"
	"strings"

	"net/http"
)

func main() {
	resp, err := http.Get("https://x-token.network/")
	if err != nil {
		fmt.Println("http(s) request error: ", err.Error())
	}
	defer resp.Body.Close()
	// body, _ := ioutil.ReadAll(resp.Body)
	var buf []byte
	buf = make([]byte, 100)
	_, err = resp.Body.Read(buf)
	s := string(buf)
	fmt.Println("s: ", s)
	ss := strings.Split(s, "\n")
	fmt.Println("ss length: ", len(ss))
	// s = strings.Join(ss, "")
	var buffer bytes.Buffer
	for _, s = range ss {
		buffer.WriteString(s)
	}
	s = strings.ReplaceAll(buffer.String(), "\r", "")
	fmt.Printf("code: %d, body: %s\n", resp.StatusCode, s)

	// client := &http.Client{}
	// request, _ := http.NewRequest("GET", "https://x-token.network/", nil)
	// request.Header.Set("Connection", "keep-alive")
	// response, _ := client.Do(request)
	// if response.StatusCode == 200 {
	// 	body, _ := ioutil.ReadAll(response.Body)
	// 	s := string(body[:100])
	// 	s = strings.ReplaceAll(s, "\n", " ")
	// 	fmt.Printf("code: %d, body: %s...\n", resp.StatusCode, s)
	// }
}
