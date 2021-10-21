package d7024e

import "fmt"

type Kademlia struct {
	network    *Network
	routingtab *RoutingTable
	hash       map[string][]byte
}

func InitKademlia(me Contact) *Kademlia {
	newKademlia := &Kademlia{
		routingtab: NewRoutingTable(me),
		hash:       make(map[string][]byte),
	}
	return newKademlia
}

var alpha int = 3

func (kademlia *Kademlia) setNetwork(network *Network) {
	kademlia.network = network
}

func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	var shortlist ContactCandidates
	shortlist.contacts = kademlia.routingtab.FindClosestContacts(target.ID, alpha)
	var visitedNodes []Contact
	var closestNode Contact
	if len(shortlist.contacts) > 0 {
		closestNode = shortlist.contacts[0]
	} else {
		return make([]Contact, 0)
	}

	for len(visitedNodes) < bucketSize && !Contains(visitedNodes, closestNode) {
		// for NODES VISITED < 20 AND CLOSESTNODE IS A NEW NODE
		// find 3 closest UNVISITED NODES and do it all over
		i := 0
		for _, contact := range shortlist.contacts {
			if !Contains(visitedNodes, contact) {
				i++
				kademlia.network.SendPingPacket(contact)
				pong := string(<-kademlia.network.findVal)
				if pong == "pong" {
					kademlia.network.SendFindNodePacket(contact, target)
					visitedNodes = append(visitedNodes, contact)
					k_contacts := <-kademlia.network.findCont
					for _, cont := range k_contacts {
						cont.CalcDistance(target.ID)
						if !Contains(shortlist.contacts, cont) {
							shortlist.contacts = append(shortlist.contacts, cont)
						}
						if cont.Less(&closestNode) {
							closestNode = cont
						}
					}
				}
			}
			if i == 3 {
				break
			}
		}
		shortlist.Sort()
	}

	if shortlist.Len() < bucketSize {
		return shortlist.contacts
	} else {
		return shortlist.contacts[:bucketSize]
	}
}

func (kademlia *Kademlia) LookupData(hash string) (string, []Contact) {
	var shortlist ContactCandidates
	targetID := NewContact(NewKademliaID(hash), "").ID
	shortlist.contacts = kademlia.routingtab.FindClosestContacts(targetID, alpha)
	var closestNode Contact
	if len(shortlist.contacts) > 0 {
		closestNode = shortlist.contacts[0]
	} else {
		return "", make([]Contact, 0)
	}
	var visitedNodes []Contact

	for len(visitedNodes) < bucketSize && !Contains(visitedNodes, closestNode) {
		// for NODES VISITED < 20 AND CLOSESTNODE IS A NEW NODE
		// find 3 closest UNVISITED NODES and do it all over
		i := 0
		for _, contact := range shortlist.contacts {
			if !Contains(visitedNodes, contact) {
				i++
				kademlia.network.SendPingPacket(contact)
				pong := string(<-kademlia.network.findVal)
				if pong == "pong" {
					kademlia.network.SendFindValuePacket(shortlist.contacts[0], hash)
					visitedNodes = append(visitedNodes, contact)
					value := string(<-kademlia.network.findVal)
					k_contacts := <-kademlia.network.findCont
					if value != "" {
						return value, make([]Contact, 0)
					}
					for _, cont := range k_contacts {
						cont.CalcDistance(targetID)
						if !Contains(shortlist.contacts, cont) {
							shortlist.contacts = append(shortlist.contacts, cont)
						}
						if cont.Less(&closestNode) {
							closestNode = cont
						}
					}
				}
			}
			if i == 3 {
				break
			}
		}
		shortlist.Sort()
	}

	if shortlist.Len() < bucketSize {
		return "", shortlist.contacts
	} else {
		return "", shortlist.contacts[:bucketSize]
	}
}

func Contains(visitedNodes []Contact, contact Contact) bool {
	for _, cont := range visitedNodes {
		if cont.ID.Equals(contact.ID) {
			return true
		}
	}
	return false
}

func (kademlia *Kademlia) Store(data []byte) {
	value, ok := kademlia.hash[kademlia.network.me.ID.String()]
	if !ok {
		kademlia.hash[kademlia.network.me.ID.String()] = data
	} else {
		fmt.Println("Immutable data: Already stored:", string(value))
		fmt.Println("Reference to Delimitations #1:")
		return
	}

	contacts := kademlia.LookupContact(kademlia.network.me)

	var res string
	for _, node := range contacts {
		if !node.ID.Equals(kademlia.network.me.ID) {
			kademlia.network.SendStorePacket(node, data)
			res = string(<-kademlia.network.findVal)
			if res == "SUCCESS" {
				fmt.Println("Saved at: ", node.ID.String())
			} else {
				fmt.Println("Could not store data at: ", node.ID.String())
			}
		}
	}
}
