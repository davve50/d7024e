package d7024e

import (
	"testing"
)

func TestRoutingTable(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "localhost:8000"))

	rt.AddContact(NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "localhost:8000"))
	rt.AddContact(NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "localhost:8001"))
	rt.AddContact(NewContact(NewKademliaID("0000000000000000000000000000000000000010"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("0000000000000000000000000000000000000100"), "localhost:8003"))
	rt.AddContact(NewContact(NewKademliaID("0000000000000000000000000000000000001000"), "localhost:8004"))

	_ = rt.FindClosestContacts(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), 20)

	var expected_contacts []Contact

	expected_contacts = append(expected_contacts,
		NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "localhost:8000"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000001"), "localhost:8001"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000010"), "localhost:8002"),
		NewContact(NewKademliaID("0000000000000000000000000000000000000100"), "localhost:8003"),
		NewContact(NewKademliaID("0000000000000000000000000000000000001000"), "localhost:8004"))

	for i := 0; i < len(expected_contacts); i++ {
		contacts := rt.FindClosestContacts(NewKademliaID("0000000000000000000000000000000000000000"), i+1)
		if len(contacts) != i+1 {
			t.Error("Error")
		}

		for c := 0; c < i+1; c++ {
			if contacts[c].ID.String() != expected_contacts[c].ID.String() {
				t.Error("Error <", contacts[c].ID.String(), ":", expected_contacts[c].ID.String(), ">")
			}
		}
	}
}
