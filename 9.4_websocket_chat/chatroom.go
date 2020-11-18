package main

type ChatRoom struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newChatRoom() *ChatRoom {
	return &ChatRoom{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (r *ChatRoom) Run() {
	for {
		select {
		case c := <-r.register:
			r.clients[c] = true
		case c := <-r.unregister:
			if _, ok := r.clients[c]; ok {
				delete(r.clients, c)
				close(c.send)
			}
		case msg := <-r.broadcast:
			for c := range r.clients {
				select {
				case c.send <- msg:
				// チャネルのバッファに空きが無い場合は default が実行される
				default:
					close(c.send)
					delete(r.clients, c)
				}
			}
		}
	}
}
