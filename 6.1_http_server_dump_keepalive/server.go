package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
)

func changeServerState(conn net.Conn, state http.ConnState) {
	var s string
	switch state {
	case http.StateNew:
		s = "StateNew"
	case http.StateActive:
		s = "StateActive"
	case http.StateIdle:
		s = "StateIdle"
	case http.StateHijacked:
		s = "StateHijacked"
	case http.StateClosed:
		s = "StateClosed"
	default:
		panic("Unknown State")
	}
	log.Println(s, conn.RemoteAddr())
}

func handler(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	fmt.Println(string(dump))
	fmt.Fprintf(w, "<html><body>hello</body></html>\n")
}

func main() {
	httpServer := &http.Server{Addr: ":8080", ConnState: changeServerState}
	http.HandleFunc("/", handler)
	fmt.Println(httpServer.ListenAndServe())
}
