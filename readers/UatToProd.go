package readers

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type UatToPros struct {
	httpClient *http.Client
	host       string
}

/*

https://strapi-cms-admin-uat-strapi-cms.tstpubaks.cosmote.gr/api/business-categories?populate=parent_categories&locale=en

https://strapi-cms-admin-uat-strapi-cms.tstpubaks.cosmote.gr/api/business-pages?pagination%5BpageSize%5D=100&populate=reusables.box.button%2Creusables.button%2Cbox.button%2Ccarousel%2Cbusiness_categories%2Creusables%2Creusables.box%2Creusables.grids%2Creusables.youtubevideos%2Creusables.youtubeItem%2Creusables.grid-alert%2Creusables.accordion%2Creusables.item%2Creusables.items%2Creusables.grid-info&locale=el

*/

func NewUatToPros() *UatToPros {
	return &UatToPros{
		httpClient: &http.Client{
			Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
			Timeout:   60 * time.Second,
		},
		host: os.Getenv("STRAPI_URL"),
	}
}

func (srv *UatToPros) Migrate() {
	//srv.categories()
	srv.pages(srv.categories())
}

func (srv *UatToPros) pages(greek map[string]int,
	english map[string]int,
	all map[string]CategoriesEnEl) {
	pagesDataEn, err := ioutil.ReadFile(os.Getenv("PAGES_EN"))
	pagesDataEl, err := ioutil.ReadFile(os.Getenv("PAGES_EL"))
	if err != nil {
		log.Fatal(err)
	}
	PagesEn := Pages{}
	PagesEl := Pages{}
	err = json.Unmarshal(pagesDataEn, &PagesEn)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(pagesDataEl, &PagesEl)
	if err != nil {
		log.Fatal(err)
	}
	mapByPageId := make(map[string]PagesEnEl)

	for _, pg := range PagesEl.PagesList {
		mapByPageId[pg.Attributes.PageID] = PagesEnEl{
			El: pg,
		}
	}

	for _, pg := range PagesEn.PagesList {
		if v, exists := mapByPageId[pg.Attributes.PageID]; exists {
			v.En = pg
			mapByPageId[pg.Attributes.PageID] = v
		} else {
			mapByPageId[pg.Attributes.PageID] = PagesEnEl{
				En: pg,
			}
		}
	}

	for k, page := range mapByPageId {
		//fmt.Println(k, " : ", page.El.Attributes.PageID, " : ", page.En.Attributes.PageID)
		insert := srv.pageToInsertablePage(page.El, greek)
		elId := srv.InsertPage(insert, k)

		if len(page.En.Attributes.PageID) > 0 {
			insrtEm := srv.pageToInsertablePage(page.En, english)
			srv.LocalizePage(insrtEm, elId, k)
		}
	}
}

func (srv *UatToPros) pageToInsertablePage(page Page, categories map[string]int) InsertPage {
	insertable := InsertPage{}
	insertable.Data.Title = page.Attributes.Title
	insertable.Data.Locale = page.Attributes.Locale
	insertable.Data.IsBusinessOne = page.Attributes.IsBusinessOne
	insertable.Data.HasContactUs = page.Attributes.HasContactUs
	insertable.Data.PageTemplate = page.Attributes.PageTemplate
	insertable.Data.PageID = page.Attributes.PageID
	if len(page.Attributes.Carousel) > 0 {
		insertable.Data.Carousel = make([]InsertPageCarousel, 0)
		insertPageCarousel := InsertPageCarousel{
			Text:              page.Attributes.Carousel[0].Text,
			ButtonLabel:       page.Attributes.Carousel[0].ButtonLabel,
			TypographyVariant: page.Attributes.Carousel[0].TypographyVariant,
			TypographyColor:   page.Attributes.Carousel[0].TypographyColor,
			LabelNext:         page.Attributes.Carousel[0].LabelNext,
			LabelPrev:         page.Attributes.Carousel[0].LabelPrev,
			IconButtonSize:    page.Attributes.Carousel[0].IconButtonSize,
			IconButtonVariant: page.Attributes.Carousel[0].IconButtonVariant,
			Href:              page.Attributes.Carousel[0].Href,
			MigratedImageURL:  page.Attributes.Carousel[0].MigratedImageURL,
		}
		insertable.Data.Carousel = append(insertable.Data.Carousel, insertPageCarousel)
	}
	if len(page.Attributes.BusinessCategories.Data) > 0 {
		categoryIds := make([]int, 0)
		for _, v := range page.Attributes.BusinessCategories.Data {
			if id, exists := categories[v.Attributes.CategoryID]; exists {
				categoryIds = append(categoryIds, id)
			}
		}
		if len(categoryIds) > 0 {
			insertable.Data.BusinessCategories = categoryIds
		}
	}

	if len(page.Attributes.Reusables) > 0 {
		reusables := make([]InsertableReusable, 0)
		for _, reuse := range page.Attributes.Reusables {
			insertableReusable := InsertableReusable{}
			insertableReusable.Title = reuse.Title
			insertableReusable.Component = reuse.Component
			insertableReusable.SecondBody = reuse.SecondBody
			insertableReusable.Template = reuse.Template
			insertableReusable.Position = reuse.Position
			insertableReusable.Name = reuse.Name
			insertableReusable.Body = reuse.Body
			insertableReusable.SecondTitle = reuse.SecondTitle
			if len(reuse.Items) > 0 {
				items := make([]InsertableItem, 0)
				for _, item := range reuse.Items {
					ini := InsertableItem{
						Title: item.Title,
						Body:  item.Body,
					}
					items = append(items, ini)
				}
				insertableReusable.Items = items
			}

			if len(reuse.Grids) > 0 {
				insGrids := make([]InsertableGrid, 0)
				for _, grid := range reuse.Grids {
					insGrid := InsertableGrid{
						Title:            grid.Title,
						Description:      grid.Description,
						URL:              grid.URL,
						ButtonText:       grid.ButtonText,
						MigratedImageURL: grid.MigratedImageURL,
					}
					insGrids = append(insGrids, insGrid)
				}
				insertableReusable.Grids = insGrids
			}

			if len(reuse.YoutubeItem) > 0 {
				ytims := make([]InsertableYouTubeItem, 0)
				for _, yti := range reuse.YoutubeItem {
					ytItm := InsertableYouTubeItem{
						VideoID:     yti.VideoID,
						Description: yti.Description,
					}
					ytims = append(ytims, ytItm)
				}
				insertableReusable.YoutubeItem = ytims
			}

			if len(reuse.Box) > 0 {
				insb := make([]InsertableBox, 0)
				for _, box := range reuse.Box {
					ins := InsertableBox{
						Body:             box.Body,
						MigratedImageURL: box.MigratedImageURL,
					}
					if len(box.Button) > 0 {
						insbtns := make([]InsertableButton, 0)
						for _, insBtn := range box.Button {
							bu := InsertableButton{
								Title: insBtn.Title,
								URL:   insBtn.Url,
							}
							insbtns = append(insbtns, bu)
						}

						ins.Button = insbtns
					}
					insb = append(insb, ins)
				}
				insertableReusable.Box = insb
			}
			reusables = append(reusables, insertableReusable)
		}
		insertable.Data.Reusables = reusables
	}

	return insertable
}

func (srv *UatToPros) InsertPage(body InsertPage, key string) int {
	jsonData, _ := json.Marshal(body)
	if key == "more_info_xlarge_4" {
		fmt.Println("Insert", string(jsonData))
	}

	uri, _ := url.JoinPath(srv.host, "api/business-pages")
	req, _ := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	var resp InsertPageResponse
	response, _ := srv.httpClient.Do(req)
	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(responseBody, &resp)
	return resp.InsertPageResponseData.ID
}

func (srv *UatToPros) LocalizePage(body InsertPage, id int, key string) int {
	jsonData, _ := json.Marshal(body.Data)
	if key == "more_info_xlarge_4" {
		fmt.Println("Insert", string(jsonData))
	}
	uri, _ := url.JoinPath(srv.host, "api/business-pages", fmt.Sprintf("%d", id), "localizations")
	req, _ := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	var resp InsertPageResponseData
	response, _ := srv.httpClient.Do(req)
	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(responseBody, &resp)
	return resp.ID
}

type InsertPage struct {
	Data InsertableData `json:"data,omitempty"`
}

type InsertableData struct {
	Title              string               `json:"title,omitempty"`
	PageID             string               `json:"pageId,omitempty"`
	PageTemplate       string               `json:"pageTemplate,omitempty"`
	BusinessCategories []int                `json:"business_categories,omitempty"`
	HasContactUs       bool                 `json:"hasContactUs,omitempty"`
	IsBusinessOne      bool                 `json:"isBusinessOne,omitempty"`
	Locale             string               `json:"locale,omitempty"`
	Reusables          []InsertableReusable `json:"reusables,omitempty"`
	Carousel           []InsertPageCarousel `json:"carousel,omitempty"`
}
type InsertableReusable struct {
	Component   string                  `json:"__component,omitempty"`
	Title       string                  `json:"title,omitempty"`
	Body        string                  `json:"body,omitempty"`
	Template    string                  `json:"template,omitempty"`
	Name        string                  `json:"name,omitempty"`
	Position    string                  `json:"position,omitempty"`
	SecondTitle string                  `json:"secondTitle,omitempty"`
	SecondBody  string                  `json:"secondBody,omitempty"`
	Grids       []InsertableGrid        `json:"grids,omitempty"`
	YoutubeItem []InsertableYouTubeItem `json:"youtubeItem,omitempty"`
	Items       []InsertableItem        `json:"items,omitempty"`
	Box         []InsertableBox         `json:"box,omitempty"`
}
type InsertableBox struct {
	Body             string             `json:"body,omitempty"`
	Button           []InsertableButton `json:"button,omitempty"`
	MigratedImageURL string             `json:"migratedImageUrl,omitempty"`
}
type InsertableButton struct {
	Title string `json:"title,omitempty"`
	URL   string `json:"url,omitempty"`
}
type InsertableItem struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}
type InsertableYouTubeItem struct {
	VideoID     string `json:"videoID,omitempty"`
	Description string `json:"description,omitempty"`
}
type InsertableGrid struct {
	Title            string `json:"title,omitempty"`
	Description      string `json:"description,omitempty"`
	URL              string `json:"url,omitempty"`
	ButtonText       string `json:"buttonText,omitempty"`
	MigratedImageURL string `json:"migratedImageUrl,omitempty"`
}
type InsertPageCarousel struct {
	Text              string `json:"text,omitempty"`
	ButtonLabel       string `json:"button_label,omitempty"`
	TypographyVariant string `json:"typographyVariant,omitempty"`
	TypographyColor   string `json:"typographyColor,omitempty"`
	LabelNext         string `json:"labelNext,omitempty"`
	LabelPrev         string `json:"labelPrev,omitempty"`
	IconButtonSize    string `json:"iconButtonSize,omitempty"`
	IconButtonVariant string `json:"iconButtonVariant,omitempty"`
	Href              string `json:"href,omitempty"`
	MigratedImageURL  string `json:"migratedImageUrl,omitempty"`
}

type InsertPageResponse struct {
	InsertPageResponseData InsertPageResponseData `json:"data"`
}
type InsertPageResponseData struct {
	ID int `json:"id"`
}

type Button struct {
	Title string `json:"title,omitempty"`
	Url   string `json:"url,omitempty"`
}

type Pages struct {
	PagesList []Page `json:"data"`
}

type Page struct {
	ID         int            `json:"-"`
	Attributes PageAttributes `json:"attributes,omitempty"`
}

type PageAttributes struct {
	Title              string             `json:"title,omitempty"`
	PageID             string             `json:"pageId,omitempty"`
	PageTemplate       string             `json:"pageTemplate,omitempty"`
	Locale             string             `json:"locale,omitempty"`
	HasContactUs       bool               `json:"hasContactUs,omitempty"`
	IsBusinessOne      bool               `json:"isBusinessOne,omitempty"`
	Carousel           []Carousel         `json:"carousel,omitempty"`
	BusinessCategories BusinessCategories `json:"business_categories,omitempty"`
	Reusables          []Reusable         `json:"reusables,omitempty"`
}
type BusinessCategories struct {
	Data []BusinessCategoriesData `json:"data,omitempty"`
}
type BusinessCategoriesData struct {
	Attributes Attributes `json:"attributes,omitempty"`
}
type Attributes struct {
	CategoryID string `json:"categoryId,omitempty"`
}
type Reusable struct {
	Component   string        `json:"__component,omitempty"`
	Title       string        `json:"title,omitempty"`
	Body        string        `json:"body,omitempty"`
	Template    string        `json:"template,omitempty"`
	SecondTitle string        `json:"secondTitle,omitempty"`
	SecondBody  string        `json:"secondBody,omitempty"`
	Name        string        `json:"name,omitempty"`
	Position    string        `json:"position,omitempty"`
	Box         []Box         `json:"box,omitempty"`
	Grids       []Grid        `json:"grids,omitempty"`
	Items       []Item        `json:"items,omitempty"`
	YoutubeItem []YoutubeItem `json:"youtubeItem,omitempty"`
}
type Box struct {
	Body             string   `json:"body,omitempty"`
	MigratedImageURL string   `json:"migratedImageUrl,omitempty"`
	Button           []Button `json:"button,omitempty"`
}
type YoutubeItem struct {
	VideoID     string `json:"videoID,omitempty"`
	Description string `json:"description,omitempty"`
}
type Item struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}
type Grid struct {
	Title            string `json:"title,omitempty"`
	Description      string `json:"description,omitempty"`
	URL              string `json:"url,omitempty"`
	ButtonText       string `json:"buttonText,omitempty"`
	MigratedImageURL string `json:"migratedImageUrl,omitempty"`
}

func (srv *UatToPros) InsertCategory(body InsertCategory) int {
	jsonData, _ := json.Marshal(body)
	uri, _ := url.JoinPath(srv.host, "api/business-categories")
	req, _ := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	var resp InsertCategoryResponse
	response, _ := srv.httpClient.Do(req)
	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(responseBody, &resp)
	return resp.Data.ID

}

func (srv *UatToPros) LocalizeCategory(body InsertCategory, id int) int {
	jsonData, _ := json.Marshal(body.Data)
	uri, _ := url.JoinPath(srv.host, "api/business-categories", fmt.Sprintf("%d", id), "localizations")
	req, _ := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	var resp InsertCategoryResponseData
	response, _ := srv.httpClient.Do(req)
	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(responseBody, &resp)
	return resp.ID
}

func (srv *UatToPros) SetParent(body InsertCategory, id int) {
	jsonData, _ := json.Marshal(body)
	uri, _ := url.JoinPath(srv.host, "api/business-categories", fmt.Sprintf("%d", id))
	req, _ := http.NewRequest(http.MethodPut, uri, bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	response, _ := srv.httpClient.Do(req)
	defer response.Body.Close()

}

func (srv *UatToPros) categories() (map[string]int, map[string]int, map[string]CategoriesEnEl) {
	categoriesDataEn, err := ioutil.ReadFile(os.Getenv("CATEGORIES_EN"))
	categoriesDataEl, err := ioutil.ReadFile(os.Getenv("CATEGORIES_EL"))
	if err != nil {
		log.Fatal(err)
	}
	CategoriesEn := Categories{}
	CategoriesEl := Categories{}
	err = json.Unmarshal(categoriesDataEn, &CategoriesEn)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(categoriesDataEl, &CategoriesEl)
	if err != nil {
		log.Fatal(err)
	}

	mapByCategoryIdEn := make(map[string]int)
	mapByCategoryIdEl := make(map[string]int)
	mapByCategoryId := make(map[string]CategoriesEnEl)
	for _, cat := range CategoriesEl.CategoriesList {
		cat.ID = srv.InsertCategory(InsertCategory{Data: InsertCategoryData{
			Title:      cat.Attributes.Title,
			URL:        cat.Attributes.URL,
			CategoryID: cat.Attributes.CategoryID,
			Locale:     cat.Attributes.Locale,
		}})
		mapByCategoryId[cat.Attributes.CategoryID] = CategoriesEnEl{El: cat}
		mapByCategoryIdEl[cat.Attributes.CategoryID] = cat.ID
	}
	for _, cat := range CategoriesEn.CategoriesList {
		v, exists := mapByCategoryId[cat.Attributes.CategoryID]
		if exists {
			cat.ID = srv.LocalizeCategory(InsertCategory{Data: InsertCategoryData{
				Title:      cat.Attributes.Title,
				URL:        cat.Attributes.URL,
				CategoryID: cat.Attributes.CategoryID,
				Locale:     cat.Attributes.Locale,
			}}, v.El.ID)
			v.En = cat
			mapByCategoryId[cat.Attributes.CategoryID] = v
			mapByCategoryIdEn[cat.Attributes.CategoryID] = cat.ID
		} else {
			cat.ID = srv.InsertCategory(InsertCategory{Data: InsertCategoryData{
				Title:      cat.Attributes.Title,
				URL:        cat.Attributes.URL,
				CategoryID: cat.Attributes.CategoryID,
				Locale:     cat.Attributes.Locale,
			}})
			mapByCategoryId[cat.Attributes.CategoryID] = CategoriesEnEl{En: cat}
			mapByCategoryIdEn[cat.Attributes.CategoryID] = cat.ID
		}
	}

	//PARENTING
	for _, v := range mapByCategoryId {
		if len(v.El.Attributes.ParentCategories.ParentCategoriesList) > 0 {
			parentsEl := make([]int, 0)
			for _, parent := range v.El.Attributes.ParentCategories.ParentCategoriesList {
				parentsEl = append(parentsEl, mapByCategoryIdEl[parent.Attributes.CategoryID])
			}
			srv.SetParent(InsertCategory{Data: InsertCategoryData{ParentCategories: parentsEl}}, v.El.ID)
		}

		if len(v.En.Attributes.ParentCategories.ParentCategoriesList) > 0 {
			parentsEn := make([]int, 0)
			for _, parent := range v.En.Attributes.ParentCategories.ParentCategoriesList {
				parentsEn = append(parentsEn, mapByCategoryIdEn[parent.Attributes.CategoryID])
			}
			srv.SetParent(InsertCategory{Data: InsertCategoryData{ParentCategories: parentsEn}}, v.En.ID)
		}
	}
	return mapByCategoryIdEn, mapByCategoryIdEl, mapByCategoryId
}

type InsertCategory struct {
	Data InsertCategoryData `json:"data,omitempty"`
}
type InsertCategoryData struct {
	Title            string `json:"title,omitempty"`
	URL              string `json:"url,omitempty"`
	CategoryID       string `json:"categoryId,omitempty"`
	Locale           string `json:"locale,omitempty"`
	ParentCategories []int  `json:"parent_categories,omitempty"`
}
type InsertCategoryResponse struct {
	Data InsertCategoryResponseData `json:"data,omitempty"`
}
type InsertCategoryResponseData struct {
	ID int `json:"id,omitempty"`
}

type PagesEnEl struct {
	En Page
	El Page
}

type CategoriesEnEl struct {
	En Category
	El Category
}

type Categories struct {
	CategoriesList []Category `json:"data,omitempty"`
}

type Category struct {
	ID         int `json:"-"`
	Attributes struct {
		Title            string `json:"title"`
		URL              string `json:"url"`
		Locale           string `json:"locale"`
		CategoryID       string `json:"categoryId"`
		ParentCategories struct {
			ParentCategoriesList []struct {
				Attributes struct {
					CategoryID string `json:"categoryId"`
				} `json:"attributes"`
			} `json:"data"`
		} `json:"parent_categories"`
	} `json:"attributes"`
}

type Carousel struct {
	Text              string `json:"text"`
	ButtonLabel       string `json:"button_label"`
	TypographyVariant string `json:"typographyVariant"`
	TypographyColor   string `json:"typographyColor"`
	LabelNext         string `json:"labelNext"`
	LabelPrev         string `json:"labelPrev"`
	IconButtonSize    string `json:"iconButtonSize"`
	IconButtonVariant string `json:"iconButtonVariant"`
	Href              string `json:"href"`
	MigratedImageURL  string `json:"migratedImageUrl"`
	SubTitle          string `json:"subTitle"`
}
