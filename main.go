package main

import (
	"fmt"
	"rasyncan/cmd"
	"time"
)

func main() {
	fmt.Println("raSyncan")

	start := time.Now()
	fmt.Println("start Time: ", time.Now())
	cmd.Execute()
	fmt.Println("TTR: ", time.Now().Sub(start).Milliseconds(), "ms")

}

//core functionality.
//File Synchronization: Implement the core functionality to compare files and directories
//between source and destination locations and synchronize them accordingly.

//compare files...
//same name
//same size,
//ue rsync algo to check if the file content is inherently the same.

//  rsync algo
