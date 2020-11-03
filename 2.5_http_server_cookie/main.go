package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

func handler(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	fmt.Println(string(dump))

	w.Header().Add("Set-Cookie", "VISIT=TRUE")
	if s, ok := r.Header["Cookie"]; ok {
		fmt.Fprintf(w, "<html><body>2回目以降: %s</body></html>\n", s)
	} else {
		fmt.Fprintf(w, "<html><body>初訪問</body></html>\n")
	}
}

// p46
// curl -c cookie.txt -b cookie.txt -b "name=value" localhost:8080
func main() {
	var httpServer http.Server
	http.HandleFunc("/", handler)
	httpServer.Addr = ":8080"
	fmt.Println(httpServer.ListenAndServe())
}
