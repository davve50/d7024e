package main

import (
	"fmt"
	"kademlia/d7024e"
)

func main() {
	fmt.Println("Starting kademlia..")
	var node *d7024e.Network
	node.Listen("localhost", 8080)
}
