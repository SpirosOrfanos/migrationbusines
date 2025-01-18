package main

import (
	"app/model"
	"app/readers"
	"fmt"
	"os"
)

func main() {
	standaloneLoader()
}

func standaloneLoader() {
	var MapFilesInDisk = readers.NewMapFilesInDisk()
	MappedFiles := MapFilesInDisk.MapFromDisk()
	var fileReader = readers.NewFileReader()
	var customHtmlReader = readers.NewCustomHtmlReaderReader()
	for _, v := range MappedFiles {
		fmt.Println(v.FileName, v.IsTemplate)
		if v.IsTemplate {
			fileReader.Splitter(model.Insertable{FilePath: v.FileName,
				PageId:     v.IdName,
				TitleEn:    v.IdName,
				TitleEl:    v.IdName,
				IsBusiness: false})
		} else {
			customHtmlReader.Splitter(model.Insertable{FilePath: v.FileName,
				PageId:     v.IdName,
				TitleEn:    v.IdName,
				TitleEl:    v.IdName,
				IsBusiness: false})
		}

	}

}
func standAlone(insertable model.Insertable) {

}

/*
	func standAlone() {
		var fileReader = readers.NewFileReader()
		fileReader.Splitter(model.Insertable{FilePath: "free-roaming-bus.xlsx",
			PageId:     "free-roaming-bus",
			TitleEn:    "free-roaming-bus_en",
			TitleEl:    "free-roaming-bus_el",
			IsBusiness: true})
	}
*/
func serializeThem() {
	var customHtmlReader = readers.NewCustomHtmlReaderReader()
	var fileReader = readers.NewFileReader()
	var MapFilesInDisk = readers.NewMapFilesInDisk()
	MappedFiles := MapFilesInDisk.MapFromDisk()
	var MigratedFiles = readers.NewMigratedFiles()
	MigratedFilesData := MigratedFiles.Read()

	dataPath := os.Getenv("FILE_PATH")
	htmlPath := os.Getenv("CUSTOM_HTML_PATH")
	for k, v := range MigratedFilesData {
		if mp, ok := MappedFiles[k]; ok {
			insertable := model.Insertable{
				PageId:       k,
				CategoriesEn: make([]int, 0),
				CategoriesEl: make([]int, 0),
			}

			if len(v.El.Levels) >= 1 {
				titleEl := v.El.Levels[len(v.El.Levels)-1]
				titleEn := v.El.Levels[len(v.El.Levels)-1]
				insertable.TitleEl = titleEl
				insertable.TitleEn = titleEn
			} else {
				insertable.TitleEl = mp.FileName
				insertable.TitleEn = mp.FileName
			}
			insertable.FilePath = mp.FilePath
			insertable.IsBusiness = v.El.IsBone

			if mp.IsTemplate {
				path := fmt.Sprintf("%s%s%s", dataPath, "\\", mp.FileName)
				fmt.Println("Template", path)
				fileReader.ReadExcelContent(mp.FileName)
				//fmt.Println(content.EN)

			} else {
				path := fmt.Sprintf("%s%s%s", htmlPath, "\\", mp.FileName)
				fmt.Println("Custom", path)
				customHtmlReader.ReadExcelContent(mp.FileName)
				//fmt.Println(content.EN)

			}
		}
	}
}
func insertThem() {
	var customHtmlReader = readers.NewCustomHtmlReaderReader()
	var fileReader = readers.NewFileReader()
	CategoryUtils := readers.NewCategoryUtils()
	levels := CategoryUtils.Read()
	allCategories := CategoryUtils.SubmitLevels(levels)
	var MapFilesInDisk = readers.NewMapFilesInDisk()
	MappedFiles := MapFilesInDisk.MapFromDisk()

	var MigratedFiles = readers.NewMigratedFiles()
	MigratedFilesData := MigratedFiles.Read()

	for k, v := range MigratedFilesData {
		if mp, ok := MappedFiles[k]; ok {
			insertable := model.Insertable{
				PageId:       k,
				CategoriesEn: make([]int, 0),
				CategoriesEl: make([]int, 0),
			}
			fmt.Println("", mp.FileName)
			if len(v.El.Levels) >= 1 {
				titleEl := v.El.Levels[len(v.El.Levels)-1]
				titleEn := v.El.Levels[len(v.El.Levels)-1]
				if len(v.El.Levels) >= 2 {
					cat := allCategories[v.El.Levels[len(v.El.Levels)-2]]
					insertable.CategoriesEn = append(insertable.CategoriesEn, cat.IdEn)
					insertable.CategoriesEl = append(insertable.CategoriesEl, cat.IdEl)
				}
				insertable.TitleEl = titleEl
				insertable.TitleEn = titleEn
			} else {
				insertable.TitleEl = mp.FileName
				insertable.TitleEn = mp.FileName
			}
			insertable.FilePath = mp.FilePath
			insertable.IsBusiness = v.El.IsBone

			if mp.IsTemplate {
				fileReader.Splitter(insertable)
			} else {
				customHtmlReader.Splitter(insertable)
			}
		}
	}
}
