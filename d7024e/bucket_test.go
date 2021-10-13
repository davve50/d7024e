package d7024e

import (
	"container/list"
	"testing"
)

func Test_newBucket(t *testing.T) {
	bucket := newBucket()
	newList := list.New()

	if bucket.list.Len() != newList.Len() {
		t.Error("Error")
	}
}

func Test_addContact(t *testing.T) {
	bucket := newBucket()

	contact := NewContact(NewKademliaID("00000000000000000000000000000000deadc0de"), "localhost:8002")
	contact2 := NewContact(NewKademliaID("00000000000000000000000000000000deadbeef"), "localhost:8002")

	bucket.AddContact(contact)
	bucket.AddContact(contact2)

	if bucket.list.Front().Value != contact2 {
		t.Error("Error")
	}

	bucket.AddContact(contact)

	if bucket.list.Front().Value != contact {
		t.Error("Error")
	}
}

func Test_Len(t *testing.T) {
	bucket := newBucket()

	contact := NewContact(NewKademliaID("00000000000000000000000000000000deadc0de"), "localhost:8002")
	bucket.AddContact(contact)

	if bucket.Len() != 1 {
		t.Error("Error")
	}
}
