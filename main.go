package main

import (
	"app/model"
	"app/readers"
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

	/*var customHtmlReader = readers.NewCustomHtmlReaderReader()
	CategoryUtils := readers.NewCategoryUtils()
	levels := CategoryUtils.Read()
	allCategories := CategoryUtils.SubmitLevels(levels)
	var MapFilesInDisk = readers.NewMapFilesInDisk()
	MappedFiles := MapFilesInDisk.MapFromDisk()

	var MigratedFiles = readers.NewMigratedFiles()
	MigratedFilesData := MigratedFiles.Read()*/

	/*for k, v := range MigratedFilesData {
		if mp, ok := MappedFiles[k]; ok {
			insertable := model.Insertable{
				PageId:       k,
				CategoriesEn: make([]int, 0),
				CategoriesEl: make([]int, 0),
			}

			if len(v.El.Levels) >= 2 {
				titleEl := allCategories[v.El.Levels[len(v.El.Levels)-2]].NameEl
				titleEn := allCategories[v.El.Levels[len(v.El.Levels)-2]].NameEn
				cat := allCategories[v.El.Levels[len(v.El.Levels)-2]]
				insertable.CategoriesEn = append(insertable.CategoriesEn, cat.IdEn)
				insertable.CategoriesEl = append(insertable.CategoriesEl, cat.IdEl)
				insertable.TitleEl = titleEl
				insertable.TitleEn = titleEn
			} else {
				insertable.TitleEl = mp.FileName
				insertable.TitleEn = mp.FileName
			}
			insertable.FilePath = mp.FilePath
			insertable.IsBusiness = v.El.IsBone
			if !mp.IsTemplate {
				customHtmlReader.Splitter(insertable)
			}
		}
	}*/
	var fileReader = readers.NewFileReader()
	fileReader.Splitter(model.Insertable{FilePath: "afoi_konstantinidis.xlsx",
		PageId:     "afoi-konstantinidis",
		TitleEn:    "afoi_konstantinidis_en",
		TitleEl:    "afoi_konstantinidis_el",
		IsBusiness: true})
}
