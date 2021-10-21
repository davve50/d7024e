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
	network := Init("localhost", 3030)
	go network.Listen()

	packet := network.CreatePacket("_", "127.0.0.1", "localhost", "localhost", nil, nil)
	network.SendPacket(packet, "127.0.0.1:8080")
}

func TestFindAllNodes(t *testing.T) {
	/*
		network := Init("localhost", 2020)
		contact := NewContact(NewKademliaID("00000000000000000000000000000000deadc0de"), "localhost:8002")
		network.kademlia.routingtab.AddContact(contact)
		_ = network.FindAllNodes(&contact)
	*/
}

func TestGetKademlia(t *testing.T) {
	network := Init("localhost", 4040)
	if network.kademlia != network.GetKademlia() {
		t.Error("Error")
	}
}
func TestGetIp(t *testing.T) {
	if GetLocalIP() != "" {
		t.Error("Error")
	}
}

func TestHandleRPC(t *testing.T) {
	network := Init("localhost", 5050)
	go network.Listen()

	// Default:

	packet := network.CreatePacket("_", "localhost", "", "localhost", nil, nil)
	_, _ = net.ResolveUDPAddr("udp", network.me.Address)
	network.HandleRPC(*packet, nil, nil)

	// Ping:

	packet = network.CreatePacket("ping", "localhost:8080", network.me.ID.String(), "00000000000000000000000000000000deadc0de", nil, nil)
	network.HandleRPC(*packet, nil, nil)

}
