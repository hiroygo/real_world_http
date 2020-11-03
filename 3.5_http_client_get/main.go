package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func main() {
	// 以下はサーバ側のログ
	/*
		GET /?query=hello+world%21&query=dog&test_param=test HTTP/1.1
		Host: localhost:8080
		Accept-Encoding: gzip
		User-Agent: Go-http-client/1.1
	*/
	values := url.Values{
		"query":      {"hello world!", "dog"},
		"test_param": {"test"},
	}

	// '?' からパラメータを開始する
	resp, err := http.Get("http://localhost:8080" + "?" + values.Encode())
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	log.Println("StatusCode:", resp.StatusCode)
	log.Println("Header:", resp.Header)
	log.Println(string(body))
}
