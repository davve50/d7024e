package d7024e

import (
	"bufio"
	"strings"
	"testing"
)

func TestKademlia_CLI(t *testing.T) {
	//contact := NewContact(NewKademliaID("00000000000000000000000000000000deadc0de"), "localhost:8002")
	network := Init("localhost", 2020)

	network.kademlia.CLI(bufio.NewScanner(strings.NewReader("ez\nput 123\nget 00000000000000000000000000000000deadc0de\nget 2\nexit")))

}
