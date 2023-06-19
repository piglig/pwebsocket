package main

import "pwebsocket/server"

func main() {
	s := server.New(":8000")
	s.Listen()
}
