package d7024e

import (
	"bufio"
	"fmt"
	"strings"
)

func (kad *Kademlia) CLI(scanner *bufio.Scanner) {
	for {
		fmt.Print("Enter a command: ")
		//scanner := bufio.NewScanner(reader)
		scanner.Scan()
		cmd := scanner.Text()
		//fmt.Print("Running command: ", cmd)

		switch {
		case strings.Contains(cmd, "put "):
			fmt.Println(kad.hash)
			hash := cmd[4:]
			kad.Store([]byte(hash)) // Mayb not correct?
			fmt.Println(kad.hash)
		case strings.Contains(cmd, "get "):
			hash := cmd[4:]
			if len(hash) != IDLength*2 {
				fmt.Println("Wrong input")
				break
			}
			value := ""
			contacts := make([]Contact, 0)
			kad.LookupData(hash, &value, &contacts)
			fmt.Println("[CLI] Value:", value)
			fmt.Println("[CLI] Contacts:", contacts)
		case strings.Contains(cmd, "exit"):
			//Exit here some otherway
			return
			//os.Exit(0)
		case strings.Contains(cmd, "sendPing"):
			contact := NewContact(NewKademliaID("00000000000000000000000000000000deadc0de"), "localhost:8080")
			kad.network.SendPingPacket(&contact)
		default:
			fmt.Println("Error: Wrong command.")
		}
	}
}
