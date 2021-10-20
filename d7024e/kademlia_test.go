package d7024e

import (
	"testing"
)

func TestInitKademlia(t *testing.T) {
	contact := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001")
	kademlia := InitKademlia(contact)
	if kademlia == nil {
		t.Error("Error")
	}
}

func TestSetNetwork(t *testing.T) {
	contact := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001")
	kademlia := InitKademlia(contact)

	oldNetwork := kademlia.network

	network := Init("localhost", 8080)

	kademlia.setNetwork(network)

	if oldNetwork == kademlia.network {
		t.Error("Error")
	}
}
