package d7024e

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
)

type Network struct {
	Socket *net.Conn
}

func Listen(ip string, port int) {
	ln, err := net.Listen("tcp", ip+":"+strconv.Itoa(port))

	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Insert handler here
//		go log.Println("Client successfully connected! :')")

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn){
	fmt.Println("Handling going on... :)")

	buffer, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		fmt.Println("Client left... :(")
		conn.Close()
		return
	}

	// Check which RPC
	fmt.Println("Client message:", string(buffer[:len(buffer)-1]))

	conn.Write([]byte("Pong\n"))

	handleConnection(conn)
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

