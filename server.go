package main

import (
	"bufio"
	"container/list"
	"flag"
	"io"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

var (
	port  = flag.String("p", "8080", "port to run the chat server")
	rooms = make(map[string]*list.List)
)

func addToRoom(room string, ws *websocket.Conn) *list.Element {
	cons, ok := rooms[room]
	if !ok {
		cons = list.New()
		rooms[room] = cons
	}
	return cons.PushBack(ws)
}

func countRoom(room string) int {
	cons, ok := rooms[room]
	if !ok {
		return 0
	}
	return cons.Len()
}

func broadcastRoom(room string, message []byte, sender *websocket.Conn) {
	cons, ok := rooms[room]
	if !ok {
		return
	}
	for e := cons.Front(); e != nil; e = e.Next() {
		ws := e.Value.(*websocket.Conn)
		if ws != sender {
			ws.Write([]byte(message))
		}
	}
}

func removeFromRoom(room string, e *list.Element) {
	cons, ok := rooms[room]
	if !ok {
		return
	}
	cons.Remove(e)
}

func chatroomHandler(ws *websocket.Conn) {
	roomName := ws.Request().URL.Path

	e := addToRoom(roomName, ws)

	log.Println("chatroom:", ws.Request().RemoteAddr, "has joined", roomName, "(", countRoom(roomName), ")")

	scanner := bufio.NewScanner(ws)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF {
			return 0, nil, io.EOF
		} else {
			return len(data), data, nil
		}
	})

	for scanner.Scan() {
		go broadcastRoom(roomName, scanner.Bytes(), ws)
	}

	removeFromRoom(roomName, e)
	log.Println("chatroom:", ws.Request().RemoteAddr, "has left", roomName, "(", countRoom(roomName), ")")
}

func main() {
	flag.Parse()
	chatroom()
}

func chatroom() {
	http.Handle("/", websocket.Handler(chatroomHandler))
	log.Println("chatroom:", "server's running at", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
