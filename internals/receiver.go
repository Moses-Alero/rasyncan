package internals

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"rasyncan/rsync"
	"rasyncan/types"
	"strings"
)

const (
	HOST = "localhost"
	PORT = "8732"
	TYPE = "tcp"
)

func Lreceiver(pipe types.Pipe){
	//var fileList = make(types.FileList, 0)
	for file := range pipe.RFileChan {
		f := strings.Split(file.Path, pipe.SDir)
		rPath := filepath.Join(pipe.RDir, f[1])
		_, err := os.Stat(rPath)
		if err != nil{
			if errors.Is(err, os.ErrNotExist) {
				fmt.Println("file does not exist")
				f, err := os.Create(rPath)
				defer f.Close()
				if err != nil{
					fmt.Println(err)
					continue
				}
				if err := os.WriteFile(rPath, []byte("rsync dunmmy file"), 0666); err != nil{
					fmt.Println(err)
					continue
				}
				fmt.Println("File created successfully")
			}else{
				continue
			}
		}
		LSync(file.Path, rPath)
		//receive the fijle diff if there's any
		//apply the delta
		
		//fileList = append(fileList, file)
		//verifyFilesToSync(file.Path)
	}
	pipe.C2 <- true
}

func LSync(srcPath string, destPath string){
	srcReader, err := os.Open(srcPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer srcReader.Close()

	rs := &rsync.RSync{}

	// here we store the whole signature in a byte slice,
	// but it could just as well be sent over a network connection for example
	sig := make([]rsync.BlockHash, 0)
	writeSignature := func(bl rsync.BlockHash) error {
		sig = append(sig, bl)
		return nil
	}

	targetReader, err := os.Open(destPath)
	if err != nil {
		fmt.Println(err)
		return
	}


	rs.CreateSignature(targetReader, writeSignature)

	opsOut := make(chan rsync.Operation)
	writeOperation := func(op rsync.Operation) error {
		opsOut <- op
		return nil
	}


	go func() {
		defer close(opsOut)
		rs.CreateDelta(srcReader, sig, writeOperation)
	}()

	srcWriter, err := os.OpenFile(destPath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	srcReader.Seek(0, io.SeekStart)

	rs.ApplyDelta(srcWriter, targetReader, opsOut)

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
