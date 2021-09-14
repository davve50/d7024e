package main

import (
	"log"
	"net"
	"bufio"
	"fmt"
)

func main(){
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	conn.Write([]byte("Ping\n"))
	buffer, err := bufio.NewReader(conn).ReadBytes('\n')

	fmt.Println("Server message:", string(buffer[:len(buffer)-1]))
	conn.Close()
}

