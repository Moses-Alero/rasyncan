package utils

import (
	"os"
	"io"
	"errors"
	"io/fs"
	"log"
	"fmt"
	"hash/adler32"
	"rasyncan/x/types"

)

func ExtractMetadata(fpath string) types.FileMetadata{
	fInfo, err := os.Stat(fpath)
	if err != nil{
		log.Fatal(err)
	}
	metadata := types.FileMetadata{
		Path: fpath,
		Size: fInfo.Size(),
		Perm: fInfo.Mode().Perm(),
		MTime: fInfo.ModTime(),
		Checksum: GenerateFileHash(fpath),
	}
	return metadata
}

func GenerateFileHash(file string) string{
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		if errors.Is(err, fs.ErrNotExist){
			return ""
		}
		log.Fatal(err)
	}
	
	h := adler32.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x",h.Sum(nil))
}



