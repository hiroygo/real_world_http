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

func handlerChunkedResponse(w http.ResponseWriter, r *http.Request) {
	// 型アサーション
	flusher, ok := w.(http.Flusher)
	if !ok {
		panic("http.Flusher へのキャストに失敗しました")
	}

	for i := 0; i < 5; i++ {
		fmt.Fprintf(w, "chunk%d\n", i)
		// ここで Flush() しないと、チャンクで送信されない
		// Flush() すると、ヘッダに `Transfer-Encoding: chunked` が付加される
		// Flush() するごとにクライアントにデータをチャンク形式で送信する
		// サーバが巨大なデータをメモリに少しづつ展開しながらクライアントにデータを送ることもできる
		flusher.Flush()
		time.Sleep(time.Second)
	}
}

func main() {
	httpServer := &http.Server{Addr: ":8080", ConnState: changeServerState}
	http.HandleFunc("/", handlerChunkedResponse)
	fmt.Println(httpServer.ListenAndServe())
}
