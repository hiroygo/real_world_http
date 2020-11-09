package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
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
	// TODO: Client から送信を開始して、切断する
	for {
		data, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		// TrimSpace は '\n' なども削除する
		log.Println("<-", string(bytes.TrimSpace(data)))

		data = append([]byte("response "), data...)
		_, err = conn.Write(data)
		if err != nil {
			panic(err)
		}
		// TODO:
		// fmt.Fprint(conn, data) では送信できない
		// サンプルの fmt.Fprintf ならできる?
	}
}
