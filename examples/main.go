//go:build !healthcheck
// +build !healthcheck

// Minimal example: detect HTTP requests with a vCWAF detector.
// Ví dụ tối giản: phát hiện HTTP request bằng detector của vCWAF.
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/hitechcloudvn/t1k-go"
)

func initDetect() *t1k.Server {
	server, err := t1k.NewWithPoolSizeWithTimeout(os.Getenv("DETECTOR_ADDR"), 10, 10*time.Second)
	if err != nil {
		return nil
	}

	server.UpdateSockErrorHandler(func(err error) {
		fmt.Printf("Socket error: %s", err.Error())
	})
	return server
}

func main() {
	server := initDetect()
	if server == nil {
		fmt.Println("Init detect error")
		return
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		res, err := server.DetectHttpRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if res.Blocked() {
			http.Error(w, fmt.Sprintf("blocked event id %s", res.EventID()), res.StatusCode())
			return
		}
		_, _ = w.Write([]byte("allowed"))
	})
	_ = http.ListenAndServe(":80", nil)
}
