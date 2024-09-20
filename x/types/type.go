package types

import (
	"os"
	"time"
)

type Pipe struct {
	Exit     chan bool
	FileChan chan FileMetadata
	SDir     string
	RDir     string
}

type FileMetadata struct {
	Path     string
	Size     int64
	Perm     os.FileMode
	MTime    time.Time
	Checksum string
}

type FileList []FileMetadata
