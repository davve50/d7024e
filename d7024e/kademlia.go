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

	var closest Contact
	alphaClosest := kademlia.routingtab.FindClosestContacts(target.ID, alpha)
	if len(alphaClosest) != 0 {
		closest = alphaClosest[0]
	}

	results := make([]Contact, 0)

	vl := &VisitList{ // Checks if SendFindNodePacket for that node has been called
		list:    make([]Contact, 0),
		visited: make(map[Contact]bool),
	}

	resVl := make(map[Contact]bool) // Checks if node has been added to visitList

	// checking my alpha closest nodes
	for _, node := range alphaClosest {
		vl.list = append(vl.list, node)
		resVl[node] = true
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
				if !resVl[node] {
					resVl[node] = true
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
		*contacts = vl.list[:bucketSize]
	} else {
		*contacts = vl.list
	}
}

func (kademlia *Kademlia) LookupData(hash string, value *string, contacts *[]Contact) {
	nodesToCheck := 0

	// We store results from SendFindValuePacket in these
	resultContacts := make([]Contact, 0)
	resultHash := ""
	var closest Contact

	target := NewContact(NewKademliaID(hash), "")
	alphaClosest := kademlia.routingtab.FindClosestContacts(target.ID, alpha)
	if len(alphaClosest) != 0 {
		closest = alphaClosest[0]
	}

	noValueList := make([]Contact, 0)

	vl := &VisitList{ // Checks if SendFindValuePacket for that node has been called
		list:    make([]Contact, 0),
		visited: make(map[Contact]bool),
	}

	resVl := make(map[Contact]bool) // Checks if node has been added to visitList

	// checking my alpha closest nodes
	for _, node := range alphaClosest {
		vl.list = append(vl.list, node)
		resVl[node] = true
	}

	//Check alpha closests nodes if someone knows value
	for nodesToCheck < len(vl.list) {
		nodesToCheck++
		kademlia.network.SendFindValuePacket(&vl.list[nodesToCheck], hash, &resultContacts, &resultHash)
		// If we found a value
		if resultHash != "" {
			if len(noValueList) != 0 {
				fmt.Println("Stored in closest node")
				kademlia.network.SendStorePacket(&noValueList[0], []byte(resultHash))
			}
			*value = resultHash
			return
		}
		for _, node := range resultContacts {
			if (node.Address != kademlia.network.me.Address) && (node.ID != nil) {
				node.CalcDistance(kademlia.network.me.ID)
				noValueList = append(noValueList, node)
			}
			noValueList = iSort(noValueList, &target)
		}
	}

	// No value found, but up to alpha*k new contacts to ask.

	for nodesToCheck > 0 {
		// Loop through and check if any of the found contacts are not added to the list of to be checked
		for _, node := range resultContacts {
			if (node.Address != kademlia.network.me.Address) && (node.ID != nil) {
				if !resVl[node] {
					resVl[node] = true
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
					kademlia.network.SendFindValuePacket(&node, hash, &resultContacts, &resultHash)
					vl.visited[node] = true
					nodesToCheck++
				}
			}
		} else {
			break
		}
	}

	fmt.Println("[LookupValue] Checking if every node has been visited:")
	for _, node := range vl.list {
		fmt.Printf("\t%s = %t\n", node.ID.String(), vl.visited[node])
	}

	if len(vl.list) > bucketSize {
		*contacts = vl.list[:bucketSize]
	} else {
		*contacts = vl.list
	}
}

func (kademlia *Kademlia) Store(data []byte) {
	kademlia.hash[kademlia.network.me.ID.String()] = data // This should be an argument
	contacts := make([]Contact, 0)
	kademlia.LookupContact(*kademlia.network.me, &contacts)
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
