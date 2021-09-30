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

var k int = 20

/* RPC 			= Command
 * SourceID		= Id of source node
 * SourceIP 	= Adress of source
 * TargetID		= Adress/file to find/store
 * Contacts 	= List of contacts used by findnode/findvalue
 * Value		= Resulting value of a call in []byte's
 */

type packet struct {
	RPC      string    `json:",omitempty"` // ping, find_value
	SourceID string    `json:",omitempty"` // NODE_123, NODE_2342
	SourceIP string    `json:",omitempty"` // 192.168.1.1, 127.0.0.1
	TargetID string    `json:",omitempty"` // NODE_235235, NODE_76457
	Contacts []Contact `json:",omitempty"` // []Contacts of closest
	Value    []byte    `json:",omitempty"` // 123, 4323452, from hashmap
}

func (network *Network) Listen(ip string, port int) {
	time.Sleep(time.Millisecond * 5)
	adress := fmt.Sprintf("%s:%d", ip, port)
	pc, err1 := net.ResolveUDPAddr("udp", adress)
	fmt.Println("Kademlia started on adress: " + adress)
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
		packet.SourceIP = connection.RemoteAddr().String()
		go network.HandleRPC(packet, connection, addr)
	}
}

// Checking which rpc recieved
func (network *Network) HandleRPC(packet packet, connection *net.UDPConn, addr *net.UDPAddr) {
	switch packet.RPC {
	case "ping":
		//do ping stuff
		network.HandlePingPongPacket(packet)
	case "store":
		//do store stuff
		network.HandleStorePacket(packet)
	case "find_node":
		//do find_node stuff
		network.HandleFindNodePacket(packet)
	case "find_value":
		//do find_value stuff
		network.HandleFindValuePacket(packet)
	default:
		//if all else fails, then something is wrong
		fmt.Println("'" + packet.RPC + "' is not an valid RPC!")
	}
}

func (network *Network) HandlePingPongPacket(packet packet) { //, connection *net.UDPConn, addr *net.UDPAddr) {
	newNode := NewContact(NewKademliaID(packet.SourceID), packet.SourceIP)
	network.kademlia.routingtab.AddContact(newNode)

	if strings.ToLower(packet.RPC) == "ping" {
		network.SendPongPacket(&newNode)
	} else {
		fmt.Println("New contact added to bucket: " + newNode.Address)
	}
}

func (network *Network) HandleStorePacket(packet packet) {
	// TOTO - Africa
	_, ok := network.kademlia.hash[packet.TargetID]

	if !ok {
		network.kademlia.hash[packet.TargetID] = packet.Value

		newPacket := network.CreatePacket("", packet.SourceID, "", nil, []byte("SUCCESS"))
		_, err := network.SendPacket(newPacket, packet.SourceIP)
		log.Println(err)
	} else {
		newPacket := network.CreatePacket("", packet.SourceID, "", nil, []byte("FAIL"))
		_, err := network.SendPacket(newPacket, packet.SourceIP)
		log.Println(err)
	}
}

func (network *Network) HandleFindNodePacket(packet packet) {
	closeContacts := network.kademlia.routingtab.FindClosestContacts(NewKademliaID(packet.TargetID), k) // <-- 20 ksk ska vara en global def

	newPacket := network.CreatePacket("", packet.SourceID, "", closeContacts, nil)
	_, err := network.SendPacket(newPacket, packet.SourceIP)
	log.Println(err)
}

func (network *Network) HandleFindValuePacket(packet packet) {
	value, ok := network.kademlia.hash[packet.TargetID]

	if !ok {
		closeContacts := network.kademlia.routingtab.FindClosestContacts(NewKademliaID(packet.TargetID), k)
		newPacket := network.CreatePacket("find_value", packet.SourceID, "", closeContacts, nil)
		_, err := network.SendPacket(newPacket, packet.SourceIP)
		log.Println(err)
	} else {
		newPacket := network.CreatePacket("find_value", packet.SourceID, "", nil, value)
		_, err := network.SendPacket(newPacket, packet.SourceIP)
		log.Println(err)
	}
}

// Creates an RPC packet containing sender data and possibility for contact array
func (network *Network) CreatePacket(rpc string, sourceid string, targetid string, contacts []Contact, value []byte) *packet {
	createdPacket := &packet{
		RPC:      rpc,
		SourceID: sourceid,
		TargetID: targetid,
		Contacts: contacts,
		Value:    value,
	}
	return createdPacket
}

func (network *Network) SendPacket(packet *packet, addr string) (*net.UDPConn, error) {
	remoteAddress, err := net.ResolveUDPAddr("udp", addr)
	connection, err := net.DialUDP("udp", nil, remoteAddress)
	log.Println(err)
	defer connection.Close()
	marshalledPacket, err := json.Marshal(packet)
	_, err = connection.Write(marshalledPacket)
	log.Println(err)
	fmt.Println("SENT: " + string(marshalledPacket))
	return connection, err
}

func (network *Network) SendPongPacket(contact *Contact) {
	packet := network.CreatePacket("pong", network.me.ID.String(), "", nil, nil)
	_, err := network.SendPacket(packet, contact.Address)
	log.Println(err)
}

func (network *Network) SendPingPacket(contact *Contact) (packet, error) {
	pack := network.CreatePacket("ping", network.me.ID.String(), "", nil, nil)
	connection, err := network.SendPacket(pack, contact.Address)
	newPacket := packet{}
	responsePacket := make([]byte, 1024)
	connection.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
	for {
		size, err := connection.Read(responsePacket)
		if err != nil {
			log.Fatal(err)
			break
		}
		err = json.Unmarshal(responsePacket[:size], &newPacket)
		newPacket.SourceIP = connection.RemoteAddr().String()
		log.Println(err)
		return newPacket, nil
	}
	return newPacket, err
}

func (network *Network) SendFindNodePacket(contact *Contact) {
	// TODO
	/*
		1. Send packet <--- Listen() => source
		2. Wait for response <--- =>>> source
		3. React to response <--- createPacket \n packet.Source = Network.me
	*/
}

func (network *Network) SendFindDataPacket(hash string) {
	// TODO
}

func (network *Network) SendStorePacket(data []byte) {
	// TODO
}
