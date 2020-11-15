// https://github.com/gorilla/websocket/blob/master/examples/echo/server.go
// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

func changeServerState(conn net.Conn, state http.ConnState) {
	log.Println(state, conn.RemoteAddr())
}

var upgrader websocket.Upgrader // use default options
func handlerEcho(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			// panic にはしない
			// 接続が切断していると error が返るため
			// 本当は WriteMessage(websocket.PingMessage, nil) で接続状態を確認するみたい?
			log.Print(err)
			break
		}
		fmt.Printf("recv: `%s`\n", message)

		if err := c.WriteMessage(mt, message); err != nil {
			panic(err)
		}
		fmt.Printf("send: `%s`\n", message)
	}
}

func handlerHome(w http.ResponseWriter, r *http.Request) {
	// ここで渡した文字列が JavaScript で websocket の接続を開くときに使われる
	// e.g. ws = new WebSocket("ws://localhost");
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
	httpServer := &http.Server{Addr: ":8080", ConnState: changeServerState}
	http.HandleFunc("/", handlerHome)
	http.HandleFunc("/echo", handlerEcho)
	log.Println(httpServer.ListenAndServe())
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
