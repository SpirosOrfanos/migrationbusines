package readers

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type MapFilesInDisk struct {
	templatePath   string
	customHtmlPath string
}

func NewMapFilesInDisk() *MapFilesInDisk {
	return &MapFilesInDisk{
		templatePath:   os.Getenv("FILE_PATH"),
		customHtmlPath: os.Getenv("CUSTOM_HTML_PATH"),
	}
}

func (srv *MapFilesInDisk) MapFromDisk() map[string]FileInfo {

	fileItems := make(map[string]FileInfo)

	templates, _ := ioutil.ReadDir(srv.templatePath)
	custom, _ := ioutil.ReadDir(srv.customHtmlPath)
	for _, file := range templates {
		fileName := file.Name()
		idName := strings.Split(file.Name(), ".")[0]
		filePath := fmt.Sprintf("%s", fileName)
		fileItems[idName] = FileInfo{
			IsTemplate: true,
			FileName:   fileName,
			IdName:     idName,
			FilePath:   filePath,
		}
	}

	for _, file := range custom {
		fileName := file.Name()
		idName := strings.Split(file.Name(), ".")[0]
		filePath := fmt.Sprintf("%s", fileName)
		fileItems[idName] = FileInfo{
			IsTemplate: false,
			FileName:   fileName,
			IdName:     idName,
			FilePath:   filePath,
		}
	}

	return fileItems
}

type FileInfo struct {
	IsTemplate bool
	FileName   string
	IdName     string
	FilePath   string
}
