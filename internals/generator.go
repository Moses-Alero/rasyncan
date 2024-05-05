package internals

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)



type FileMetadata struct{
	path  string
	size  int64
	perm  os.FileMode
	mTime time.Time
	checksum string
}

type FileList []FileMetadata

var fileList FileList = make(FileList, 0)

func GenerateFileList(dirPath string) FileList{
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files{
		fPath := filepath.Join(dirPath, file.Name())
		if file.IsDir(){
			GenerateFileList(fPath)
		}else{
			fileList = append(fileList, extractMetadata(fPath))
		}
	}
	return fileList
}

func extractMetadata(fpath string) FileMetadata{
	fInfo, err := os.Stat(fpath)
	if err != nil{
		log.Fatal(err)
	}
	metadata := FileMetadata{
		path: fpath,
		size: fInfo.Size(),
		perm: fInfo.Mode().Perm(),
		mTime: fInfo.ModTime(),
		checksum: generateFileHash(fpath),
	}
	return metadata
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

