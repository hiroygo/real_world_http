package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	dialer := &net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 10 * time.Second,
	}

	conn, err := dialer.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// http リクエストをソケットに直接書き込む
	request, err := http.NewRequest("GET", "http://localhost:8080", nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Connection", "Upgrade")
	request.Header.Set("Upgrade", "TestProtocol")
	err = request.Write(conn)
	if err != nil {
		panic(err)
	}

	// ソケットを読み、http レスポンスとして解析する
	response, err := http.ReadResponse(reader, request)
	if err != nil {
		panic(err)
	}
	log.Println("Status: ", response.Status)
	log.Println("Headers: ", response.Header)

	// TestProtocol を送受信する
	for {
		recv, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		// TrimSpace は '\n' なども削除する
		log.Printf("recv: %s\n", string(bytes.TrimSpace(recv)))

		send := append([]byte("echo"), recv...)
		_, err = fmt.Fprint(conn, string(send))
		if err != nil {
			panic(err)
		}
		log.Printf("send: %s\n", string(bytes.TrimSpace(send)))
	}
}
