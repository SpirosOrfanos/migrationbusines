package readers

import (
	"app/adapter"
	"app/model"
	"fmt"
	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"strconv"
	"strings"
)

type FileReader struct {
	FilePath      string
	StrapiUrl     string
	ImagePath     string
	StrapiAdapter *adapter.StrapiAdapter
}

func NewFileReader() *FileReader {
	return &FileReader{
		FilePath:      os.Getenv("FILE_PATH"),
		StrapiUrl:     os.Getenv("STRAPI_URL"),
		ImagePath:     os.Getenv("IMAGE_PATH"),
		StrapiAdapter: adapter.NewStrapiAdapter(),
	}
}

func (reader *FileReader) Splitter(insertable model.Insertable) {
	fmt.Println(insertable.FilePath)
	pages := reader.HandleContent(reader.ReadExcel(insertable.FilePath), insertable)
	grResp := reader.StrapiAdapter.Insert(model.BusinessPageInsert{Data: *pages.GrPage})

	reader.StrapiAdapter.Localizations(*pages.EnPage, grResp.Data.Id)
}

func (reader *FileReader) HandleContent(ch2 []model.Excelized, insertable model.Insertable) model.BusinessPageWrapper {
	sections := make(map[int]*sll.List)

	for _, excItem := range ch2 {
		v, ok := sections[excItem.GR.Order]
		if ok {
			v.Add(excItem)
			sections[excItem.GR.Order] = v
		} else {
			list := sll.New()
			list.Add(excItem)
			sections[excItem.GR.Order] = list
		}
	}

	reusables := make([]model.Reusable, 0)
	carouselImageIf := sections[0].Values()[0]

	carouselImage := carouselImageIf.(model.Excelized)
	carousels := make([]model.Carousel, 1)
	carousels[0] = model.Carousel{
		Text:             "",
		MigratedImageURL: fmt.Sprintf("%s%s", reader.ImagePath, carouselImage.GR.Content),
	}

	if len(sections[0].Values()) > 1 {
		carouselImageIf2 := sections[0].Values()[1]
		carouselImage2 := carouselImageIf2.(model.Excelized)
		carousels[0].Text = carouselImage2.GR.Content
	}
	businessPageWrapper := model.BusinessPageWrapper{}
	businessPage := &model.BusinessPage{
		Title:              insertable.TitleEl,
		PageID:             insertable.PageId,
		PageTemplate:       "default",
		BusinessCategories: insertable.CategoriesEl,
		Locale:             "el",
		Carousel:           carousels,
		IsBusinessOne:      insertable.IsBusiness,
	}

	for _, item := range sections {
		plann := item.Values()[0]
		excl := plann.(model.Excelized)
		if excl.GR.ReusableType == "reusable-html" {
			reusables = append(reusables, reader.handleReusableHtml(item, true))
		}

		if excl.GR.ReusableType == "reusable-accordion-item" {
			reusables = append(reusables, reader.handleAccordion(item, true))
		}

		if excl.GR.ReusableType == "reusable-contact" {
			reusables = append(reusables, reader.handleContact(item, true))
		}

		if excl.GR.ReusableType == "reusable-video" {
			reusables = append(reusables, handleVideo(item, true))
		}
		if excl.GR.ReusableType == "reusable-grids" {
			reusables = append(reusables, handleReusableGrids(item, true))
		}
		if excl.GR.ReusableType == "reusable-grid-alert" {
			reusables = append(reusables, handleReusableGridAlert(item, true))
		}

		if excl.GR.ReusableType == "reusable-testimonial" {
			reusables = append(reusables, handleReusableTestimonial(item, true))
		}

		if excl.GR.ReusableType == "hasContactUs" {
			withContact := hasContact(item, false)
			businessPage.HasContactUs = withContact
		}

	}

	businessPage.Reusables = reusables
	businessPageWrapper.GrPage = businessPage

	reusablesEn := make([]model.Reusable, 0)

	carouselsEn := make([]model.Carousel, 1)
	carouselsEn[0] = model.Carousel{
		Text:             "",
		MigratedImageURL: fmt.Sprintf("%s%s", reader.ImagePath, carouselImage.EN.Content),
	}

	if len(sections[0].Values()) > 1 {
		carouselImageIf2 := sections[0].Values()[1]
		carouselImage2 := carouselImageIf2.(model.Excelized)
		carouselsEn[0].Text = carouselImage2.EN.Content
	}

	businessPageEn := &model.BusinessPage{
		Title:              insertable.TitleEn,
		PageID:             insertable.PageId,
		PageTemplate:       "default",
		BusinessCategories: insertable.CategoriesEl,
		Locale:             "en",
		Carousel:           carouselsEn,
		IsBusinessOne:      insertable.IsBusiness,
	}

	for _, item := range sections {
		plann := item.Values()[0]
		excl := plann.(model.Excelized)
		if excl.EN.ReusableType == "reusable-html" {
			reusablesEn = append(reusablesEn, reader.handleReusableHtml(item, false))

		}
		if excl.EN.ReusableType == "reusable-accordion-item" {
			reusablesEn = append(reusablesEn, reader.handleAccordion(item, false))
		}
		if excl.EN.ReusableType == "reusable-contact" {
			reusablesEn = append(reusablesEn, reader.handleContact(item, false))
		}
		if excl.EN.ReusableType == "reusable-video" {
			reusablesEn = append(reusablesEn, handleVideo(item, false))
		}
		if excl.EN.ReusableType == "reusable-grids" {
			reusablesEn = append(reusablesEn, handleReusableGrids(item, false))
		}
		if excl.EN.ReusableType == "reusable-grid-alert" {
			reusablesEn = append(reusablesEn, handleReusableGridAlert(item, false))
		}
		if excl.GR.ReusableType == "reusable-testimonial" {
			reusablesEn = append(reusablesEn, handleReusableTestimonial(item, false))
		}

		if excl.EN.ReusableType == "hasContactUs" {
			withContact := hasContact(item, false)
			businessPageEn.HasContactUs = withContact
		}
	}

	businessPageEn.Reusables = reusablesEn
	businessPageWrapper.EnPage = businessPageEn
	return businessPageWrapper

}
func (reader *FileReader) handleReusableHtml(items *sll.List, isGr bool) model.Reusable {
	reusable := model.Reusable{
		Component: "reusables.html",
		Title:     "",
		Body:      "",
	}
	if isGr {
		for _, item := range items.Values() {
			itm := item.(model.Excelized)
			if itm.GR.TypeOfRecord == "html" {
				reusable.Body = itm.GR.Content
			}
			if itm.GR.TypeOfRecord == "title" {
				reusable.Title = itm.GR.Content
			}
		}
	} else {
		for _, item := range items.Values() {
			itm := item.(model.Excelized)
			if itm.EN.TypeOfRecord == "html" {
				reusable.Body = itm.EN.Content
			}
			if itm.EN.TypeOfRecord == "title" {
				reusable.Title = itm.EN.Content
			}
		}
	}

	return reusable
}

func (reader *FileReader) handleAccordion(items *sll.List, isGr bool) model.Reusable {
	reusable := model.Reusable{
		Component: "reusables.accordion",
		Title:     "",
		Body:      "",
		Items:     make([]model.Item, 0),
	}
	if isGr {
		for _, item := range items.Values() {
			itm := item.(model.Excelized)
			if itm.GR.TypeOfRecord == "html" || itm.GR.TypeOfRecord == "html " {
				reusable.Body = itm.GR.Content
				continue
			}
			if itm.GR.TypeOfRecord == "title" || itm.GR.TypeOfRecord == "title " {
				reusable.Title = itm.GR.Content
				continue
			}
			if itm.GR.TypeOfRecord == "title accordion" || itm.GR.TypeOfRecord == "title-accordion" ||
				itm.GR.TypeOfRecord == "accordion title" || itm.GR.TypeOfRecord == "accordion-title" {
				reusable.Items = append(reusable.Items, model.Item{
					Title: itm.GR.Content,
				})
				continue
			}
			if itm.GR.TypeOfRecord == "html accordion" || itm.GR.TypeOfRecord == "html-accordion" ||
				itm.GR.TypeOfRecord == "accordion html" || itm.GR.TypeOfRecord == "accordion-html" {
				reusable.Items[len(reusable.Items)-1].Body = itm.GR.Content
			}
		}
	} else {
		for _, item := range items.Values() {
			itm := item.(model.Excelized)
			if itm.EN.TypeOfRecord == "html" || itm.EN.TypeOfRecord == "html " {
				reusable.Body = itm.EN.Content
				continue
			}
			if itm.EN.TypeOfRecord == "title" || itm.EN.TypeOfRecord == "title " {
				reusable.Title = itm.EN.Content
				continue
			}
			if itm.EN.TypeOfRecord == "title accordion" || itm.EN.TypeOfRecord == "title-accordion" ||
				itm.EN.TypeOfRecord == "accordion title" || itm.EN.TypeOfRecord == "accordion-title" ||
				itm.EN.TypeOfRecord == "title accordion " || itm.EN.TypeOfRecord == "title-accordion " ||
				itm.EN.TypeOfRecord == "accordion title " || itm.EN.TypeOfRecord == "accordion-title " {
				reusable.Items = append(reusable.Items, model.Item{
					Title: itm.EN.Content,
				})
				continue
			}
			if itm.EN.TypeOfRecord == "html accordion" || itm.EN.TypeOfRecord == "html-accordion" ||
				itm.EN.TypeOfRecord == "accordion html" || itm.EN.TypeOfRecord == "accordion-html" ||
				itm.EN.TypeOfRecord == "html accordion " || itm.EN.TypeOfRecord == "html-accordion " ||
				itm.EN.TypeOfRecord == "accordion html " || itm.EN.TypeOfRecord == "accordion-html " ||
				itm.EN.TypeOfRecord == "body accordion " || itm.EN.TypeOfRecord == "body-accordion " ||
				itm.EN.TypeOfRecord == "accordion body " || itm.EN.TypeOfRecord == "accordion-body " {
				reusable.Items[len(reusable.Items)-1].Body = itm.EN.Content
				continue
			}
		}
	}

	return reusable
}

func (reader *FileReader) handleContact(items *sll.List, isGr bool) model.Reusable {
	reusable := model.Reusable{
		Component: "reusables.contact-box",
		Title:     "",
		Body:      "",
	}
	boxFound := false
	if isGr {
		for _, item := range items.Values() {
			itm := item.(model.Excelized)
			if (itm.GR.TypeOfRecord == "text" || itm.GR.TypeOfRecord == "title") && !boxFound {
				reusable.Title = itm.GR.Content
				boxFound = true
				reusable.Box = make([]model.Box, 0)
				continue
			}
			if itm.GR.TypeOfRecord == "text" || itm.GR.TypeOfRecord == "text " {
				reusable.Box = append(reusable.Box, model.Box{Body: itm.GR.Content})
				continue
			}

			if itm.GR.TypeOfRecord == "img" || itm.GR.TypeOfRecord == "icon" ||
				itm.GR.TypeOfRecord == "img " || itm.GR.TypeOfRecord == "icon " {
				reusable.Box[len(reusable.Box)-1].MigratedImageUrl = fmt.Sprintf("%s%s", reader.ImagePath, itm.GR.Content)
				continue
			}
			if itm.GR.TypeOfRecord == "button" || itm.GR.TypeOfRecord == "button " {
				reusable.Box[len(reusable.Box)-1].Button = make([]model.Button, 1)
				reusable.Box[len(reusable.Box)-1].Button[0].Url = itm.GR.Content
				continue
			}

			if strings.Contains(itm.GR.TypeOfRecord, "name") || strings.Contains(itm.GR.TypeOfRecord, "name ") {
				reusable.Box[len(reusable.Box)-1].Button[0].Title = itm.GR.Content
				continue
			}
		}
	} else {
		for _, item := range items.Values() {
			itm := item.(model.Excelized)
			if (itm.EN.TypeOfRecord == "text" || itm.EN.TypeOfRecord == "title" ||
				itm.EN.TypeOfRecord == "text " || itm.EN.TypeOfRecord == "title ") && !boxFound {
				reusable.Title = itm.EN.Content
				boxFound = true
				reusable.Box = make([]model.Box, 0)
				continue
			}
			if itm.EN.TypeOfRecord == "text" || itm.EN.TypeOfRecord == "text " {
				reusable.Box = append(reusable.Box, model.Box{Body: itm.EN.Content})
				continue
			}

			if itm.EN.TypeOfRecord == "img" || itm.EN.TypeOfRecord == "icon" ||
				itm.EN.TypeOfRecord == "img " || itm.EN.TypeOfRecord == "icon " {
				reusable.Box[len(reusable.Box)-1].MigratedImageUrl = fmt.Sprintf("%s%s", reader.ImagePath, itm.EN.Content)
				continue

			}
			if itm.EN.TypeOfRecord == "button" || itm.EN.TypeOfRecord == "button " {
				reusable.Box[len(reusable.Box)-1].Button = make([]model.Button, 1)
				reusable.Box[len(reusable.Box)-1].Button[0].Url = itm.EN.Content
				continue
			}

			if strings.Contains(itm.EN.TypeOfRecord, "name") ||
				strings.Contains(itm.EN.TypeOfRecord, "name ") {
				reusable.Box[len(reusable.Box)-1].Button[0].Title = itm.EN.Content
				continue
			}
		}
	}

	return reusable
}

func hasContact(items *sll.List, isGr bool) bool {
	item := items.Values()[0]
	itm := item.(model.Excelized)
	return itm.GR.Content == "Yes" || itm.GR.Content == "YES" ||
		itm.GR.Content == "Yes " || itm.GR.Content == "YES "

}

func handleReusableTestimonial(items *sll.List, isGr bool) model.Reusable {
	reusable := model.Reusable{
		Component: "reusables.grid-info",
		Title:     "",
	}
	if isGr {
		for _, item := range items.Values() {
			itm := item.(model.Excelized)

			if itm.GR.TypeOfRecord == "quote" || itm.GR.TypeOfRecord == "quote " {
				reusable.Body = itm.GR.Content
				continue
			}
			if itm.GR.TypeOfRecord == "name" || itm.GR.TypeOfRecord == "name " {
				reusable.Name = itm.GR.Content
				continue
			}
			if itm.GR.TypeOfRecord == "position" || itm.GR.TypeOfRecord == "position " {
				reusable.Position = itm.GR.Content
				continue
			}
		}
	} else {
		for _, item := range items.Values() {
			itm := item.(model.Excelized)

			if itm.EN.TypeOfRecord == "quote" || itm.EN.TypeOfRecord == "quote" {
				reusable.Body = itm.EN.Content
				continue
			}
			if itm.EN.TypeOfRecord == "name" || itm.EN.TypeOfRecord == "name " {
				reusable.Name = itm.EN.Content
				continue
			}
			if itm.EN.TypeOfRecord == "position" || itm.EN.TypeOfRecord == "position " {
				reusable.Position = itm.EN.Content
				continue
			}
		}
	}
	return reusable
}

func handleReusableGridAlert(items *sll.List, isGr bool) model.Reusable {
	reusable := model.Reusable{
		Component: "reusables.grid-alert",
		Title:     "",
	}
	if isGr {
		for _, item := range items.Values() {
			itm := item.(model.Excelized)
			if len(reusable.Template) == 0 && (itm.GR.TypeOfRecord == "text" || itm.GR.TypeOfRecord == "title" ||
				itm.GR.TypeOfRecord == "text " || itm.GR.TypeOfRecord == "title ") {
				reusable.Title = itm.GR.Content
				reusable.Template = "default"
				reusable.Grids = make([]model.ReusableGridItem, 0)
				continue
			}
			if len(reusable.Template) > 0 && (itm.GR.TypeOfRecord == "body" || itm.GR.TypeOfRecord == "body") {
				reusable.Body = itm.GR.Content
				continue
			}
			if len(reusable.Template) > 0 && (itm.GR.TypeOfRecord == "secondTitle" || itm.GR.TypeOfRecord == "secondTitle ") {
				reusable.SecondTitle = itm.GR.Content
				continue
			}
			if len(reusable.Template) > 0 && (itm.GR.TypeOfRecord == "secondBody" || itm.GR.TypeOfRecord == "secondBody ") {
				reusable.SecondBody = itm.GR.Content
				continue
			}
		}
	} else {
		for _, item := range items.Values() {
			itm := item.(model.Excelized)
			if len(reusable.Template) == 0 && (itm.EN.TypeOfRecord == "text" || itm.EN.TypeOfRecord == "title" ||
				itm.EN.TypeOfRecord == "text " || itm.EN.TypeOfRecord == "title ") {
				reusable.Title = itm.EN.Content
				reusable.Template = "default"
				reusable.Grids = make([]model.ReusableGridItem, 0)
				continue
			}
			if len(reusable.Template) > 0 && (itm.EN.TypeOfRecord == "body" || itm.EN.TypeOfRecord == "body") {
				reusable.Body = itm.EN.Content
				continue
			}
			if len(reusable.Template) > 0 && (itm.EN.TypeOfRecord == "secondTitle" || itm.EN.TypeOfRecord == "secondTitle ") {
				reusable.SecondTitle = itm.EN.Content
				continue
			}
			if len(reusable.Template) > 0 && (itm.EN.TypeOfRecord == "secondBody" || itm.EN.TypeOfRecord == "secondBody ") {
				reusable.SecondBody = itm.EN.Content
				continue
			}
		}
	}
	return reusable
}

func handleReusableGrids(items *sll.List, isGr bool) model.Reusable {
	reusable := model.Reusable{
		Component: "reusables.grids",
		Title:     "",
	}
	if isGr {
		for _, item := range items.Values() {
			itm := item.(model.Excelized)
			if len(reusable.Template) == 0 && (itm.GR.TypeOfRecord == "text" || itm.GR.TypeOfRecord == "title" ||
				itm.GR.TypeOfRecord == "text " || itm.GR.TypeOfRecord == "title ") {
				reusable.Title = itm.GR.Content
				reusable.Template = "default"
				reusable.Grids = make([]model.ReusableGridItem, 0)
				continue
			}
			if len(reusable.Template) > 0 && (itm.GR.TypeOfRecord == "text" || itm.GR.TypeOfRecord == "title" ||
				itm.GR.TypeOfRecord == "text " || itm.GR.TypeOfRecord == "title ") {
				reusable.Grids = append(reusable.Grids, model.ReusableGridItem{
					Title: itm.GR.Content,
				})
				continue
			}
			if len(reusable.Template) > 0 && (itm.GR.TypeOfRecord == "body" || itm.GR.TypeOfRecord == "body ") {
				reusable.Grids[len(reusable.Grids)-1].Description = itm.GR.Content
				continue
			}
		}
	} else {
		for _, item := range items.Values() {
			itm := item.(model.Excelized)
			if len(reusable.Template) == 0 && (itm.EN.TypeOfRecord == "text" || itm.EN.TypeOfRecord == "title" ||
				itm.EN.TypeOfRecord == "text " || itm.EN.TypeOfRecord == "title ") {
				reusable.Title = itm.EN.Content
				reusable.Template = "default"
				reusable.Grids = make([]model.ReusableGridItem, 0)
				continue
			}
			if len(reusable.Template) > 0 && (itm.EN.TypeOfRecord == "text" || itm.EN.TypeOfRecord == "title" ||
				itm.EN.TypeOfRecord == "text " || itm.EN.TypeOfRecord == "title ") {
				reusable.Grids = append(reusable.Grids, model.ReusableGridItem{
					Title: itm.EN.Content,
				})
				continue
			}
			if len(reusable.Template) > 0 && (itm.EN.TypeOfRecord == "body" || itm.EN.TypeOfRecord == "body" ||
				itm.EN.TypeOfRecord == "body " || itm.EN.TypeOfRecord == "body ") {
				reusable.Grids[len(reusable.Grids)-1].Description = itm.EN.Content
				continue
			}
		}
	}

	return reusable
}
func handleVideo(items *sll.List, isGr bool) model.Reusable {
	reusable := model.Reusable{
		Component: "reusables.youtubevideos",
		Title:     "",
	}
	if isGr {
		for _, item := range items.Values() {
			itm := item.(model.Excelized)
			if itm.GR.TypeOfRecord == "text" || itm.GR.TypeOfRecord == "title" ||
				itm.GR.TypeOfRecord == "text " || itm.GR.TypeOfRecord == "title " {
				reusable.Title = itm.GR.Content
				reusable.YouTubeItem = make([]model.YouTubeItem, 0)
				continue
			}
			if strings.Contains(itm.GR.TypeOfRecord, "video") || strings.Contains(itm.GR.TypeOfRecord, "video ") ||
				strings.Contains(itm.GR.TypeOfRecord, "link") || strings.Contains(itm.GR.TypeOfRecord, "link ") {
				reusable.YouTubeItem = append(reusable.YouTubeItem, model.YouTubeItem{
					VideoID: itm.GR.Content,
				})
				continue
			}
		}
	} else {
		for _, item := range items.Values() {
			itm := item.(model.Excelized)
			if itm.EN.TypeOfRecord == "text" || itm.EN.TypeOfRecord == "title" ||
				itm.GR.TypeOfRecord == "text " || itm.GR.TypeOfRecord == "title " {
				reusable.Title = itm.EN.Content
				reusable.YouTubeItem = make([]model.YouTubeItem, 0)
				continue
			}
			if strings.Contains(itm.GR.TypeOfRecord, "video") || strings.Contains(itm.GR.TypeOfRecord, "video ") ||
				strings.Contains(itm.GR.TypeOfRecord, "link") || strings.Contains(itm.GR.TypeOfRecord, "link ") {
				reusable.YouTubeItem = append(reusable.YouTubeItem, model.YouTubeItem{
					VideoID: itm.EN.Content,
				})
				continue
			}
		}
	}
	return reusable
}

func (reader *FileReader) ReadExcel(path string) []model.Excelized {
	set := make([]model.Excelized, 0)
	fmt.Println(reader.FilePath + "\\\\" + path)
	f, _ := excelize.OpenFile(reader.FilePath + "\\\\" + path)
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal("Could mot close file", path)
		}
	}()
	rows, _ := f.GetRows("GR")
	rowsEn, _ := f.GetRows("EN")
	for index, row := range rows {
		if index == 0 {
			continue
		}
		var typeOfRecord = valueStringOtNull(row[3])
		var order = valueIntOrNull(row[2])
		var innerOrder = index
		var reusableType = valueStringOtNull(rows[index][1])
		var content = valueStringOtNull(row[4])

		var typeOfRecordEn = valueStringOtNull(rowsEn[index][3])
		var orderEn = valueIntOrNull(rowsEn[index][2])
		var innerOrderEn = index
		var reusableTypeEn = valueStringOtNull(rowsEn[index][1])
		var contentEn = valueStringOtNull(rowsEn[index][4])
		set = append(set, model.Excelized{
			GR: model.ExItem{
				TypeOfRecord: typeOfRecord,
				Order:        order,
				InnerOrder:   innerOrder,
				ReusableType: reusableType,
				Content:      content,
			}, EN: model.ExItem{
				TypeOfRecord: typeOfRecordEn,
				Order:        orderEn,
				InnerOrder:   innerOrderEn,
				ReusableType: reusableTypeEn,
				Content:      contentEn,
			},
		})
	}
	return set
}

func valueStringOtNull(val string) string {
	if len(val) == 0 {
		return ""
	}
	return val
}

func valueIntOrNull(val string) int {
	if len(val) == 0 {
		return 0
	}
	res, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return res
}
