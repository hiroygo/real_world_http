package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
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

	// http リクエストをソケットに直接書き込む
	request, err := http.NewRequest("GET", "http://localhost:8080", nil)
	if err != nil {
		panic(err)
	}
	err = request.Write(conn)
	if err != nil {
		panic(err)
	}

	// ソケットを読み、http レスポンスとして解析する
	// チャンク形式が有効ならヘッダだけ読み込まれるはず
	reader := bufio.NewReader(conn)
	response, err := http.ReadResponse(reader, request)
	if err != nil {
		panic(err)
	}

	// response.Header.Get("Transfer-Encoding") だと動作しない
	if response.TransferEncoding[0] != "chunked" {
		panic("レスポンスヘッダに `Transfer-Encoding: chunked` が存在しません")
	}

	ChunkSize := func() (int, error) {
		b, err := reader.ReadBytes('\n')
		if err != nil {
			return 0, err
		}

		n, err := strconv.ParseInt(string(bytes.TrimSpace(b)), 16, 64)
		if err != nil {
			return 0, err
		}

		return int(n), nil
	}
	Chunk := func(size int) ([]byte, error) {
		b := make([]byte, size)
		if _, err := io.ReadFull(reader, b); err != nil {
			return nil, err
		}

		// "\r\n" を読み飛ばす
		if _, err := reader.Discard(2); err != nil {
			return nil, err
		}

		return b, nil
	}

	// チャンク形式
	// https://tools.ietf.org/html/rfc2616#section-3.6.1

	// チャンク形式データを逐次受信する
	for {
		chunkSize, err := ChunkSize()
		if err != nil {
			panic(err)
		}
		// すべてのチャンク形式データの送信が終わるとサイズ 0 が送信される
		if chunkSize == 0 {
			// TODO: 最後の改行を読むのは必要?
			// reader.ReadBytes('\n')
			break
		}

		chunk, err := Chunk(chunkSize)
		if err != nil {
			panic(err)
		}

		log.Printf("`%s`\n", string(chunk))
	}
}
