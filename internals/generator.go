package internals

import (
	"errors"
	"fmt"
	"hash/adler32"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"rasyncan/types"
)

var fileList = make(types.FileList, 0)

func GenerateFileList(dir string) types.FileList{
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files{
		fPath := filepath.Join(dir, file.Name())
		if file.IsDir(){
			GenerateFileList(fPath)
		}else{
			fileList = append(fileList, extractMetadata(fPath))
		}
	}
	return fileList
}

func extractMetadata(fpath string) types.FileMetadata{
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

