package main

import (
	"fmt"
	"net/http"
)

func handleHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HelloWorld\n")
}

func handleTest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "test\n")
}

func main() {
	// ListenAndServe の第2引数に Handler を渡すと、この設定は動作しない
	http.HandleFunc("/", handleTest)

	// すべてのリクエストが handleHello で処理される
	// http://localhost:8080 でも HelloWorld が返る
	// http://localhost:8080/test.jpg でも HelloWorld が返る
	http.ListenAndServe(":8080", http.HandlerFunc(handleHello))
}
