package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if rdump, err := httputil.DumpRequest(r, true); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	} else {
		log.Println(string(rdump))
	}

	addr, err := net.ResolveIPAddr("ip", r.FormValue("q"))
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	st := struct {
		IPAddress string `json:"ipAddress"`
	}{
		IPAddress: addr.String(),
	}
	if err := json.NewEncoder(w).Encode(st); err != nil {
		log.Println("Error:", err)
		http.Error(w, "", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
