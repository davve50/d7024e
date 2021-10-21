package d7024e

import (
	"bufio"
	"fmt"
	"strings"
	"time"
)

func (kad *Kademlia) CLI(test bool, scanner *bufio.Scanner) {
	for {
		if test {
			fmt.Print("[ROOT] Enter a command: ")
		} else {
			fmt.Print("Enter a command: ")
		}

		scanner.Scan()
		cmd := scanner.Text()
		fmt.Println("Running command:", cmd)

		switch {
		case strings.Contains(cmd, "put "):
			hash := cmd[4:]
			kad.Store([]byte(hash))
		case strings.Contains(cmd, "get "):
			hash := cmd[4:]
			if len(hash) != IDLength*2 {
				fmt.Println("[ERROR] Faulty hash")
				break
			}
			value, contacts := kad.LookupData(hash)
			fmt.Println("[CLI] Value:", value)
			fmt.Println("[CLI] Contacts:", contacts)
		case strings.Contains(cmd, "exit"):
			packet := kad.network.CreatePacket("stop_rpc", "", "", "", nil, nil)
			kad.network.SendPacket(packet, kad.network.me.Address)
			time.Sleep(time.Millisecond * 200)
			return
		case strings.Contains(cmd, "list"):
			fmt.Println("Values stored: ")
			for key, element := range kad.hash {
				fmt.Println("\tKey:", key, "=>", "Element:", string(element))
			}
		case strings.Contains(cmd, "me"):
			fmt.Println("ID: ", kad.network.me.ID.String(), " IP: ", kad.network.me.Address)
		case strings.Contains(cmd, "contacts"):
			fmt.Println("Saved contacts: ")
			for _, bucket := range kad.routingtab.buckets {
				for _, cont := range bucket.GetContactAndCalcDistance(NewKademliaID("0000000000000000000000000000000000000000")) {
					fmt.Println("\tID: ", cont.ID.String(), " IP: ", cont.Address)
				}
			}
		default:
			fmt.Println("Error: Wrong command.")
		}
	}
}
