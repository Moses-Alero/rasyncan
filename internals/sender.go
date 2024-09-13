package internals

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"rasyncan/x/types"
)

func Lsender(pipe types.Pipe, fList types.FileList){

	for _, file:= range fList{
		pipe.RFileChan <- file
	}
	close(pipe.RFileChan)
}

func sender(pipe types.Pipe) {
	conn, err := net.Dial(TYPE, HOST+":"+PORT)

	fmt.Println("Connecting to servers")
	if err != nil {
		log.Fatal(err)
	}
		
	var fileList  = make(types.FileList, 0)

	for file := range pipe.SFileChan { 
		fileList = append(fileList, file)
		buf := new(bytes.Buffer)
		enc := gob.NewEncoder(buf)
		if err := enc.Encode(file); err != nil{
			log.Fatal(err)
		}
		fmt.Println("data to be sent", len(buf.Bytes()))
		conn.Write(buf.Bytes())
		conn.Read(make([]byte, 1024))
		conn.Close()

	}
	fmt.Println("Connected and file list sent")
	
}
