package model

import "github.com/emirpasic/gods/sets"

type Excelized struct {
	GR ExItem
	EN ExItem
}

type ExItem struct {
	TypeOfRecord string
	Order        int
	InnerOrder   int
	ReusableType string
	Content      string
	Title        string
}

type NModel struct {
	Name string
}

type BusinessPageWrapper struct {
	GrId   int
	EnId   int
	GrPage *BusinessPage
	EnPage *BusinessPage
}

type BusinessPageInsert struct {
	Data BusinessPage `json:"data,omitempty"`
}
type BusinessPageResponseWrapper struct {
	Data BusinessPageResponse `json:"data,omitempty"`
}
type BusinessPageResponse struct {
	Id int `json:"id,omitempty"`
}
type BusinessPage struct {
	Title              string     `json:"title,omitempty"`
	PageID             string     `json:"pageId,omitempty"`
	PageTemplate       string     `json:"pageTemplate,omitempty"`
	BusinessCategories []int      `json:"business_categories,omitempty"`
	Locale             string     `json:"locale,omitempty"`
	Carousel           []Carousel `json:"carousel,omitempty"`
	Reusables          []Reusable `json:"reusables,omitempty"`
}

type Carousel struct {
	Text             string `json:"text,omitempty"`
	MigratedImageURL string `json:"migratedImageUrl,omitempty"`
}

type Reusable struct {
	ID          *int          `json:"id,omitempty"`
	Component   string        `json:"__component,omitempty"`
	Title       string        `json:"title,omitempty"`
	Body        string        `json:"body,omitempty"`
	Items       []Item        `json:"items,omitempty"`
	Box         []Box         `json:"box,omitempty"`
	YouTubeItem []YouTubeItem `json:"youtubeItem,omitempty"`
}

type Item struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}

type YouTubeItem struct {
	VideoID     string `json:"videoID,omitempty"`
	Description string `json:"description,omitempty"`
}

type Box struct {
	Body             string   `json:"body,omitempty"`
	MigratedImageUrl string   `json:"migratedImageUrl,omitempty"`
	Button           []Button `json:"button,omitempty"`
}

type Button struct {
	Title string `json:"title,omitempty"`
	Url   string `json:"url,omitempty"`
}

type StrapiCategoryCreateWrapper struct {
	Data StrapiCategoryCreate `json:"data,omitempty"`
}
type StrapiCategoryCreateResponseWrapper struct {
	Data StrapiCategoryCreateResponse `json:"data,omitempty"`
}

type StrapiCategoryParentingWrapper struct {
	Data StrapiCategoryParenting `json:"data,omitempty"`
}
type StrapiCategoryCreate struct {
	Title      string `json:"title,omitempty"`
	CategoryId string `json:"categoryId,omitempty"`
	Locale     string `json:"locale,omitempty"`
}
type StrapiCategoryCreateResponse struct {
	Id int `json:"id,omitempty"`
}

type StrapiCategoryParenting struct {
	ParentCategories []int `json:"parent_categories,omitempty"`
}

type Levels struct {
	Level0Map map[string]Category
	Level1Map map[string]Category
	Level2Map map[string]Category
	Level3Map map[string]Category
	Level4Map map[string]Category
	Level5Map map[string]Category
}

type Category struct {
	IdEl   int
	IdEn   int
	NameEl string
	NameEn string
	CatId  string
	Parent sets.Set
}

type Insertable struct {
	PageId       string
	TitleEl      string
	TitleEn      string
	CategoriesEl []int
	CategoriesEn []int
	FilePath     string
	IsBusiness   bool
}
