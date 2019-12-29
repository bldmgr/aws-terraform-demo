package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Usage: go run client.go <server-IP>:<PORT> [<commands>]\n")
		fmt.Println("Example: go run client.go 3.82.245.52:8081 'pwd, ls -la'")
		return
	}

	links := []string{
		arguments[1],
	}

	if len(os.Args) > 2 {
		arg := os.Args[2]

		c := make(chan string)
		retry := 1

		for _, link := range links {
			go liveCheck(link, c, arg)
		}

		for l := range c {
			if retry != 10 {

				go func() {
					retry += 1
					time.Sleep(5 * time.Second)
					liveCheck(l, c, arg)
				}()
			} else {
				fmt.Println("Error connecting to server ")
				os.Exit(0)
			}
		}

	}

}

func liveCheck(link string, c chan string, send_arg string) {
	conn, err := net.Dial("tcp", link)
	if err != nil {
		fmt.Printf("Error when connecting to server: %v\n", err)
		c <- link
		return
	} else {
		fmt.Println(link, "is up! process might take some time, please wait...")
		_, err := fmt.Fprintf(conn, send_arg+"\n")
		if err != nil {
			fmt.Printf("Error when sending message to server: %v\n", err)
		} else {
			// Listen for reply. 4 is the ASCII value for EOT ("End Of Transmission")
			message, err := bufio.NewReader(conn).ReadString(4)
			if err != nil {
				fmt.Printf("Error when receiving message from server: %s\n", err)
			} else {
				//  fmt.Print(message)
				scanner := bufio.NewScanner(strings.NewReader(message))
				for scanner.Scan() {

					if strings.Contains(scanner.Text(), "ERRCODE:01") {
						err := ioutil.WriteFile("server.err", []byte(message), 0755)
						if err != nil {
							fmt.Printf("Unable to write file: %v", err)
						}
					}
					fmt.Println(scanner.Text())

				}

				if err := scanner.Err(); err != nil {
					// Handle the error
				}

			}
			os.Exit(0)
		}
	}
}
