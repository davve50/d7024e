package d7024e

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func (kad *Kademlia) CLI() {
	for {
		fmt.Println("Enter a command:")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		cmd := scanner.Text()

		switch {
		case strings.Contains(cmd, "put "):
			hash := cmd[4:]
			kad.Store([]byte(hash)) // Mayb not correct?
		case strings.Contains(cmd, "get "):
			hash := cmd[4:]
			value := ""
			contacts := make([]Contact, 0)
			kad.LookupData(hash, &value, &contacts)
		case strings.Contains(cmd, "exit"):
			os.Exit(0)
		default:
			fmt.Println("Error: Wrong command.")
		}
	}
}
