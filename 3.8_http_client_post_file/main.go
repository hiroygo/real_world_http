package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	file, err := os.Open("main.go")
	if err != nil {
		panic(err)
	}

	resp, err := http.Post("http://localhost:8080", "text/plain", file)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println("Status:", resp.Status)
}
