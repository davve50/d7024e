package main

import (
	"log"
	"net"
)

func main3() {
	ln, err := net.Listen("tcp", "localhost:8080")

	if err != nil {
		log.Fatal(err)
	}
	for {
		_, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go log.Println("Client successfully connected! :')")
	}
}
