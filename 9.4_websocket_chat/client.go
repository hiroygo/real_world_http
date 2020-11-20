package main

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	wrTimeout = 10 * time.Second

	rdTimeout = 30 * time.Second

	// Must be less than rdTimeout.
	pingInterval = (rdTimeout * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type client struct {
	room *chatRoom
	conn *websocket.Conn
	send chan []byte
}

func (c *client) connectionMessageToChatRoom() {
	defer func() {
		c.room.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	// pong が来たら、受信のタイムアウト時刻を更新する
	c.conn.SetPongHandler(func(string) error { return c.conn.SetReadDeadline(time.Now().Add(rdTimeout)) })

	if err := c.conn.SetReadDeadline(time.Now().Add(rdTimeout)); err != nil {
		log.Println(err)
		return
	}

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println(err)
			}
			log.Printf("websocket close: %v\n", c.conn.UnderlyingConn().RemoteAddr())
			return
		}
		msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
		c.room.broadcast <- broadcastMsg{src: c, msg: msg}
	}
}

func (c *client) chatRoomMessageToConnection() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				log.Printf("%v: チャネルが閉じられているため websocket を切断します", c.conn.UnderlyingConn().RemoteAddr())
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.SetWriteDeadline(time.Now().Add(wrTimeout)); err != nil {
				log.Println(err)
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Println(err)
				return
			}
			if _, err := w.Write(msg); err != nil {
				log.Println(err)
				return
			}

			// len でバッファ内の値の数を取得する
			remainingMsg := len(c.send)
			for i := 0; i < remainingMsg; i++ {
				if _, err := w.Write(newline); err != nil {
					log.Println(err)
					return
				}
				if _, err := w.Write(<-c.send); err != nil {
					log.Println(err)
					return
				}
			}

			if err := w.Close(); err != nil {
				log.Println(err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(wrTimeout))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println(err)
				return
			}
		}
	}
}
