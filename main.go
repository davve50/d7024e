package main

import (
	"fmt"
	"kademlia/d7024e"
)

func main() {
	fmt.Println("Starting kademlia..")
	node := d7024e.Init("localhost", 8080)
	go node.Listen()
	node.GetKademlia().CLI()
}
