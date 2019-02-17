package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	// Initialization
	config, err := loadConfig("conf.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("config info: %+v\n", *config)

	listener, err := net.Listen("tcp", config.ListenAddress)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("server listening on", config.ListenAddress, "...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}

// handleConn handle the connection
func handleConn(conn net.Conn) {
	processConn(conn)

}
