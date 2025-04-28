package main

import (
	"fmt"
	"net"
	"os"
)

func establishConnections(addresses []string) []net.Conn {
	conns := make([]net.Conn, len(addresses))

	for i := 0; i < len(addresses); i++ {
		var err error
		fmt.Println("establishing connection " + addresses[i])
		conns[i], err = net.Dial("tcp", addresses[i])
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			os.Exit(1)
		}
	}
	return conns
}

// Propagates Sync Reqs to Other Replicas
func broadcast(conns []net.Conn, content []byte) {
	content = append(content, 0x00)
	for i := 0; i < len(conns); i++ {
		//w := bufio.NewWriter(conns[i])
		// Send to socket connection.
		//conns[i].Write([]byte(content))
		//_, err := w.Write([]byte(content))
		//fmt.Println("Writing to socket ", i)
		//fmt.Println("content len is ", len(content))

		// Append a custom delimiter (0x00) to the content

		_, err := conns[i].Write([]byte(content))
		if err != nil {
			fmt.Println("Error socket writing:", err.Error())
			os.Exit(1)
		}
		//fmt.Println("Wrote to socket bytes", n)
	}
}
