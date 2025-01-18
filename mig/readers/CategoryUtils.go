package readers

import (
	"app/adapter"
	"app/model"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/xuri/excelize/v2"
	"os"
)

type CategoryUtils struct {
	filePath string
	adapter  *adapter.StrapiAdapter
}

func NewCategoryUtils() *CategoryUtils {
	return &CategoryUtils{filePath: os.Getenv("MIGRATED_FILES_PATH"),
		adapter: adapter.NewStrapiAdapter()}
}

func (reader *CategoryUtils) Create(request model.StrapiCategoryCreateWrapper) int {

	return reader.adapter.CreateCategory(request).Data.Id
}
func (reader *CategoryUtils) Localize(id int, category model.StrapiCategoryCreate) int {
	return reader.adapter.Localize(id, category).Id
}

func (reader *CategoryUtils) Parent(id int, parenting model.StrapiCategoryParenting) {
	reader.adapter.Parent(id, parenting)
}

func (reader *CategoryUtils) Read() model.Levels {
	f, _ := excelize.OpenFile(reader.filePath)
	rowsEl, _ := f.GetRows("el")
	rowsEn, _ := f.GetRows("en")
	level0Map := make(map[string]model.Category)
	level1Map := make(map[string]model.Category)
	level2Map := make(map[string]model.Category)
	level3Map := make(map[string]model.Category)
	level4Map := make(map[string]model.Category)
	level5Map := make(map[string]model.Category)

	for index, val := range rowsEl {
		if index == 0 {
			continue
		}

		if len(val[0]) > 0 {
			if len(val[0]) == 0 {
				continue
			}
			mv := ToGreeklish(val[0])
			if v, ok := level0Map[mv]; ok {
				if len(rowsEn[index][0]) > 0 {
					v.NameEn = rowsEn[index][0]
					level0Map[mv] = v
				}
			} else {
				cat := model.Category{
					IdEl:   -1,
					IdEn:   -1,
					NameEl: val[0],
					NameEn: "",
					CatId:  mv,
					Parent: hashset.New(),
				}
				if len(rowsEn[index][0]) > 0 {
					cat.NameEn = rowsEn[index][0]
				} else {
					cat.NameEn = rowsEl[index][0]
				}
				level0Map[mv] = cat
			}

		}
		if len(val[1]) > 0 {
			if len(val[2]) == 0 {
				continue
			}
			mv := ToGreeklish(val[1])
			mvp := ToGreeklish(val[0])
			if v, ok := level1Map[mv]; ok {
				if len(rowsEn[index][1]) > 0 {
					v.NameEn = rowsEn[index][1]
				}
				v.Parent.Add(mvp)
				level1Map[mv] = v
			} else {
				cat := model.Category{
					IdEl:   -1,
					IdEn:   -1,
					NameEl: val[1],
					NameEn: "",
					CatId:  mv,
					Parent: hashset.New(),
				}
				if len(rowsEn[index][1]) > 0 {
					cat.NameEn = rowsEn[index][1]
				} else {
					cat.NameEn = rowsEl[index][1]
				}
				cat.Parent.Add(mvp)
				level1Map[mv] = cat
			}
		}
		if len(val[2]) > 0 {
			if len(val[3]) == 0 {
				continue
			}
			mv := ToGreeklish(val[2])
			mvp := ToGreeklish(val[1])
			if v, ok := level2Map[val[2]]; ok {
				if len(rowsEn[index][2]) > 0 {
					v.NameEn = rowsEn[index][2]
				}
				v.Parent.Add(mvp)
				level2Map[mv] = v
			} else {
				cat := model.Category{
					IdEl:   -1,
					IdEn:   -1,
					NameEl: val[2],
					NameEn: "",
					CatId:  mv,
					Parent: hashset.New(),
				}
				if len(rowsEn[index][2]) > 0 {
					cat.NameEn = rowsEn[index][2]
				} else {
					cat.NameEn = rowsEl[index][2]
				}
				cat.Parent.Add(mvp)
				level2Map[mv] = cat
			}

		}
		if len(val[3]) > 0 {
			if len(val[4]) == 0 {
				continue
			}
			mv := ToGreeklish(val[3])
			mvp := ToGreeklish(val[2])
			if v, ok := level3Map[mv]; ok {
				if len(rowsEn[index][3]) > 0 {
					v.NameEn = rowsEn[index][3]
				}
				v.Parent.Add(mvp)
				level3Map[mv] = v
			} else {
				cat := model.Category{
					IdEl:   -1,
					IdEn:   -1,
					NameEl: val[3],
					NameEn: "",
					CatId:  mv,
					Parent: hashset.New(),
				}
				if len(rowsEn[index][3]) > 0 {
					cat.NameEn = rowsEn[index][3]
				} else {
					cat.NameEn = rowsEl[index][3]
				}
				cat.Parent.Add(mvp)
				level3Map[mv] = cat
			}

		}
		if len(val[4]) > 0 {
			if len(val[5]) == 0 {
				continue
			}
			mv := ToGreeklish(val[4])
			mvp := ToGreeklish(val[3])
			if v, ok := level4Map[mv]; ok {
				if len(rowsEn[index][4]) > 0 {
					v.NameEn = rowsEn[index][4]
				}
				v.Parent.Add(mvp)
				level4Map[mv] = v
			} else {
				cat := model.Category{
					IdEl:   -1,
					IdEn:   -1,
					NameEl: val[4],
					NameEn: "",
					CatId:  mv,
					Parent: hashset.New(),
				}
				if len(rowsEn[index][4]) > 0 {
					cat.NameEn = rowsEn[index][4]
				} else {
					cat.NameEn = rowsEl[index][4]
				}
				cat.Parent.Add(mvp)
				level4Map[mv] = cat
			}

		}
		if len(val[5]) > 0 {
			if len(val[6]) == 0 {
				continue
			}
			mv := ToGreeklish(val[5])
			mvp := ToGreeklish(val[4])
			if v, ok := level5Map[mv]; ok {
				if len(rowsEn[index][5]) > 0 {
					v.NameEn = rowsEn[index][5]
				}
				v.Parent.Add(mvp)
				level5Map[mv] = v
			} else {
				cat := model.Category{
					IdEl:   -1,
					IdEn:   -1,
					NameEl: val[5],
					NameEn: "",
					CatId:  mv,
					Parent: hashset.New(),
				}
				if len(rowsEn[index][5]) > 0 {
					cat.NameEn = rowsEn[index][5]
				} else {
					cat.NameEn = rowsEl[index][5]
				}
				cat.Parent.Add(mvp)
				level5Map[mv] = cat
			}
		}
	}

	return model.Levels{
		Level0Map: level0Map,
		Level1Map: level1Map,
		Level2Map: level2Map,
		Level3Map: level3Map,
		Level4Map: level4Map,
		Level5Map: level5Map,
	}
}

func (reader *CategoryUtils) SubmitLevels(levels model.Levels) map[string]model.Category {
	allCategories := make(map[string]model.Category)
	/** Level0Map **/
	for k, v := range levels.Level0Map {
		requestEl := model.StrapiCategoryCreateWrapper{
			Data: model.StrapiCategoryCreate{
				Title:      v.NameEl,
				CategoryId: k,
				Locale:     "el",
			},
		}
		res := reader.Create(requestEl)
		v.IdEl = res
		if len(v.NameEn) > 0 {
			loc := reader.Localize(res, model.StrapiCategoryCreate{
				Title:      v.NameEn,
				CategoryId: k,
				Locale:     "en",
			})
			v.IdEn = loc
		}
		levels.Level0Map[k] = v
		allCategories[k] = v
	}
	/** Level1Map **/
	for k, v := range levels.Level1Map {
		requestEl := model.StrapiCategoryCreateWrapper{
			Data: model.StrapiCategoryCreate{
				Title:      v.NameEl,
				CategoryId: k,
				Locale:     "el",
			},
		}
		res := reader.Create(requestEl)
		v.IdEl = res
		if v.Parent.Size() > 0 {
			parentIds := make([]int, 0)
			for _, prnt := range v.Parent.Values() {
				pt := levels.Level0Map[fmt.Sprintf("%v", prnt)].IdEl
				parentIds = append(parentIds, pt)
			}
			reader.Parent(res, model.StrapiCategoryParenting{ParentCategories: parentIds})
		}
		if len(v.NameEn) > 0 {
			loc := reader.Localize(res, model.StrapiCategoryCreate{
				Title:      v.NameEn,
				CategoryId: k,
				Locale:     "en",
			})
			v.IdEn = loc
			if v.Parent.Size() > 0 {
				parentIds := make([]int, 0)
				for _, prnt := range v.Parent.Values() {

					pt := levels.Level0Map[fmt.Sprintf("%v", prnt)].IdEn
					parentIds = append(parentIds, pt)
				}
				reader.Parent(loc, model.StrapiCategoryParenting{ParentCategories: parentIds})
			}
		} else {
			loc := reader.Localize(res, model.StrapiCategoryCreate{
				Title:      v.NameEl,
				CategoryId: k,
				Locale:     "en",
			})
			v.IdEn = loc
			if v.Parent.Size() > 0 {
				parentIds := make([]int, 0)
				for _, prnt := range v.Parent.Values() {
					pt := levels.Level0Map[fmt.Sprintf("%v", prnt)].IdEn
					parentIds = append(parentIds, pt)
				}
				reader.Parent(loc, model.StrapiCategoryParenting{ParentCategories: parentIds})
			}
		}
		levels.Level1Map[k] = v
		allCategories[k] = v
	}
	/** Level2Map **/
	for k, v := range levels.Level2Map {
		requestEl := model.StrapiCategoryCreateWrapper{
			Data: model.StrapiCategoryCreate{
				Title:      v.NameEl,
				CategoryId: k,
				Locale:     "el",
			},
		}
		res := reader.Create(requestEl)
		v.IdEl = res
		if v.Parent.Size() > 0 {
			parentIds := make([]int, 0)
			for _, prnt := range v.Parent.Values() {
				pt := levels.Level1Map[fmt.Sprintf("%v", prnt)].IdEl
				parentIds = append(parentIds, pt)
			}
			reader.Parent(res, model.StrapiCategoryParenting{ParentCategories: parentIds})
		}
		if len(v.NameEn) > 0 {
			loc := reader.Localize(res, model.StrapiCategoryCreate{
				Title:      v.NameEn,
				CategoryId: k,
				Locale:     "en",
			})
			v.IdEn = loc
			if v.Parent.Size() > 0 {
				parentIds := make([]int, 0)
				for _, prnt := range v.Parent.Values() {

					pt := levels.Level1Map[fmt.Sprintf("%v", prnt)].IdEn
					parentIds = append(parentIds, pt)
				}
				reader.Parent(loc, model.StrapiCategoryParenting{ParentCategories: parentIds})
			}
		} else {
			loc := reader.Localize(res, model.StrapiCategoryCreate{
				Title:      v.NameEl,
				CategoryId: k,
				Locale:     "en",
			})
			v.IdEn = loc
			if v.Parent.Size() > 0 {
				parentIds := make([]int, 0)
				for _, prnt := range v.Parent.Values() {
					pt := levels.Level1Map[fmt.Sprintf("%v", prnt)].IdEn
					parentIds = append(parentIds, pt)
				}
				reader.Parent(loc, model.StrapiCategoryParenting{ParentCategories: parentIds})
			}
		}
		levels.Level2Map[k] = v
		allCategories[k] = v
	}
	/** Level3Map **/
	for k, v := range levels.Level3Map {
		requestEl := model.StrapiCategoryCreateWrapper{
			Data: model.StrapiCategoryCreate{
				Title:      v.NameEl,
				CategoryId: k,
				Locale:     "el",
			},
		}
		res := reader.Create(requestEl)
		v.IdEl = res
		if v.Parent.Size() > 0 {
			parentIds := make([]int, 0)
			for _, prnt := range v.Parent.Values() {
				pt := levels.Level2Map[fmt.Sprintf("%v", prnt)].IdEl
				parentIds = append(parentIds, pt)
			}
			reader.Parent(res, model.StrapiCategoryParenting{ParentCategories: parentIds})
		}
		if len(v.NameEn) > 0 {
			loc := reader.Localize(res, model.StrapiCategoryCreate{
				Title:      v.NameEn,
				CategoryId: k,
				Locale:     "en",
			})
			v.IdEn = loc
			if v.Parent.Size() > 0 {
				parentIds := make([]int, 0)
				for _, prnt := range v.Parent.Values() {

					pt := levels.Level2Map[fmt.Sprintf("%v", prnt)].IdEn
					parentIds = append(parentIds, pt)
				}
				reader.Parent(loc, model.StrapiCategoryParenting{ParentCategories: parentIds})
			}
		} else {
			loc := reader.Localize(res, model.StrapiCategoryCreate{
				Title:      v.NameEl,
				CategoryId: k,
				Locale:     "en",
			})
			v.IdEn = loc
			if v.Parent.Size() > 0 {
				parentIds := make([]int, 0)
				for _, prnt := range v.Parent.Values() {
					pt := levels.Level2Map[fmt.Sprintf("%v", prnt)].IdEn
					parentIds = append(parentIds, pt)
				}
				reader.Parent(loc, model.StrapiCategoryParenting{ParentCategories: parentIds})
			}
		}
		levels.Level3Map[k] = v
		allCategories[k] = v
	}
	/** Level4Map **/
	for k, v := range levels.Level4Map {
		requestEl := model.StrapiCategoryCreateWrapper{
			Data: model.StrapiCategoryCreate{
				Title:      v.NameEl,
				CategoryId: k,
				Locale:     "el",
			},
		}
		res := reader.Create(requestEl)
		v.IdEl = res
		if v.Parent.Size() > 0 {
			parentIds := make([]int, 0)
			for _, prnt := range v.Parent.Values() {
				pt := levels.Level3Map[fmt.Sprintf("%v", prnt)].IdEl
				parentIds = append(parentIds, pt)
			}
			reader.Parent(res, model.StrapiCategoryParenting{ParentCategories: parentIds})
		}
		if len(v.NameEn) > 0 {
			loc := reader.Localize(res, model.StrapiCategoryCreate{
				Title:      v.NameEn,
				CategoryId: k,
				Locale:     "en",
			})
			v.IdEn = loc
			if v.Parent.Size() > 0 {
				parentIds := make([]int, 0)
				for _, prnt := range v.Parent.Values() {

					pt := levels.Level3Map[fmt.Sprintf("%v", prnt)].IdEn
					parentIds = append(parentIds, pt)
				}
				reader.Parent(loc, model.StrapiCategoryParenting{ParentCategories: parentIds})
			}
		} else {
			loc := reader.Localize(res, model.StrapiCategoryCreate{
				Title:      v.NameEl,
				CategoryId: k,
				Locale:     "en",
			})
			v.IdEn = loc
			if v.Parent.Size() > 0 {
				parentIds := make([]int, 0)
				for _, prnt := range v.Parent.Values() {
					pt := levels.Level3Map[fmt.Sprintf("%v", prnt)].IdEn
					fmt.Println("Retrieve parent en", prnt, pt)
					parentIds = append(parentIds, pt)
				}
				reader.Parent(loc, model.StrapiCategoryParenting{ParentCategories: parentIds})
			}
		}
		levels.Level4Map[k] = v
		allCategories[k] = v
	}
	/** Level5Map **/
	for k, v := range levels.Level5Map {
		requestEl := model.StrapiCategoryCreateWrapper{
			Data: model.StrapiCategoryCreate{
				Title:      v.NameEl,
				CategoryId: k,
				Locale:     "el",
			},
		}
		res := reader.Create(requestEl)
		v.IdEl = res
		if v.Parent.Size() > 0 {
			parentIds := make([]int, 0)
			for _, prnt := range v.Parent.Values() {
				pt := levels.Level4Map[fmt.Sprintf("%v", prnt)].IdEl
				fmt.Println("Retrieve parent el", prnt, pt)
				parentIds = append(parentIds, pt)
			}
			reader.Parent(res, model.StrapiCategoryParenting{ParentCategories: parentIds})
		}
		if len(v.NameEn) > 0 {
			loc := reader.Localize(res, model.StrapiCategoryCreate{
				Title:      v.NameEn,
				CategoryId: k,
				Locale:     "en",
			})
			v.IdEn = loc
			if v.Parent.Size() > 0 {
				parentIds := make([]int, 0)
				for _, prnt := range v.Parent.Values() {
					pt := levels.Level4Map[fmt.Sprintf("%v", prnt)].IdEn
					fmt.Println("Retrieve parent en", prnt, pt)
					parentIds = append(parentIds, pt)
				}
				reader.Parent(loc, model.StrapiCategoryParenting{ParentCategories: parentIds})
			}
		} else {
			loc := reader.Localize(res, model.StrapiCategoryCreate{
				Title:      v.NameEl,
				CategoryId: k,
				Locale:     "en",
			})
			v.IdEn = loc
			if v.Parent.Size() > 0 {
				parentIds := make([]int, 0)
				for _, prnt := range v.Parent.Values() {
					pt := levels.Level4Map[fmt.Sprintf("%v", prnt)].IdEn
					parentIds = append(parentIds, pt)
				}
				reader.Parent(loc, model.StrapiCategoryParenting{ParentCategories: parentIds})
			}
		}
		levels.Level5Map[k] = v
		allCategories[k] = v
	}

	return allCategories
}
