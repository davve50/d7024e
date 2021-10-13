package d7024e

import (
	"testing"
)

func TestString(t *testing.T) {
	contact := NewContact(NewKademliaID("ffffffffffffffffffffffffffffffffffffffff"), "localhost:8000")
	if contact.String() != `contact("ffffffffffffffffffffffffffffffffffffffff", "localhost:8000")` {
		t.Error("Error")
	}
}
