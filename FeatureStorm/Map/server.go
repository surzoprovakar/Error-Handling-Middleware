package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Application constants, defining host, port, and protocol.
const (
	//connHost = "localhost"
	//connPort = "8080"
	connType = "tcp"
)

var hosts []string
var mapReplica *Map

var conns []net.Conn

func do_actions(actions []string) {

	//sleep for 5 secs, so other replicase
	//have time to get started
	time.Sleep(5 * time.Second)
	fmt.Println("Starting to do_actions")
	for _, action := range actions {
		if strings.Contains(action, ":") {
			action_data := strings.Split(action, ":")
			opt := action_data[0]
			key := action_data[1]
			if opt == "Add" {
				value, _ := strconv.Atoi(action_data[2])
				mapReplica.Add(key, value)
				mapReplica.Print()
			} else if opt == "Delete" {
				mapReplica.Delete(key)
				mapReplica.Print()
			} else if opt == "Update" {
				value, _ := strconv.Atoi(action_data[2])
				mapReplica.Update(key, value)
				mapReplica.Print()
			}
		} else if action == "Broadcast" {
			// // func before boradcast
			// b1 := mapReplica.ToMarshal()
			// rid, rupdates := FromMarshalData(b1)
			// fmt.Println("test before: rid= ", rid)
			// fmt.Println("test before: rupdates= ", rupdates)

			fmt.Println("processing Broadcast")
			if conns == nil { //establish connecitons on first broadcast
				conns = establishConnections(hosts)
			}
			//conns = establishConnections(hosts)
			fmt.Println("About to broadcast Map")
			mapReplica.Print()
			broadcast(conns, mapReplica.ToMarshal())
		} else { //assume it is delay
			var err error
			var number int
			if number, err = strconv.Atoi(action); err != nil {
				panic(err)
			}

			time.Sleep(time.Duration(number) * time.Second)
		}

	}
}

func main() {

	input := os.Args[1:]
	if len(input) != 4 {
		println("Usage: map_id ip_address crdt_socket_server Replicas'_Addresses.txt Actions.txt")
		os.Exit(1)
	}

	//establish connections using the addresses from the first input file
	//read the execution steps from the second input file
	//execute the script

	var err error
	var id int
	if id, err = strconv.Atoi(input[0]); err != nil {
		panic(err)
	}
	mapReplica = NewMap(id)
	ip_address := input[1]
	hosts = ReadFile(input[2])
	actions := ReadFile(input[3])

	// Start the server and listen for incoming connections.
	fmt.Println("Starting " + connType + " server on " + ip_address)
	l, err := net.Listen(connType, ip_address)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()

	go do_actions(actions)

	// run loop forever, until exit.
	for {
		// Listen for an incoming connection.
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			return
		}
		fmt.Println("Client connected.")

		// Print client connection address.
		fmt.Println("Client " + c.RemoteAddr().String() + " connected.")

		// Handle connections concurrently in a new goroutine.
		go handleConnection(c)
	}
}

// handleConnection handles logic for a single connection request.
func handleConnection(conn net.Conn) {
	// Buffer client input until a newline.
	//buffer, err := bufio.NewReader(conn).ReadBytes('\n')

	/*
		buffer := make([]byte, 1088)
		//c := bufio.NewReader(conn)
		fmt.Println("starting to read")
		_, err := conn.Read(buffer)
		// Close left clients.
		if err != nil {
			fmt.Println("Client left.")
			conn.Close()
			return
		}
	*/

	reader := bufio.NewReader(conn)
	// reqs, err := reader.ReadBytes('\n')
	var message []byte

	for {
		b, err := reader.ReadByte()
		if err != nil {
			fmt.Println("Client left.")
			conn.Close()
			return
		}

		// Check for delimiter
		if b == 0x00 {
			break
		}

		message = append(message, b)
	}
	// fmt.Println("reqs: ", reqs)

	rid, updates := FromMarshalData(message)
	// fmt.Println(updates)
	mapReplica.Merge(rid, updates)

	// Restart the process.
	handleConnection(conn)
}
