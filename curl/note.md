# curl コマンド
## 独自ヘッダを追加する
* **-H "X-TEST: Hello"**
* `curl -H "X-Test: Hello" http://localhost:8080`

## メソッドを指定する
* **-X POST**
* `curl -X POST http://localhost:8080`

## POST でデータを送信する
* **-d**
* **-d** はデフォルトだと `Content-Type: application/x-www-form-urlencoded` になってしまう
* `curl -d "{\"hello\": \"world\"}" -H "Content-Type: application/json" http://localhost:8080`

## base64 を送る
* **-d @-** で標準入力を受け取る
* `base64 ./dog.jpg | curl -d @- http://localhost:8080`

## form タグでファイル送信する動き
* `curl -v --http1.0 -F "attachment-file=@test.txt;filename=sample.txt" http://localhost:8080`
```
POST / HTTP/1.0
Host: localhost:8080
Connection: close
Accept: */*
Content-Length: 213
Content-Type: multipart/form-data; boundary=------------------------fe4f11f4cc056e8d
User-Agent: curl/7.58.0

--------------------------fe4f11f4cc056e8d
Content-Disposition: form-data; name="attachment-file"; filename="sample.txt"
Content-Type: text/plain

test test aaa

--------------------------fe4f11f4cc056e8d--
```

## クッキーを送る
* `curl -c cookie.txt -b cookie.txt -b "name=value" localhost:8080`
