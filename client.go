package main

import (
	//	"bufio"
	"fmt"
	"log"
	"net"
)

func main2() {

	raddr, _ := net.ResolveUDPAddr("udp", "localhost:8080")

	for true {
		fmt.Println("Sending ping")
		conn, err := net.DialUDP("udp", nil, raddr)
		if err != nil {
			log.Fatal(err)
		}
		conn.Write([]byte("ping"))
		conn.Close()

		conn, _ = net.ListenUDP("udp", raddr)
		buffer2 := make([]byte, 1024)
		_, _, _ = conn.ReadFromUDP(buffer2)

		fmt.Println(string(buffer2))
		conn.Close()
	}

	//	fmt.Println("Server message:", string(buffer[:len(buffer)-1]))
	//	conn.Close()
}
