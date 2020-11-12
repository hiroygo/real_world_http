package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

func changeServerState(conn net.Conn, state http.ConnState) {
	log.Println(state, conn.RemoteAddr())
}

func handler(w http.ResponseWriter, r *http.Request) {
	//	dump, err := httputil.DumpRequest(r, true)
	//	if err != nil {
	//		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
	//		return
	//	}
	//	fmt.Println(string(dump))
	fmt.Fprintf(w, "<html><body>%s</body></html>\n", time.Now())
}

func main() {
	httpServer := &http.Server{Addr: ":8080", ConnState: changeServerState}
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/date.html", handler)
	fmt.Println(httpServer.ListenAndServe())
}
