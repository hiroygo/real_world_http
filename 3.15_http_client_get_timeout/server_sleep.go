package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	fmt.Println(string(dump))
	fmt.Fprintf(w, "<html><body>hello</body></html>\n")
	time.Sleep(time.Second * 10)
}

func main() {
	var httpServer http.Server
	http.HandleFunc("/", handler)
	httpServer.Addr = ":8080"
	fmt.Println(httpServer.ListenAndServe())
}
