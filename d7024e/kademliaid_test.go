package d7024e

import (
	"testing"
)

func TestNewKademliaID(t *testing.T) {
	newID := NewKademliaID("ffffffff00000000000000000000000000000000")

	if newID.String() != "ffffffff00000000000000000000000000000000" {
		t.Error("Error")
	}
}

func TestEquals(t *testing.T) {
	idSame := NewKademliaID("0000000000000000000000000000000000000001")
	idSame2 := NewKademliaID("0000000000000000000000000000000000000001")
	idDiff := NewKademliaID("1000000000000000000000000000000000000000")

	if !idSame.Equals(idSame2) {
		t.Error("Error")
	}

	if idDiff.Equals(idSame) {
		t.Error("Error")
	}

}

func TestLess(t *testing.T) {
	id := NewKademliaID("0000000000000000000000000000000000000000")
	id2 := NewKademliaID("ffffffffffffffffffffffffffffffffffffffff")
	id3 := NewKademliaID("ffffffffffffffffffffffffffffffffffffffff")

	if id2.Less(id) {
		t.Error("Error")
	}

	if id2.Less(id3) {
		t.Error("Error")
	}
}

func TestNewRandomKademliaID(t *testing.T) {
	newRandomID := NewRandomKademliaID()
	newRandomID2 := NewRandomKademliaID()

	if newRandomID.Equals(newRandomID2) {
		t.Error("Error")
	}
}
