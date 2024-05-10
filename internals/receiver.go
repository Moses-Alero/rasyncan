package internals

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"rasyncan/types"
)

const (
	HOST = "localhost"
	PORT = "8732"
	TYPE = "tcp"
)

func Lreceiver(pipe types.Pipe){
	//var fileList = make(types.FileList, 0)
	for file := range pipe.RFileChan {
		fmt.Println(file)
		//fileList = append(fileList, file)
		//verifyFilesToSync(file.Path)
	}
	pipe.C2 <- true
}

func receiver(pipe types.Pipe) {
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Println("Starting the server")

	// close listener
	defer listen.Close()

	pipe.C1 <- true
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		go handleRequest(conn, pipe)
	}
}

func handleRequest(conn net.Conn, pipe types.Pipe) {
	// incoming request
	var fmetadata types.FileMetadata
	d := gob.NewDecoder(conn)

	if err := d.Decode(&fmetadata); err != nil {
        	log.Fatal(err)
    	}
	//verify what files to sync
	fileList = append(fileList, fmetadata)
	//verifyFilesToSync(fmetadata, pipe)
	pipe.C2 <- true
	conn.Write([]byte(""))
	conn.Close()
}

func verifyFilesToSync(f types.FileMetadata){
	fmt.Println(f.Path)
	//generate checksum for both files 
	//check if the checksums are the same
	//if checksum for the file mathches 
	//skip else sync
	return
}
