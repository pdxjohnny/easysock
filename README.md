easysock
---

This go package builds upon gorilla websocket's chat example

It provides a handler function for net/http `ServeWs` that along with calling
`easysock.Hub.Run` makes all websockets connected receive each others messages

Example Usage
---

```go
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/pdxjohnny/easysock"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	go easysock.Hub.Run()
	http.HandleFunc("/ws", easysock.ServeWs)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
```
