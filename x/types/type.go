package types

import(
	"os"
	"time"
)


type Pipe struct{
	C1 chan bool 
	C2 chan bool
	SFileChan chan FileMetadata
	RFileChan chan FileMetadata
	SDir string
	RDir string
}

type FileMetadata struct{
	Path  string
	Size  int64
	Perm  os.FileMode
	MTime time.Time
	Checksum string
}

type FileList []FileMetadata




