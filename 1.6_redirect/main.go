package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func handleDog(w http.ResponseWriter, r *http.Request) {
	d, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(os.Stderr, string(d))
	fmt.Fprintf(w, "dog\n")
}

func main() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/dog", handleDog)

	// StatusMovedPermanently = 301 では大抵のクライアントは
	// メソッドを GET にしてリダイレクトする
	http.Handle("/301", http.RedirectHandler("/dog", http.StatusMovedPermanently))

	// StatusPermanentRedirect = 308 では大抵のクライアントは
	// メソッドを引き継いでくれる
	// つまり "/308" に POST したら "/dog" にも POST してくれる
	http.Handle("/308", http.RedirectHandler("/dog", http.StatusPermanentRedirect))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
