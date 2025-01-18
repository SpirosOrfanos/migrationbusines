package readers

import (
	"github.com/xuri/excelize/v2"
	"os"
	"regexp"
)

type MigratedFiles struct {
	filePath string
}

func NewMigratedFiles() *MigratedFiles {
	return &MigratedFiles{filePath: os.Getenv("MIGRATED_FILES_PATH")}
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func (reader *MigratedFiles) Read() map[string]BilingualRowData {
	response := make(map[string]BilingualRowData)
	f, _ := excelize.OpenFile(reader.filePath)
	rowsEl, _ := f.GetRows("el")
	rowsEn, _ := f.GetRows("en")
	for index, _ := range rowsEl {
		if index == 0 {
			continue
		}
		elRow := rowsEl[index]
		enRow := rowsEn[index]
		bilingualRowData := BilingualRowData{
			Id:  elRow[8],
			Row: index,
			El:  parseRow(elRow, elRow),
			En:  parseRow(enRow, elRow),
		}
		response[bilingualRowData.Id] = bilingualRowData
	}
	return response

}

func parseRow(row []string, elRow []string) RowData {

	rowData := RowData{
		Levels: make([]string, 0),
	}
	count := 0
	for i := 0; i < 7; i++ {
		if len(row[i]) > 0 && row[i] != "" {
			rowData.Levels = append(rowData.Levels, ToGreeklish(elRow[i]))
			count++
		}
	}

	if len(row[7]) > 0 && row[7] == "Yes" {
		rowData.IsBone = true
	}
	rowData.FileName = row[8]
	if count > 0 {
		rowData.IdName = ToGreeklish(rowData.Levels[count-1])
	} else {
		rowData.IdName = rowData.FileName
	}

	return rowData
}

type BilingualRowData struct {
	Id  string
	Row int
	El  RowData
	En  RowData
}

type RowData struct {
	IdName   string
	FileName string
	Levels   []string
	IsBone   bool
}
