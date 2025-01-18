package readers

import (
	"app/adapter"
	"app/model"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"strings"
)

type CustomHtmlReader struct {
	FilePath      string
	StrapiUrl     string
	ImagePath     string
	StrapiAdapter *adapter.StrapiAdapter
}

func NewCustomHtmlReaderReader() *CustomHtmlReader {
	return &CustomHtmlReader{
		FilePath:      os.Getenv("CUSTOM_HTML_PATH"),
		StrapiUrl:     os.Getenv("STRAPI_URL"),
		ImagePath:     os.Getenv("IMAGE_PATH"),
		StrapiAdapter: adapter.NewStrapiAdapter(),
	}
}

/*func (reader *CustomHtmlReader) Execute() {
	files, err := ioutil.ReadDir(reader.FilePath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		reader.Splitter(file.Name())
	}
}*/

func (reader *CustomHtmlReader) Splitter(insertable model.Insertable) {
	pages := reader.HandleContent(reader.ReadExcel(insertable.FilePath), insertable)
	grResp := reader.StrapiAdapter.Insert(model.BusinessPageInsert{Data: *pages.GrPage})
	reader.StrapiAdapter.Localizations(*pages.EnPage, grResp.Data.Id)
}

func (reader *CustomHtmlReader) HandleContent(excl model.Excelized, insertable model.Insertable) model.BusinessPageWrapper {

	//name = strings.TrimSpace(name)
	//names := strings.Split(name, ".")

	reusables := make([]model.Reusable, 0)
	reusables = append(reusables, reader.handleReusableHtml(excl.GR.Content))
	businessPageWrapper := model.BusinessPageWrapper{}
	businessPage := &model.BusinessPage{
		Title:              excl.GR.Title,
		PageID:             insertable.PageId,
		BusinessCategories: insertable.CategoriesEl,
		PageTemplate:       "blank",
		Locale:             "el",
		Reusables:          reusables,
	}
	businessPageWrapper.GrPage = businessPage

	reusablesEn := make([]model.Reusable, 0)
	reusablesEn = append(reusablesEn, reader.handleReusableHtml(excl.EN.Content))
	businessPageEn := &model.BusinessPage{
		Title:              excl.EN.Title,
		PageID:             insertable.PageId,
		BusinessCategories: insertable.CategoriesEn,
		PageTemplate:       "blank",
		Locale:             "en",
		Reusables:          reusablesEn,
	}
	businessPageWrapper.EnPage = businessPageEn
	return businessPageWrapper

}

func (reader *CustomHtmlReader) ReadExcelContent(path string) model.Content {
	f, _ := excelize.OpenFile(reader.FilePath + "\\\\" + path)
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal("Could mot close file", path)
		}
	}()

	content := make([]string, 0)
	contentEn := make([]string, 0)
	rows, errGr := f.GetRows("GR")
	rowsEn, errEn := f.GetRows("EN")
	if errGr == nil {
		for index, row := range rows {
			if index == 0 {
				continue
			}
			if index == 1 {
				content = append(content, valueStringOtNull(row[4]))
			}
			if index > 1 {
				break
			}

		}
	}

	if errEn == nil {
		for index, row := range rowsEn {
			if index == 0 {
				continue
			}
			if index == 1 {
				contentEn = append(contentEn, valueStringOtNull(row[4]))
			}
			if index > 1 {
				break
			}

		}
	}

	return model.Content{
		GR: strings.Join(content, " "),
		EN: strings.Join(contentEn, " "),
	}

}

func (reader *CustomHtmlReader) ReadExcel(path string) model.Excelized {
	f, _ := excelize.OpenFile(reader.FilePath + "\\\\" + path)
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal("Could mot close file", path)
		}
	}()
	res := model.Excelized{
		GR: model.ExItem{},
		EN: model.ExItem{},
	}
	rows, errGr := f.GetRows("GR")
	rowsEn, errEn := f.GetRows("EN")
	if errGr == nil {
		for index, row := range rows {
			if index == 0 {
				continue
			}
			/*if index == 1 {
				res.GR.Title = row[4]
			}*/
			if index == 1 {
				res.GR.Content = row[4]
			}

			if index > 1 {
				break
			}
		}
	}

	if errEn == nil {
		for index, row := range rowsEn {
			if index == 0 {
				continue
			}
			/*if index == 1 {
				res.EN.Title = row[4]
			}*/
			if index == 1 {
				res.EN.Content = row[4]
			}

			if index > 1 {
				break
			}
		}
	}

	return res
}

func (reader *CustomHtmlReader) handleReusableHtml(content string) model.Reusable {

	reusable := model.Reusable{
		Component: "reusables.html",
		Title:     "",
		Body:      content,
	}
	return reusable
}
