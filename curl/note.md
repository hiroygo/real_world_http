# curl コマンド
## 独自ヘッダを追加する
* **-H "X-TEST: Hello"**
* `curl -H "X-Test: Hello" http://localhost:8080`

## メソッドを指定する
* **-X POST**
* `curl -X POST http://localhost:8080`

## POST でデータを送信する
* **-d** はデフォルトだと `Content-Type: application/x-www-form-urlencoded` になってしまう
  * `curl -d "{\"hello\": \"world\"}" -H "Content-Type: application/json" http://localhost:8080`

### **-d**
* **-d** は **--data** と同じ
* **--data-urlencode** との違いは必要に応じて文字をエスケープするかどうか

#### **-d** 単体だと POST の動きになる
* `curl -d "test" localhost:8080`
```
POST / HTTP/1.1
Host: localhost:8080
Accept: */*
Content-Length: 4
Content-Type: application/x-www-form-urlencoded
User-Agent: curl/7.47.0

test
```

#### **-d** と GET の組み合わせで URL 末尾にクエリを付与する
* `curl -G -d "test" localhost:8080`
```
GET /?test HTTP/1.1
Host: localhost:8080
Accept: */*
User-Agent: curl/7.47.0

```

### **--data-urlencode**
* **--data-urlencode** 単体だと POST の動きになる
* `curl --data-urlencode "hello world" localhost:8080`
```
POST / HTTP/1.1
Host: localhost:8080
Accept: */*
Content-Length: 13
Content-Type: application/x-www-form-urlencoded
User-Agent: curl/7.47.0

hello%20world
```

#### **--data-urlencode** と GET の組み合わせで URL 末尾にクエリを付与する
* `curl -G --data-urlencode "hello world" localhost:8080`
```
GET /?hello%20world HTTP/1.1
Host: localhost:8080
Accept: */*
User-Agent: curl/7.47.0

```

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
