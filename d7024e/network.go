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

func Init(ip string, port int) *Network {
	me := NewContact(NewRandomKademliaID(), fmt.Sprintf("%s:%d", ip, port))
	newNetwork := &Network{
		me:       &me,
		kademlia: InitKademlia(me),
	}
	return newNetwork
}

func (network *Network) JoinNetwork(ip string) {
	networkContact := NewContact(nil, ip)
	_, err := network.SendPingPacket(&networkContact)
	log.Println(err)
	// Use FindAllNodes and iterate through the results and populate your bucket??
}

func (network *Network) FindAllNodes(target *Contact) []Contact {
	contacts := make([]Contact, 0)
	network.kademlia.LookupContact(network.me, *target, &contacts)
	return contacts
}

func (network *Network) Listen() {
	time.Sleep(time.Millisecond * 5)
	pc, err1 := net.ResolveUDPAddr("udp", network.me.Address)
	fmt.Println("Kademlia started on adress: " + network.me.Address)
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
	_, ok := network.kademlia.hash[packet.TargetID]

	if !ok {
		network.kademlia.hash[packet.TargetID] = packet.Value

		createdPacket := network.CreatePacket("", packet.SourceID, "", nil, []byte("SUCCESS"))
		_, err := network.SendPacket(createdPacket, packet.SourceIP)
		log.Println(err)
	} else {
		createdPacket := network.CreatePacket("", packet.SourceID, "", nil, []byte("FAIL"))
		_, err := network.SendPacket(createdPacket, packet.SourceIP)
		log.Println(err)
	}
}

func (network *Network) HandleFindNodePacket(packet packet) {
	closeContacts := network.kademlia.routingtab.FindClosestContacts(NewKademliaID(packet.TargetID), k) // <-- 20 ksk ska vara en global def

	createdPacket := network.CreatePacket("", packet.SourceID, "", closeContacts, nil)
	_, err := network.SendPacket(createdPacket, packet.SourceIP)
	log.Println(err)
}

func (network *Network) HandleFindValuePacket(packet packet) {
	value, ok := network.kademlia.hash[packet.TargetID]

	if !ok {
		closeContacts := network.kademlia.routingtab.FindClosestContacts(NewKademliaID(packet.TargetID), k)
		createdPacket := network.CreatePacket("find_value", packet.SourceID, "", closeContacts, nil)
		_, err := network.SendPacket(createdPacket, packet.SourceIP)
		log.Println(err)
	} else {
		createdPacket := network.CreatePacket("find_value", packet.SourceID, "", nil, value)
		_, err := network.SendPacket(createdPacket, packet.SourceIP)
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
	createdPacket := network.CreatePacket("pong", network.me.ID.String(), "", nil, nil)
	_, err := network.SendPacket(createdPacket, contact.Address)
	log.Println(err)
}

func (network *Network) SendPingPacket(contact *Contact) (packet, error) {
	createdPacket := network.CreatePacket("ping", network.me.ID.String(), "", nil, nil)
	connection, err := network.SendPacket(createdPacket, contact.Address)
	newPacket := packet{}
	responsePacket := make([]byte, 1024)
	connection.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
	for {
		size, err := connection.Read(responsePacket)
		// Fix error handler :)))  FIX THIS FIX THIS FIX THIS FIX THIS FIX THIS FIX THIS FIX THIS FIX THIS FIX THIS FIX THIS FIX THIS FIX THIS FIX THIS FIX THIS
		if err != nil {
			log.Fatal(err)
			break
		}
		err = json.Unmarshal(responsePacket[:size], &newPacket)
		newPacket.SourceIP = connection.RemoteAddr().String()
		newNode := NewContact(NewKademliaID(newPacket.SourceID), newPacket.SourceIP)
		network.kademlia.routingtab.AddContact(newNode)
		log.Println(err)
		return newPacket, nil
	}
	return newPacket, err
}

func (network *Network) SendFindNodePacket(contact *Contact, found chan []Contact) {
	var contacts []Contact
	createdPacket := network.CreatePacket("find_node", network.me.ID.String(), contact.ID.String(), nil, nil)
	connection, err := network.SendPacket(createdPacket, contact.Address)
	log.Println(err)
	newPacket := packet{}
	responsePacket := make([]byte, 1024)
	connection.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
	for {
		length, err := connection.Read(responsePacket)
		log.Println(err)
		err = json.Unmarshal(responsePacket[:length], &newPacket)
	}
	found <- contacts // Channel here maybe???
}

func (network *Network) SendFindValuePacket(contact *Contact, hash string) {
	createdPacket := network.CreatePacket("find_value", network.me.ID.String(), "", nil, []byte(hash))
	connection, err := network.SendPacket(createdPacket, contact.Address)
	log.Println(err)
	newPacket := packet{}
	responsePacket := make([]byte, 1024)
	connection.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
	length, err := connection.Read(responsePacket)
	err = json.Unmarshal(responsePacket[:length], &newPacket)

	// We need to decide what to do after recieving the answer???

	// PROBABLY NEED SOME CHANNELS FOR THIS SO WILL PROBABLY NEED TO DO
	// KADEMLIA.GO PROGRAMMING BEFORE WE START THIS MESS
}

func (network *Network) SendStorePacket(contact *Contact, data []byte) {
	createdPacket := network.CreatePacket("store", network.me.ID.String(), "", nil, data)
	connection, err := network.SendPacket(createdPacket, contact.Address)
	log.Println(err)
	newPacket := packet{}
	responsePacket := make([]byte, 1024)
	connection.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
	length, err := connection.Read(responsePacket)
	err = json.Unmarshal(responsePacket[:length], &newPacket)
	fmt.Println(string(newPacket.Value))
}
