package main

import (
	"fmt"
	"kademlia/d7024e"
)

func main() {
	fmt.Println("Kademlia started..")
	d7024e.Listen("localhost", 8080)
}
