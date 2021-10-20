package main

import (
	"bufio"
	"fmt"
	"kademlia/d7024e"
	"math/rand"
	"net"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	conn, _ := net.Dial("ip:icmp", "google.com")
	fmt.Println("IP address is: ", conn.LocalAddr())

	//If connection is root
	if conn.LocalAddr().String() == "172.27.0.2" {
		fmt.Println("[ROOT] Starting kademlia..")
		node := d7024e.InitRoot("00000000000000000000000000000000deadc0de", "172.27.0.2", 8080)
		time.Sleep(time.Millisecond * 5)
		go node.GetKademlia().CLI(true, bufio.NewScanner(os.Stdin))
		node.Listen()
	} else {
		time.Sleep(time.Second * 2)
		fmt.Println("Starting kademlia..")
		node := d7024e.Init(d7024e.GetLocalIP(), 8080)
		time.Sleep(time.Millisecond * 5)
		//scan := bufio.NewScanner(strings.NewReader("put 123"))
		go node.GetKademlia().CLI(false, bufio.NewScanner(os.Stdin))
		go node.JoinNetwork("00000000000000000000000000000000deadc0de", "172.27.0.2:8080")
		node.Listen()
	}
	fmt.Println("End of main")
}
