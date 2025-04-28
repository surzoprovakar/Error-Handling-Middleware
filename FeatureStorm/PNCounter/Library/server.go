package main

import (
	"bufio"
	"encoding/json"
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
var counter *Counter

var conns []net.Conn

func do_actions(actions []string) {

	//sleep for 10 secs, so other replicase
	//have time to get started
	time.Sleep(10 * time.Second)
	fmt.Println("Starting to do_actions")
	for _, action := range actions {
		if action == "Inc" {
			counter.Inc()
			fmt.Println(counter.Print())
		} else if action == "Dec" {
			counter.Dec()
			fmt.Println(counter.Print())
		} else if action == "Broadcast" {
			fmt.Println("processing Broadcast")
			if conns == nil { //establish connecitons on first broadcast
				conns = establishConnections(hosts)
			}
			//conns = establishConnections(hosts)
			fmt.Println("About to broadcast Counter" + counter.Print())
			broadcast(conns, counter.ToMarshal())
		} else if strings.Contains(action, "undo") { //undo executing
			fmt.Println("Executing local undo")
			undo_data := strings.Split(action, " ")
			if undo_data[1] == "Inc" {
				fmt.Println("Undoing Inc of LT:", undo_data[2])
				counter.Dec()
				fmt.Println(counter.Print())
			} else if undo_data[1] == "Dec" {
				fmt.Println("Undoing Dec of LT:", undo_data[2])
				counter.Inc()
				fmt.Println(counter.Print())
			}
			// fmt.Println("undo data:", undo_data)
			undo_data_wid := append([]string{strconv.Itoa(counter.Id())}, undo_data...)
			undo_json_data, err := json.Marshal(undo_data_wid)

			if err != nil {
				fmt.Println("Error while marshaling undo data")
			}

			if conns == nil { //establish connecitons on first broadcast
				conns = establishConnections(hosts)
			}
			//conns = establishConnections(hosts)
			fmt.Println("About to broadcast undo information to other replicas")
			broadcast(conns, undo_json_data)

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
		println("Usage: counter_id ip_address crdt_socket_server Replicas'_Addresses.txt Actions.txt")
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
	counter = NewCounter(id)

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
	reqs, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Println("Client left.")
		conn.Close()
		return
	}

	rid, updates := FromMarshalData(reqs)
	// fmt.Println(updates)
	counter.Merge(rid, updates)

	// Restart the process.
	handleConnection(conn)
}
