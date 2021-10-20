package d7024e

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

type Network struct {
	me       *Contact
	kademlia *Kademlia
	findCont chan []Contact
	findVal  chan []byte
}

/* RPC 			= Command
 * SourceID		= Id of source node
 * SourceIP 	= Adress of source
 * TargetID		= Adress/file to find/store
 * Contacts 	= List of contacts used by findnode/findvalue
 * Value		= Resulting value of a call in []byte's
 */

type packet struct {
	RPC      string    `json:",omitempty"` 
	SourceID string    `json:",omitempty"` 
	SourceIP string    `json:",omitempty"` 
	TargetID string    `json:",omitempty"` 
	Contacts []Contact `json:",omitempty"` 
	Value    []byte    `json:",omitempty"` 
}

// Initialize regular node
func Init(ip string, port int) *Network {
	me := NewContact(NewRandomKademliaID(), fmt.Sprintf("%s:%d", ip, port))
	newNetwork := &Network{
		me:       &me,
		kademlia: InitKademlia(me),
		findCont: make(chan []Contact),
		findVal:  make(chan []byte),
	}

	fmt.Println("Kademlia ID: ", me.ID.String())
	newNetwork.kademlia.setNetwork(newNetwork)
	return newNetwork
}

// Initialize root node
func InitRoot(id string, ip string, port int) *Network {
	me := NewContact(NewKademliaID(id), fmt.Sprintf("%s:%d", ip, port))
	newNetwork := &Network{
		me:       &me,
		kademlia: InitKademlia(me),
		findCont: make(chan []Contact),
		findVal:  make(chan []byte),
	}

	fmt.Println("Kademlia ID: ", me.ID.String())
	newNetwork.kademlia.setNetwork(newNetwork)
	return newNetwork
}

// Ping & Pong to add node to root
func (network *Network) JoinNetwork(id string, ip string) {
	time.Sleep(time.Second * 1)
	networkContact := NewContact(NewKademliaID(id), ip)
	network.SendPingPacket(networkContact)
	if string(<-network.findVal) == "pong" {
		network.kademlia.routingtab.AddContact(networkContact)
	}

	// FIND_NODE to populate buckets
	contacts := network.kademlia.LookupContact(network.me)

	for _, cont := range contacts {
		if !cont.ID.Equals(network.me.ID) {
			network.kademlia.routingtab.AddContact(cont)
		}
	}
}

// Returns a kademlia reference
func (network *Network) GetKademlia() *Kademlia {
	return network.kademlia
}

// Pretty printing for debug
func (packet *packet) String() string {
	printString := string("\n\tRPC: " + packet.RPC +
		"\n\tSourceID: " + packet.SourceID +
		"\n\tSourceIP: " + packet.SourceIP +
		"\n\tTargetID: " + packet.TargetID)
	contactsString := "\n\tContacts: \n"
	for _, contact := range packet.Contacts {
		contactsString = contactsString + "\t\t" + contact.String() + "\n"
	}
	printString = printString + contactsString + "\tValue: " + string(packet.Value)

	return printString
}

// Listen for packet
func (network *Network) Listen() {
	pc, err1 := net.ResolveUDPAddr("udp", network.me.Address)
	fmt.Println("Kademlia is listening on address: " + network.me.Address)
	connection, err2 := net.ListenUDP("udp", pc)
	if (err1 != nil) || (err2 != nil) {
		fmt.Println("Error :'(")
	}
	defer connection.Close()
	for {
		packet := packet{}
		buffer := make([]byte, 65536)
		//fmt.Println("Listening")
		n, addr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal(buffer[:n], &packet)
		//fmt.Println("\033[32m", "RECIEVED <", addr, ">: ", packet.String(), "\033[0m")
		if packet.RPC == "stop_rpc" {
			break
		}
		newContact := NewContact(NewKademliaID(packet.SourceID), packet.SourceIP)
		if !newContact.ID.Equals(network.me.ID) {
			network.kademlia.routingtab.AddContact(newContact)
		}
		network.HandleRPC(packet, connection, addr)
	}
}

// Checking which rpc was recieved
func (network *Network) HandleRPC(packet packet, connection *net.UDPConn, addr *net.UDPAddr) {
	switch packet.RPC {
	case "ping":
		network.HandlePingPacket(packet)
	case "pong":
		network.findVal <- []byte("pong")
	case "store":
		network.HandleStorePacket(packet)
	case "store_reply":
		network.findVal <- packet.Value
	case "find_node":
		network.HandleFindNodePacket(packet)
	case "find_node_reply":
		network.findCont <- packet.Contacts
	case "find_value":
		network.HandleFindValuePacket(packet)
	case "find_value_reply":
		network.findVal <- packet.Value
		network.findCont <- packet.Contacts
	default:
		fmt.Println("'" + packet.RPC + "' is not an valid RPC!")
	}
}

func (network *Network) HandlePingPacket(packet packet) {
	newNode := NewContact(NewKademliaID(packet.SourceID), packet.SourceIP)
	if !newNode.ID.Equals(network.me.ID) {
		network.kademlia.routingtab.AddContact(newNode)
	}
	network.SendPongPacket(newNode)
}

func (network *Network) HandleStorePacket(packet packet) {
	_, ok := network.kademlia.hash[packet.SourceID]

	if !ok {
		network.kademlia.hash[packet.SourceID] = packet.Value

		createdPacket := network.CreatePacket("store_reply", network.me.Address, network.me.ID.String(), packet.SourceID, nil, []byte("SUCCESS"))
		network.SendPacket(createdPacket, packet.SourceIP)
	} else {
		createdPacket := network.CreatePacket("store_reply", network.me.Address, network.me.ID.String(), packet.SourceID, nil, []byte("FAIL"))
		network.SendPacket(createdPacket, packet.SourceIP)
	}
}

func (network *Network) HandleFindNodePacket(packet packet) {
	closeContacts := network.kademlia.routingtab.FindClosestContacts(NewKademliaID(packet.TargetID), bucketSize)

	createdPacket := network.CreatePacket("find_node_reply", network.me.Address, network.me.ID.String(), packet.SourceID, closeContacts, []byte(""))
	network.SendPacket(createdPacket, packet.SourceIP)
}

func (network *Network) HandleFindValuePacket(packet packet) {
	value, ok := network.kademlia.hash[string(packet.Value)]
	//fmt.Println("\033[31m", ok)

	if !ok {
		closeContacts := network.kademlia.routingtab.FindClosestContacts(NewKademliaID(string(packet.Value)), bucketSize)
		createdPacket := network.CreatePacket("find_value_reply", network.me.Address, network.me.ID.String(), packet.SourceID, closeContacts, []byte(""))
		network.SendPacket(createdPacket, packet.SourceIP)
	} else {
		createdPacket := network.CreatePacket("find_value_reply", network.me.Address, network.me.ID.String(), packet.SourceID, nil, value) 
		network.SendPacket(createdPacket, packet.SourceIP)
	}
}

// Creates and return a packet
func (network *Network) CreatePacket(rpc string, sourceip string, sourceid string, targetid string, contacts []Contact, value []byte) *packet {
	createdPacket := &packet{
		RPC:      rpc,
		SourceID: sourceid,
		SourceIP: sourceip,
		TargetID: targetid,
		Contacts: contacts,
		Value:    value,
	}
	return createdPacket
}


func (network *Network) SendPacket(packet *packet, addr string) {
	// Returns an address of UDP end point
	remoteAddress, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Println(err)
	}

	// Dial the UDP address
	connection, err := net.DialUDP("udp", nil, remoteAddress)
	if err != nil {
		log.Println(err)
	}
	defer connection.Close()

	marshalledPacket, err := json.Marshal(packet)
	if err != nil {
		log.Println(err)
	}

	// Write to connection
	_, err = connection.Write(marshalledPacket)
	if err != nil {
		log.Println(err)
	}
	//fmt.Println("\033[33m", "SENT <", addr, ">: ", packet.String(), "\033[0m")
}

// Create and send ping packet
func (network *Network) SendPingPacket(contact Contact) {
	createdPacket := network.CreatePacket("ping", network.me.Address, network.me.ID.String(), contact.ID.String(), nil, nil)
	network.SendPacket(createdPacket, contact.Address)
}

// Create and send pong packet
func (network *Network) SendPongPacket(contact Contact) {
	createdPacket := network.CreatePacket("pong", network.me.Address, network.me.ID.String(), contact.ID.String(), nil, nil)
	network.SendPacket(createdPacket, contact.Address)
}

// Create and send FIND_NODE packet
func (network *Network) SendFindNodePacket(contact Contact, target *Contact) {
	createdPacket := network.CreatePacket("find_node", network.me.Address, network.me.ID.String(), target.ID.String(), nil, nil)
	network.SendPacket(createdPacket, contact.Address)
}

// Create and send FIND_VALUE packet
func (network *Network) SendFindValuePacket(contact Contact, hash string) {
	createdPacket := network.CreatePacket("find_value", network.me.Address, network.me.ID.String(), contact.ID.String(), nil, []byte(hash))
	network.SendPacket(createdPacket, contact.Address)
}

// Create and send STORE packet
func (network *Network) SendStorePacket(contact Contact, data []byte) {
	createdPacket := network.CreatePacket("store", network.me.Address, network.me.ID.String(), contact.ID.String(), nil, data)
	network.SendPacket(createdPacket, contact.Address)
}

// Gotten from https://github.com/holwech/UDP-module/blob/e03eccee9bfb5585d2c27c7e153fef273285099c/communication.go#L15
func GetLocalIP() string {
	var localIP string
	addr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("GetLocalIP in communication failed")
		return "localhost"
	}
	for _, val := range addr {
		if ip, ok := val.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				localIP = ip.IP.String()
			}
		}
	}
	return localIP
}
