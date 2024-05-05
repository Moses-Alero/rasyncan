package internals

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	HOST = "localhost"
	PORT = "8732"
	TYPE = "tcp"
)

func Lreceiver(c, c2 chan bool, c3 chan int64) {
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Println("Starting the server")

	// close listener
	defer listen.Close()

	c <- true
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		go handleRequest(conn, c2, c3)
	}
}

func handleRequest(conn net.Conn, c2 chan bool, c3 chan int64) {
	// incoming request
	val := <-c3
	fmt.Println("handling connection", val)
	buffer := make([]byte, val)
	conn.Read(buffer)
	var fileList FileList
	fmt.Println(buffer)	
	if err := json.Unmarshal(buffer, &fileList); err != nil {
        	log.Fatal(err)
    	}

	fmt.Println("This is the file list", len(fileList))

	// write data to response
	c2 <- true
	// close afterb the sybc has been completed conn
	conn.Close()
}
