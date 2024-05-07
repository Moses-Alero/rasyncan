package internals

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"rasyncan/types"
)


func Lsender(pipe types.Pipe) {
	conn, err := net.Dial(TYPE, HOST+":"+PORT)

	fmt.Println("Connecting to servers")
	if err != nil {
		log.Fatal(err)
	}

	fileList := GenerateFileList(pipe.SDir)		
	
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(fileList); err != nil{
		log.Fatal(err)
	}

	fmt.Println("data to be sent", len(buf.Bytes()))
	conn.Write(buf.Bytes())

	fmt.Println("Connected and file list sent")
	conn.Close()

}
