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

func TestISort(t *testing.T) {
	contact4 := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "localhost:8001")
	contact2 := NewContact(NewKademliaID("0000000000000000000000000001000000000000"), "localhost:8001")
	contact3 := NewContact(NewKademliaID("1000000000000000000000000000000000000000"), "localhost:8001")
	contact := NewContact(NewKademliaID("0000000000100000000000000000000000000000"), "localhost:8001")
	target := NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "localhost:8001")

	unorderedContacts := []Contact{contact, contact2, contact3, contact4}
	orderedContacts := []Contact{contact, contact2, contact3, contact4}
	iSort(orderedContacts, &target)

	/*
		for i := 0; i < len(unorderedContacts); i++ {
			fmt.Println("[TestISort] UO:", unorderedContacts[i].ID, " O:", orderedContacts[i].ID)
		}
		fmt.Println("[TestISort] Target:", target.ID)
	*/

	if unorderedContacts[0] == orderedContacts[0] {
		t.Error("Error")
	}
}
