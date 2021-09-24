package d7024e

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type Network struct {
	me       *Contact
	kademlia *Kademlia
}

type packet struct {
	Rpc      string    `json:",omitempty"`
	Id       string    `json:",omitempty"`
	Contacts []Contact `json:",omitempty"`
}

func (network *Network) Listen(ip string, port int) {
	time.Sleep(time.Millisecond * 5)
	pc, err1 := net.ResolveUDPAddr("udp", ":8080") // TODO: Port + ip
	connection, err2 := net.ListenUDP("udp", pc)
	if (err1 != nil) || (err2 != nil) {
		fmt.Println("Error :'(")
	}
	defer connection.Close()
	packet := packet{}
	for {
		buffer := make([]byte, 1024)
		n, addr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal(buffer[:n], &packet)
		go network.HandleRPC(packet, connection, addr)
	}
}

// Checking which rpc recieved
func (network *Network) HandleRPC(packet packet, connection *net.UDPConn, addr *net.UDPAddr) {
	switch packet.Rpc {
	case "ping":
		//do ping stuff
		network.HandlePingPongPacket(packet, connection, addr)
	case "store":
		//do store stuff
		network.HandleStorePacket(packet, connection, addr)
	case "find_node":
		//do find_node stuff
		network.HandleFindNodePacket(packet, connection, addr)
	case "find_value":
		//do find_value stuff
		network.HandleFindValuePacket(packet, connection, addr)
	default:
		//if all else fails, then something is wrong
		fmt.Println("'" + packet.Rpc + "' is not an valid RPC!")
	}
}

func (network *Network) HandlePingPongPacket(packet packet, connection *net.UDPConn, addr *net.UDPAddr) {
	newNode := NewContact(NewKademliaID(packet.Id), addr.IP.String())
	network.kademlia.routingtab.AddContact(newNode)

	if strings.ToLower(packet.Rpc) == "ping" {
		network.SendPongPacket(connection, addr, network.me)
	} else {
		fmt.Println("New contact added to bucket: " + newNode.Address)
	}
}

func (network *Network) HandleStorePacket(packet packet, connection *net.UDPConn, addr *net.UDPAddr) {
	// TOTO - Africa
}

func (network *Network) HandleFindNodePacket(packet packet, connection *net.UDPConn, addr *net.UDPAddr) {
	closestContacts := network.kademlia.routingtab.FindClosestContacts(NewKademliaID(packet.Id), 20) // <-- 20 ksk ska vara en global def
	newPacket, err := json.Marshal(network.CreatePacket("", nil, closestContacts))
	log.Println(err)
	_, err = connection.WriteToUDP(newPacket, addr)
	log.Println(err)
	fmt.Println("SENT: contacts")

	/*
		Routingtable skall ändras bara av en go routine
		Tänk mer simpelt än komplext
	*/
}

func (network *Network) HandleFindValuePacket(packet packet, connection *net.UDPConn, addr *net.UDPAddr) {
	// TODO
}

// Creates an RPC packet containing sender data and possibility for contact array
func (network *Network) CreatePacket(rpc string, contact *Contact, contacts []Contact) *packet {
	createdPacket := &packet{
		Rpc:      rpc,
		Id:       contact.ID.String(),
		Contacts: contacts,
	}
	return createdPacket
}

func (network *Network) SendPongPacket(connection *net.UDPConn, addr *net.UDPAddr, me *Contact) {
	packet, err := json.Marshal(network.CreatePacket("pong", me, nil))
	log.Println(err)
	_, err = connection.WriteToUDP(packet, addr)
	log.Println(err)
	fmt.Println("SENT: pong ", me)
}

func (network *Network) SendPingPacket(contact *Contact) {
	// TODO
}

func (network *Network) SendFindNodePacket(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataPacket(hash string) {
	// TODO
}

func (network *Network) SendStorePacket(data []byte) {
	// TODO
}
