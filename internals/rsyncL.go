package internals

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"rasyncan/rsync"
	"rasyncan/x/types"
	"rasyncan/x/utils"
	"strings"
)

func Lreceiver(pipe types.Pipe) {
	for file := range pipe.RFileChan {
		f := strings.Split(file.Path, pipe.SDir)
		rPath := filepath.Join(pipe.RDir, f[1])
		_, err := os.Stat(rPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				setupFileDir(pipe.RDir, f[1])
			} else {
				continue
			}
		}
		rFile := utils.ExtractMetadata(rPath)
		if verifyFilesToSync(file, rFile) {
			sync(file.Path, rPath)
		}
	}
	pipe.C2 <- true
	close(pipe.C2)
}

func Lsender(pipe types.Pipe, fList types.FileList) {

	for _, file := range fList {
		pipe.RFileChan <- file
	}
	close(pipe.RFileChan)
}

func verifyFilesToSync(f1, f2 types.FileMetadata) bool {
	//compare checksum for both files
	if f1.Checksum == f2.Checksum &&
		f1.Perm == f2.Perm &&
		f1.Size == f2.Size {
		return false
	}

	fmt.Println("Pathm: ", f1.Path, " <-> ", f2.Path)
	fmt.Println("Checksum: ", f1.Checksum, " <-> ", f2.Checksum)
	fmt.Println("Perm: ", f1.Perm, " <-> ", f2.Perm)
	fmt.Println("Size: ", f1.Size, " <-> ", f2.Size)

	return true
}

// this is not a very descriptive name for this function
// this func creates the directory if it does not exist and also creates the file
// if the directory exists
func setupFileDir(root, path string) {
	file, err := os.Create(filepath.Join(root, path))
	defer file.Close()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			folderList := strings.Split(path, "/")
			fmt.Println(folderList)
			f := ""
			for _, folder := range folderList[:len(folderList)-1] {
				f = filepath.Join(f, folder)
			}
			fmt.Println("this is a directory >>>", filepath.Join(root, f))
			if err := os.MkdirAll(filepath.Join(root, f), 0750); err != nil {
				fmt.Println(err)
				return
			}
			if err = os.WriteFile(filepath.Join(root, path), []byte{}, 0660); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
	return
}

func sync(srcPath string, destPath string) {
	srcReader, err := os.Open(srcPath)
	defer srcReader.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

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
	defer srcWriter.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	srcReader.Seek(0, io.SeekStart)

	rs.ApplyDelta(srcWriter, targetReader, opsOut)
}
