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
	Rpc      		string     		`json:",omitempty"`
	Id       		string     		`json:",omitempty"`
	Contacts 		[]Contact  		`json:",omitempty"`
	Target 			string 			`json:",omitempty"`
	Value 			[]byte  		`json:",omitempty"`
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
		network.SendPongPacket(&newNode)
	} else {
		fmt.Println("New contact added to bucket: " + newNode.Address)
	}
}

func (network *Network) HandleStorePacket(packet packet, contact *Contact) {
	// TOTO - Africa
	value, ok := network.kademlia.hash[packet.Target]

	if (!ok) {
		network.kademlia.hash[packet.Target] = packet.Value
		// skicka OK
	} else{
		
	}

}

func (network *Network) HandleFindNodePacket(packet packet, connection *net.UDPConn, addr *net.UDPAddr) {
	closeContacts := network.kademlia.routingtab.FindClosestContacts(NewKademliaID(packet.Id), 20) // <-- 20 ksk ska vara en global def
	newPacket, err := json.Marshal(network.CreatePacket("", nil, closeContacts, "", nil))
	log.Println(err)
	_, err = connection.WriteToUDP(newPacket, addr)
	log.Println(err)
	fmt.Println("SENT: contacts")

	/*
		Routingtable skall ändras bara av en go routine
		Tänk mer simpelt än komplext
	*/
}

func (network *Network) HandleFindValuePacket(packet packet, contact *Contact) {
	value, ok := network.kademlia.hash[packet.Id]

	if (!ok) {
		closeContacts := network.kademlia.routingtab.FindClosestContacts(NewKademliaID(packet.Id), 20)
		newPacket := network.CreatePacket("find_value", nil, closeContacts, "", nil)
		_, err := network.SendPacket(newPacket, contact.Address)
		log.Println(err)
	} else{
		newPacket := network.CreatePacket("find_value", nil, nil, "", value)
		_, err := network.SendPacket(newPacket, contact.Address)
		log.Println(err)
	}
}

// Creates an RPC packet containing sender data and possibility for contact array
func (network *Network) CreatePacket(rpc string, contact *Contact, contacts []Contact, target string, value []byte) *packet {
	createdPacket := &packet{
		Rpc:      rpc,
		Id:       contact.ID.String(),
		Contacts: contacts,
		Target: target,
		Value: value,
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
	packet := network.CreatePacket("pong", network.me, nil, "", nil)
	_, err := network.SendPacket(packet, contact.Address)
	log.Println(err)
}

func (network *Network) SendPingPacket(contact *Contact) (packet, error) {
	pack := network.CreatePacket("ping", network.me, nil, "", nil)
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
		log.Println(err)
		return newPacket, nil
	}
	return newPacket, err
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
