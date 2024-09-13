package internals

import (
	"log"
	"os"
	"path/filepath"
	"rasyncan/x/types"
	"rasyncan/x/utils"
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
			fileList = append(fileList, utils.ExtractMetadata(fPath))
		}
	}
	return fileList
}
