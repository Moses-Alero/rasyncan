package cmd

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"rasyncan/internals"
	"rasyncan/types"

	"github.com/spf13/cobra"
)

var FileSync = &cobra.Command{
	Use: "a",
	Short: "sync file somewhat",
	Long: "flesh this stuff out a bit more",
	Run: func(cmd *cobra.Command, args []string){
		fmt.Println("File sync has startted")
		pipe := types.Pipe{
			C1: make(chan bool),
			C2: make(chan bool),
			RFileChan: make(chan types.FileMetadata),
			SFileChan: make(chan types.FileMetadata),
			SDir: args[0],
			RDir: args[1],
		}
		//validate the directories exist locally
		//the sender directory must exist 
		//the receiver directory may be optional and created if it does not exist 
		//user must specify if they for the file to exist
		fileList :=  internals.GenerateFileList(pipe.SDir)		
		go internals.Lsender(pipe, fileList)
		go internals.Lreceiver(pipe)
			
		//keep main routine running until receiver go routine is completed
		for {
			if <-pipe.C2 {
				break
			}
		}
	},
}

func compareDir(dirs ...string){
	files, err := os.ReadDir(dirs[0])
		if err != nil{
			log.Fatal(err)
		}

		for _ ,file := range files{
			fileDir := make([]string, 2)
			fileDir = []string {
				filepath.Join(dirs[0], file.Name()),
				filepath.Join(dirs[1], file.Name()),
			}

			if file.IsDir(){
				compareDir(fileDir...)
			}else {
				f1, f2 := fileDir[0], fileDir[1]
				//do file comparison logic
				compareFile(f1, f2)
			}

	}
}

func generateFileHash(file string) string{
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		if errors.Is(err, fs.ErrNotExist){
			return ""
		}
		log.Fatal(err)
	}
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x",h.Sum(nil))
}

func compareFile(f1, f2 string) {
	if generateFileHash(f1) == generateFileHash(f2){
		fmt.Println("Files match")
	}else{
		fmt.Println("No luck")
	}	
}
