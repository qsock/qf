package main

import (
	"fmt"
	"github.com/qsock/qf/net/ws"
	"net/http"
)

func main() {
	s := ws.New()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if err := s.HandleRequest(w, r); err != nil {
			panic(err)
		}
	})
	s.HandleConnect(func(s *ws.Session) {
		fmt.Println(s.GetId())
	})
	s.HandleMessage(func(s *ws.Session, b []byte) {
		fmt.Println("recv", string(b), s.GetId())
		s.Write([]byte("nonono"))
	})
	s.HandleDisconnect(func(s *ws.Session) {
		fmt.Println("dis", s.GetId())
	})
	s.HandleClose(func(s *ws.Session, n int, str string) {
		fmt.Println("close", s.GetId(), n, str)
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}
