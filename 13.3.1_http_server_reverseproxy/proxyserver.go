package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
)

func roundRobinServerAddress(ctx context.Context, beginPort, numberOfPort int) <-chan string {
	ch := make(chan string)

	go func() {
		n := 0
		for {
			select {
			case <-ctx.Done():
				close(ch)
				return
			case ch <- fmt.Sprintf(":%d", beginPort+n):
				n = (n + 1) % numberOfPort
			}
		}
	}()

	return ch
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	addrCh := roundRobinServerAddress(ctx, 8080, 3)

	// リクエストの書き換えを行うハンドラ
	director := func(req *http.Request) {
		addr, ok := <-addrCh
		if !ok {
			panic("チャネルが閉じています")
		}

		req.URL.Scheme = "http"
		req.URL.Host = addr
		req.Header.Add("X-Test", "FromProxy")
	}

	// レスポンスの書き換えを行うハンドラ
	// ハンドラを指定しないときは、実際のサーバの返答をそのまま返す
	modifier := func(resp *http.Response) error {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		newBody := bytes.NewBuffer(body)
		newBody.WriteString("<!-- passed proxy -->")
		// bytes.NewBuffer はメモリ上に存在するバッファなので閉じる必要がない。そのため Close を持たない
		// Body は io.ReadCloser なので、Close が必要。そのため、何もしない Close を追加する
		resp.Body = ioutil.NopCloser(newBody)
		resp.Header.Set("Content-Length", strconv.Itoa(newBody.Len()))
		return nil
	}

	// 通信部分もカスタマイズする場合は http.RoundTripper を使う
	rp := &httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: modifier,
	}

	log.Fatal(http.ListenAndServe(":9090", rp))
}
