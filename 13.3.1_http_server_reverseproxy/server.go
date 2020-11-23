package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if d, err := httputil.DumpRequest(r, true); err != nil {
		panic(err)
	} else {
		log.Println(string(d))
	}

	fmt.Fprintf(w, "<html><body>%s</body></html>", os.Args[1])
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(os.Args[1], nil))
}
