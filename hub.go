// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package easysock

import (
	"log"
	"net/http"
)

type EventFunc func (Hub *WebSocketManger, Conn *Connection)

// WebSocketManger maintains the set of active Connections and broadcasts messages to the
// Connections.
type WebSocketManger struct {
	// Registered Connections.
	Connections map[*Connection]bool

	// Inbound messages from the Connections.
	Broadcast chan []byte

	// Register requests from the Connections.
	Register chan *Connection

	// Unregister requests from Connections.
	Unregister chan *Connection

	// On Connect function
	OnConnect EventFunc

	// On Close function
	OnClose EventFunc
}

var Hub = WebSocketManger{
	Broadcast:   make(chan []byte, maxMessageSize),
	Register:    make(chan *Connection, maxMessageSize),
	Unregister:  make(chan *Connection, maxMessageSize),
	Connections: make(map[*Connection]bool, maxMessageSize),
	OnConnect: nil,
	OnClose: nil,
}

func (Hub *WebSocketManger) Run() {
	for {
		select {
		case c := <-Hub.Register:
			Hub.Connections[c] = true
			if Hub.OnConnect != nil {
				go Hub.OnConnect(Hub, c)
			}
		case c := <-Hub.Unregister:
			if _, ok := Hub.Connections[c]; ok {
				delete(Hub.Connections, c)
				close(c.Send)
				if Hub.OnClose != nil {
					go Hub.OnClose(Hub, c)
				}
			}
		case m := <-Hub.Broadcast:
			for c := range Hub.Connections {
				select {
				case c.Send <- m:
				default:
					close(c.Send)
					delete(Hub.Connections, c)
				}
			}
		}
	}
}

// serverWs handles websocket requests from the peer.
func (Hub *WebSocketManger) ServeWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &Connection{Send: make(chan []byte, 256), ws: ws}
	Hub.Register <- c
	go c.writePump()
	c.readPump()
}
