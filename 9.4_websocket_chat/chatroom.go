package main

import "log"

type broadcastMsg struct {
	src *client
	msg []byte
}

type chatRoom struct {
	clients    map[*client]bool
	broadcast  chan broadcastMsg
	register   chan *client
	unregister chan *client
}

func newChatRoom() *chatRoom {
	return &chatRoom{
		broadcast:  make(chan broadcastMsg),
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[*client]bool),
	}
}

func (r *chatRoom) run() {
	for {
		select {
		case c := <-r.register:
			r.clients[c] = true
			log.Printf("client registered: %v\n", c.conn.UnderlyingConn().RemoteAddr())
		case c := <-r.unregister:
			if _, ok := r.clients[c]; ok {
				delete(r.clients, c)
				close(c.send)
				log.Printf("client unregistered: %v\n", c.conn.UnderlyingConn().RemoteAddr())
			}
		case brmsg := <-r.broadcast:
			log.Printf("%s: %s\n", brmsg.src.conn.UnderlyingConn().RemoteAddr(), brmsg.msg)
			for c := range r.clients {
				select {
				case c.send <- brmsg.msg:
				// チャネルのバッファに空きが無い場合は default が実行される
				default:
					close(c.send)
					delete(r.clients, c)
				}
			}
		}
	}
}
