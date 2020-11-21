package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if rdump, err := httputil.DumpRequest(r, true); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	} else {
		log.Println(string(rdump))
	}

	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	addrs := make(map[string][]string)
	// 複数のホスト名を解決するときはホスト名を半角スペースで区切る
	for _, name := range strings.Fields(r.FormValue("q")) {
		if ss, err := net.LookupHost(name); err == nil {
			addrs[name] = append(addrs[name], ss...)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(addrs); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
