package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
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

func main() {
	httpServer := &http.Server{Addr: ":8080", ConnState: changeServerState}
	http.Handle("/", http.FileServer(http.Dir(".")))
	fmt.Println(httpServer.ListenAndServe())
}
