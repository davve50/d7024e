package d7024e

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

func TestKademlia_CLI(t *testing.T) {
	fmt.Println("[Test CLI]")
	node := InitRoot("00000000000000000000000000000000deadc0de", "localhost", 8080)
	go node.Listen()

	node.kademlia.network.kademlia.CLI(true, bufio.NewScanner(strings.NewReader("ez\nget 00000000000000000000000000000000deadc0de\nget 2\nexit")))

}
