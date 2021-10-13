package main

import (
	"bufio"
	"fmt"
	"kademlia/d7024e"
	"os"
	"time"
)

func main() {
	fmt.Println("Starting kademlia..")
	node := d7024e.Init("localhost", 8080)
	go node.Listen()
	time.Sleep(time.Millisecond * 5)
	node.GetKademlia().CLI(bufio.NewScanner(os.Stdin))
}
