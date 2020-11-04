package main

import (
	"log"
	"net/http"
	"net/url"
)

func main() {
	// 以下はサーバ側のログ
	/*
		POST / HTTP/1.1
		Host: localhost:8080
		Accept-Encoding: gzip
		Content-Length: 46
		Content-Type: application/x-www-form-urlencoded
		User-Agent: Go-http-client/1.1

		query=hello+world%21&query=dog&test_param=test
	*/
	values := url.Values{
		"query":      {"hello world!", "dog"},
		"test_param": {"test"},
	}

	resp, err := http.PostForm("http://localhost:8080", values)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println("Status:", resp.Status)
}
