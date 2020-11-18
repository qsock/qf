package main

import (
	"fmt"
	"github.com/gorilla/websocket"
)

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:8080/ws", nil)
	if err != nil {
		panic(err)
	}

	if err := conn.WriteMessage(websocket.TextMessage, []byte("hello where")); err != nil {
		panic(err)
	}
	{
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				panic(err)
			}
			fmt.Println(msg)
		}
	}
}
