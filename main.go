package main

import (
	"app/readers"
	"fmt"
)

func main() {

	/*	var fileReader = readers.NewFileReader()
			fileReader.Execute()
			var customHtmlReader = readers.NewCustomHtmlReaderReader()
			customHtmlReader.Execute()
		var envVarReader = readers.NewEnvVarReader()
		envVarReader.Read()*/
	/*var MapFilesInDisk = readers.NewMapFilesInDisk()
	MappedFiles := MapFilesInDisk.MapFromDisk()

	var MigratedFiles = readers.NewMigratedFiles()
	MigratedFilesData := MigratedFiles.Read()
	for k, v := range MigratedFilesData {
		if _, ok := MappedFiles[k]; ok {
			fmt.Println(k)
			fmt.Println(v.El)
			fmt.Println(v.En)
			fmt.Println("_________")
		}
	}*/
	CategoryUtils := readers.NewCategoryUtils()
	levels := CategoryUtils.Read()
	allCategories := CategoryUtils.SubmitLevels(levels)
	var MapFilesInDisk = readers.NewMapFilesInDisk()
	MappedFiles := MapFilesInDisk.MapFromDisk()

	var MigratedFiles = readers.NewMigratedFiles()
	MigratedFilesData := MigratedFiles.Read()

	for k, v := range MigratedFilesData {
		if mp, ok := MappedFiles[k]; ok {
			fmt.Println("PageId", k)
			if len(v.El.Levels) >= 2 {
				titleEl := allCategories[v.El.Levels[len(v.El.Levels)-2]].NameEl
				titleEn := allCategories[v.El.Levels[len(v.El.Levels)-2]].NameEn
				cat := allCategories[v.El.Levels[len(v.El.Levels)-2]]
				fmt.Println(fmt.Sprintf("El Title:%s En Title:%s", titleEl, titleEn))
				fmt.Println(fmt.Sprintf("El category:%d En category:%d", cat.IdEl, cat.IdEn))
			} else {
				fmt.Println(fmt.Sprintf("El Title:%s En Title:%s", mp.FileName, mp.FileName))
				fmt.Println("No category")
			}
			fmt.Println("IsTemplate", mp.IsTemplate)
			fmt.Println("Path", mp.FilePath)
			fmt.Println("_________")
		}
	}
}
