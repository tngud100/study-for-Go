package main

import (
	"net/http"
)

type requestresult struct {
	url    string
	status string
}

func main() {
	results := make(map[string]string)
	c := make(chan requestresult)
	urls := []string{
		"https://www.airbnb.com/",
		"https://www.google.com/",
		"https://www.naver.com/",
		"https://www.reddit.com/",
		"https://www.google.com/",
		"https://www.soundcloud.com/",
		"https://www.facebook.com/",
		"https://www.instagram.com/",
		"https://instagram.com/",
		"https://instagram.com/",
	}
	for _, url := range urls {
		go hitURL(url, c)
	}

	for i := 0; i < len(urls); i++ {
		result := <-c
		results[result.url] = result.status
	}
}

func hitURL(url string, c chan<- requestresult) {
	resp, err := http.Get(url)
	status := "OK"
	if err != nil || resp.StatusCode >= 400 {
		c <- requestresult{url: url, status: "Failed"}
	}
	c <- requestresult{url: url, status: status}
}
