package main

import (
	"flag"
	"io"
	"log"
	"os"

	"golang.org/x/net/websocket"
)

func main() {
	flag.Parse()

	origin := "http://localhost/"
	url := flag.Arg(0)

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	go io.Copy(os.Stdout, ws)
	io.Copy(ws, os.Stdin)
}
