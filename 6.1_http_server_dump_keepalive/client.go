package main

import (
	"log"
	"net/http"
	"net/http/httputil"
)

// curl や go の http.Client ではデフォルトで Keep-Alive が有効になっている
func main() {
	client := &http.Client{}

	for i := 0; i < 2; i++ {
		func() {
			request, err := http.NewRequest("GET", "http://localhost:8080", nil)
			if err != nil {
				panic(err)
			}

			response, err := client.Do(request)
			if err != nil {
				panic(err)
			}
			defer response.Body.Close()

			dump, err := httputil.DumpResponse(response, true)
			if err != nil {
				panic(err)
			}
			log.Println(string(dump))
		}()
	}

}
