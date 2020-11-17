// 以下を参考に作成
// https://github.com/gorilla/websocket/tree/master/examples/chat

package main

import (
	"flag"
	"log"
	"net/http"
)

func handlerHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {
	flag.Parse()
	addr := flag.String("addr", ":8080", "http service address")

	hub := newHub()
	go hub.run()

	http.HandleFunc("/", handlerHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	log.Fatal(http.ListenAndServe(*addr, nil))
}
