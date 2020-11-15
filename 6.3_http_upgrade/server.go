package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
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
	for i := 0; i < 5; i++ {
		send := fmt.Sprintf("hello%d\n", i)
		if _, err := fmt.Fprint(readWriter, send); err != nil {
			panic(err)
		}
		// type Writer は必ず Flush() すること
		readWriter.Flush()
		log.Printf("send: %s\n", strings.TrimSpace(send))

		recv, err := readWriter.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		log.Printf("recv: %s\n", string(bytes.TrimSpace(recv)))

		time.Sleep(time.Second)
	}
}

func main() {
	httpServer := &http.Server{Addr: ":8080", ConnState: changeServerState}
	http.HandleFunc("/", handlerUpgrade)
	fmt.Println(httpServer.ListenAndServe())
}
