// 以下を参考に作成
// https://github.com/gorilla/websocket/tree/master/examples/chat

package main

import (
	"flag"
	"log"
	"net"
	"net/http"
)

func changeServerState(conn net.Conn, state http.ConnState) {
	log.Println(state, conn.RemoteAddr())
}

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "index.html")
}

func handlerJoinChat(c *ChatRoom, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := &Client{chatRoom: c, conn: conn, send: make(chan []byte, 256)}
	client.chatRoom.register <- client

	go client.ConnectionMessageToChatRoom()
	go client.ChatRoomMessageToConnection()
}

func main() {
	flag.Parse()
	addr := flag.String("addr", ":8080", "http service address")

	cr := newChatRoom()
	go cr.Run()

	httpServer := &http.Server{Addr: *addr, ConnState: changeServerState}
	http.HandleFunc("/", handlerIndex)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handlerJoinChat(cr, w, r)
	})
	log.Fatal(httpServer.ListenAndServe())
}
