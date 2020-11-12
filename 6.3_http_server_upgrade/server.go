package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

func changeServerState(conn net.Conn, state http.ConnState) {
	log.Println(state, conn.RemoteAddr())
}

func handlerUpgrade(w http.ResponseWriter, r *http.Request) {
	validRequest := r.Header.Get("Connection") == "Upgrade" && r.Header.Get("Upgrade") == "TestProtocol"
	if !validRequest {
		w.WriteHeader(400)
		return
	}
	log.Println("Upgrade to TestProtocol")

	// 型アサーション
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		panic("http.Hijacker へのキャストに失敗しました")
	}
	// ソケットの取得
	conn, readWriter, err := hijacker.Hijack()
	if err != nil {
		panic(err)
	}
	// TODO: あえてクローズしないでやってみる
	defer conn.Close()

	// プロトコルが変わるというレスポンスを送信する
	// 101 Switching Protocols
	response := http.Response{
		StatusCode: 101,
		Header:     make(http.Header),
	}
	response.Header.Set("Upgrade", "TestProtocol")
	response.Header.Set("Connection", "Upgrade")
	if err := response.Write(conn); err != nil {
		panic(err)
	}

	// ソケット通信の開始
	for i := 1; i <= 10; i++ {
		if _, err := fmt.Fprintf(readWriter, "%d\n", i); err != nil {
			panic(err)
		}
		// type Writer は必ず Flush() すること
		readWriter.Flush()
		log.Println("->", i)

		recv, err := readWriter.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		log.Printf("<- %s", string(recv))
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	httpServer := &http.Server{Addr: ":8080", ConnState: changeServerState}
	http.HandleFunc("/", handlerUpgrade)
	fmt.Println(httpServer.ListenAndServe())
}
