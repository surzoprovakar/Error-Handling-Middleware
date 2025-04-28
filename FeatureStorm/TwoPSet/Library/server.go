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
var set *Set

var conns []net.Conn

var id2trust map[int]float64

func do_actions(actions []Operation2Value) {

	//sleep for 10 secs, so other replicase
	//have time to get started
	time.Sleep(10 * time.Second)
	fmt.Println("Starting to do_actions")

	for _, act2val := range actions {
		if act2val.action == "Add" || act2val.action == "Remove" {
			set.SetVal(act2val.action, act2val.operator)
			set.Print()
		} else if act2val.action == "Broadcast" {
			fmt.Println("processing Broadcast")
			if conns == nil { //establish connecitons on first broadcast
				conns = establishConnections(hosts)
			}
			fmt.Print("About to broadcast Counter")
			set.Print()
			broadcast(conns, set.ToMarshal())
		} else if strings.Contains(act2val.action, "undo") { //undo executing
			fmt.Println("Executing local undo")
			undo_data := strings.Split(act2val.action, "_")
			// fmt.Println("undo-data:", undo_data)
			// fmt.Println(len(undo_data))
			if undo_data[1] == "Add" {
				fmt.Println("Undoing Add of element", undo_data[2], " at LT:", act2val.operator)
				undo_val, _ := strconv.Atoi(undo_data[2])
				set.SetVal("Remove", undo_val)
				set.Print()
			} else if undo_data[1] == "Remove" {
				fmt.Println("Undoing Remove of element", undo_data[2], " at LT:", act2val.operator)
				undo_val, _ := strconv.Atoi(undo_data[2])
				set.SetVal("Add", undo_val)
				set.Print()
			}
			// fmt.Println("undo data:", undo_data)
			replica_2d := []string{strconv.Itoa(set.Id()), ""}
			var undo_sync [][]string

			undo_LT := strconv.Itoa(act2val.operator)
			undo_updates := []string{undo_data[0], undo_data[1] + "_" + undo_data[2] + "_" + undo_LT}
			undo_sync = append(undo_sync, undo_updates)
			undo_sync = append(undo_sync, replica_2d)

			// fmt.Println("undo sync")
			// fmt.Println(undo_sync)

			undo_json_data, err := json.Marshal(undo_sync)

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
			if number, err = strconv.Atoi(act2val.action); err != nil {
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
	set = NewSet(id)

	ip_address := input[1]

	hosts = ReadFile(input[2])
	actions := ReadFile_Actions(input[3])
	//id2trust = ReadFile_TrustGradient(input[4])
	//trusts_list := get_trust_from_map(hosts)

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
		buffer := make([]byte, 704)
		//c := bufio.NewReader(conn)
		fmt.Println("starting to read")
		n, err := conn.Read(buffer)
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
	set.Merge(rid, updates)

	// Restart the process.
	handleConnection(conn)
}
