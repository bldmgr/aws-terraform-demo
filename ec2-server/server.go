package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type BarnMessage struct {
	Name    string
	Version string
	Message []string
	Id      int64 `json:"ref"`
	Created time.Time
}

func handleConnection(c net.Conn) {
	startTime := time.Now()
	client := c.RemoteAddr().String()
	fmt.Printf("Serving %s\n", client)
	for {
		// Listen for a message.
		msg, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Printf("[%s] Error received: %s\n", client, err)
			fmt.Printf("[%s] Client will be disconnected.\n", client)
			return
		}
		// Output the received message.
		msgStr := string(msg)
		fmt.Printf("[%s] Message received: %s", client, msgStr)
		// Trim the newline from the message, and check for the special 'STOP' message.
		msgStrTrimmed := strings.TrimSpace(msgStr)
		if msgStrTrimmed == "STOP" {
			fmt.Printf("[%s] Closing connection with client due to STOP signal.\n", client)
			break
		}
		// msgStrTrimmed is a comma separated list of command strings, e.g. "echo hello, ls -la, pwd".
		// Set cmdStrArray to an array with these command strings as elements.
		cmdStrArray := strings.Split(msgStrTrimmed, ",")
		cmdNum := len(cmdStrArray)
		returnMsg := ""
		resultStr := time.Now().Format(time.RFC850)
		for i, cmdStr := range cmdStrArray {
			cmdStrTrimmed := strings.TrimSpace(cmdStr)
			// To allow for complex command strings, like "ps aux | grep hello", without having to do
			// any parsing, run each command in its own bash process.
			cmd := exec.Command("bash", "-c", cmdStrTrimmed)
			fmt.Printf("[%s] Running command %v/%v, '%s', and waiting for it to finish...\n", client, i+1, cmdNum, cmdStrTrimmed)
			returnMsg += "\n\nCommand (" + strconv.Itoa(i+1) + "/" + strconv.Itoa(cmdNum) + "): " + cmdStrTrimmed + "\n"
			outErr, err := cmd.CombinedOutput()
			fmt.Print(string(outErr))
			if err != nil {
				fmt.Printf("[%s] Command %v/%v finished with error: %v\n", client, i+1, cmdNum, err)
				returnMsg += "ERRCODE:01"
				resultStr = "ERRCODE:01"
			} else {
				fmt.Printf("[%s] Command %v/%v finished without error.\n", client, i+1, cmdNum)
				returnMsg += "SUCCEEDED"
			}
			returnMsg += "\n---\n" + string(outErr)
		}

		// Write standard out
		returnMsg = resultStr + returnMsg + "\n"
		returnArray := []byte(returnMsg)
		c.Write(returnArray)

		// Write BarnMessage
		timeArray := taskDuration(startTime, resultStr)
		c.Write(timeArray)

		// 4 is the ASCII value for EOT ("End Of Transmission")
		// returnArray = append(returnArray, 4)
		endArray := []byte("\n")
		endArray = append(endArray, 4)
		c.Write(endArray)
	}
	c.Close()
}

func taskDuration(startTime time.Time, resultStr string) []byte {
	endTime := time.Now()
	sessionTime := endTime.Sub(startTime)

	barnmessage := BarnMessage{
		Name:    "Ec2 Worker",
		Version: "1.0.0",
		Message: []string{resultStr, fmt.Sprintf("%v", sessionTime)},
		Id:      999,
		Created: time.Now(),
	}

	var jsonData []byte
	jsonData, error := json.MarshalIndent(barnmessage, "", "    ")
	if error != nil {
		log.Println(error)
	}
	return jsonData
}

func main() {
	arguments := os.Args
	PORT := "8081"
	if len(arguments) == 1 {
		fmt.Println("No port provided. Port 8081 will be used by default.")
	} else {
		PORT = arguments[1]
	}
	l, err := net.Listen("tcp4", ":"+PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	fmt.Println("Starting server, listening on port " + PORT)

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}
