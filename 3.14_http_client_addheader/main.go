package main

import (
	"log"
	"net/http"
	"net/http/httputil"
)

func main() {
	client := &http.Client{}

	request, err := http.NewRequest("GET", "http://localhost:8080", nil)
	if err != nil {
		panic(err)
	}

	request.Header.Add("X-Test", "test")

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	dump, err := httputil.DumpResponse(response, true)
	if err != nil {
		panic(err)
	}
	log.Println(string(dump))
}
