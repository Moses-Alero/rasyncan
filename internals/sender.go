package internals

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)


func Lsender(dir string, c3 chan int64) {
//	tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)
	
	conn, err := net.Dial(TYPE, HOST+":"+PORT)

	fmt.Println("Connecting to servers")
	if err != nil {
		log.Fatal(err)
	}

	fileList := GenerateFileList(dir)		

	buffer, _ := json.Marshal(fileList)

	fmt.Println("data to be sent", len(buffer))
	c3 <- int64(len(buffer))
	conn.Write(buffer)
	

	fmt.Println("Connected and file list sent")
	conn.Close()

}
