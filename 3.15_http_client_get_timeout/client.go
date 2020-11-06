package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080", nil)
	if err != nil {
		panic(err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	dump, err := httputil.DumpResponse(response, true)
	if err != nil {
		panic(err)
	}
	log.Println(string(dump))
}
