package d7024e

type Kademlia struct {
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

func (kademlia *Kademlia) LookupContact(me *Contact, target Contact, contacts *[]Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
