package d7024e

import (
	"net"
	"testing"
)

func TestInit(t *testing.T) {
	network := Init("localhost", 2020)

	if network.me.Address != "localhost:2020" {
		t.Error("Error")
	}
}

func TestListen(t *testing.T) {
	network := Init("localhost", 2020)
	go network.Listen()

	packet := network.CreatePacket("_", "localhost", "localhost", nil, nil)
	network.SendPacket(packet, "127.0.0.1:8080")
}

func TestFindAllNodes(t *testing.T) {
	network := Init("localhost", 2020)
	contact := NewContact(NewKademliaID("00000000000000000000000000000000deadc0de"), "localhost:8002")
	network.kademlia.routingtab.AddContact(contact)
	_ = network.FindAllNodes(&contact)
}

func TestGetKademlia(t *testing.T) {
	network := Init("localhost", 2020)
	if network.kademlia != network.GetKademlia() {
		t.Error("Error")
	}
}

func TestHandleRPC(t *testing.T) {
	network := Init("localhost", 2020)

	// Default:
	packet := network.CreatePacket("_", "localhost", "localhost", nil, nil)
	_, _ = net.ResolveUDPAddr("udp", network.me.Address)
	network.HandleRPC(*packet, nil, nil)

	// Ping:
	packet = network.CreatePacket("ping", network.me.ID.String(), "000000000000000000010", nil, nil)
	network.HandleRPC(*packet, nil, nil)
}
