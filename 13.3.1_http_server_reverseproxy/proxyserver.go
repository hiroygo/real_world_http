package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
)

func main() {
	// リクエストの書き換えを行うハンドラ
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = ":8080"
		req.Header.Add("X-Test", "FromProxy")
	}

	// レスポンスの書き換えを行うハンドラ
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
