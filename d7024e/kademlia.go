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

type VisitList struct {
	list    []Contact
	visited map[Contact]bool
}

var alpha int = 3

func (kademlia *Kademlia) setNetwork(network *Network) {
	kademlia.network = network
}

func (kademlia *Kademlia) LookupContact(target Contact, contacts *[]Contact) {
	nodesToCheck := 0

	alphaClosest := kademlia.routingtab.FindClosestContacts(target.ID, alpha)
	closest := alphaClosest[0]

	results := make([]Contact, 0)

	vl := &VisitList{
		list:    make([]Contact, 0),
		visited: make(map[Contact]bool),
	}

	// checking my alpha closest nodes
	for _, node := range alphaClosest {
		vl.list = append(vl.list, node)
		vl.visited[node] = true
	}

	// getting contacts from my alpha closest nodes
	for nodesToCheck < len(vl.list) {
		kademlia.network.SendFindNodePacket(&vl.list[nodesToCheck], &target, &results)

		vl.visited[vl.list[nodesToCheck]] = true
		nodesToCheck++
	}

	// Loop through the nodes we need to check to find the shortest path
	for nodesToCheck > 0 {
		// Loop through and check if any of the found contacts are not added to the list of to be checked
		for _, node := range results {
			if (node.Address != kademlia.network.me.Address) && (node.ID != nil) {
				if !vl.visited[node] {
					vl.visited[node] = true
					node.CalcDistance(kademlia.network.me.ID)
					vl.list = append(vl.list, node)
				}
			}
		}

		// Sort the list by lowest distance
		vl.list = iSort(vl.list, &target)
		nodesToCheck--

		// If there is a new closest node to the target we need to check the alpha closest nodes
		if closest.ID.String() != vl.list[0].ID.String() {
			closest = vl.list[0]
			for i, node := range vl.list {
				if i >= alpha {
					break
				}

				if !vl.visited[node] {
					kademlia.network.SendFindNodePacket(&node, &target, &results)

					vl.visited[node] = true
					nodesToCheck++
				}
			}
		} else {
			break
		}
	}

	fmt.Println("[LookupContact] Checking if every node has been visited:")
	for _, node := range vl.list {
		fmt.Printf("\t%s = %t\n", node.ID.String(), vl.visited[node])
	}

	if len(vl.list) > bucketSize {
		*contacts = vl.list[:20]
	} else {
		*contacts = vl.list
	}
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(key *KademliaID, data []byte) {
	contacts := make([]Contact, 0)
	contact := NewContact(key, "")
	kademlia.LookupContact(contact, &contacts)
	for _, node := range contacts {
		kademlia.network.SendStorePacket(&node, data)
	}
}

// inspired by golangprograms.com insertion sort
func iSort(items []Contact, target *Contact) []Contact {
	var n = len(items)

	for i := 1; i < n; i++ {
		j := i
		for j > 0 {
			distL := items[j-1].ID.CalcDistance(target.ID)
			distR := items[j].ID.CalcDistance(target.ID)

			if distL.Less(distR) {
				items[j-1], items[j] = items[j], items[j-1]
			}
			j = j - 1
		}
	}
	return items
}
