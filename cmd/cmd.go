package cmd

import (
	"fmt"
	"log"
	"os"
	"rasyncan/internals"
	"rasyncan/x/types"

	"github.com/spf13/cobra"
)

var FileSync = &cobra.Command{
	Use: "a",
	Short: "sync file somewhat",
	Long: "flesh this stuff out a bit more",
	Run: func(cmd *cobra.Command, args []string){
		fmt.Println("File sync has startted")
		if len(args) < 2{
			log.Fatal("Please input the directories you want to sync")
		}
		validateDir(args[0])
		validateDir(args[1])

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

func validateDir(path string){
		file, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if !file.IsDir(){
			log.Fatalf("%s is not a directory", path)
		}
		return
}
