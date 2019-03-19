package main

import (
	"net/http"
	"time"
)

func sendPhoto(imgToSend *[]byte) error {
	return nil
}

func checkForAPI() {
	for {
		resp, err := http.Get("https://api.telegram.org")
		if err != nil {
			time.Sleep(time.Minute * 2)
			continue
		}
		resp.Body.Close()
		if resp.StatusCode != 200 {
			time.Sleep(time.Minute * 2)
			continue
		}
		return
	}
}
