// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package easysock

// WebSocketManger maintains the set of active connections and broadcasts messages to the
// connections.
type WebSocketManger struct {
	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection
}

var Hub = WebSocketManger{
	broadcast:   make(chan []byte, maxMessageSize),
	register:    make(chan *connection, maxMessageSize),
	unregister:  make(chan *connection, maxMessageSize),
	connections: make(map[*connection]bool, maxMessageSize),
}

func (Hub *WebSocketManger) Run() {
	for {
		select {
		case c := <-Hub.register:
			Hub.connections[c] = true
		case c := <-Hub.unregister:
			if _, ok := Hub.connections[c]; ok {
				delete(Hub.connections, c)
				close(c.send)
			}
		case m := <-Hub.broadcast:
			for c := range Hub.connections {
				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(Hub.connections, c)
				}
			}
		}
	}
}
